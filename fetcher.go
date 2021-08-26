package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/osoderholm/svenska-yle-bot/database"
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const timeOut = 5 // Minutes

const feedURL = "https://svenska.yle.fi/nyheter/senaste-nytt.rss"

// CreateArticleFetcher starts a goroutine that with specified intervals fetches the latest articles using an RSS feed.
// The fetcher returns a chan and can be killed by closing the returned chan.
func CreateArticleFetcher(db *database.DB) chan struct{} {
	quit := make(chan struct{})

	fp := gofeed.NewParser()

	go func() {
		i := 0
		for {
			select {
			case <-quit:
				log.Println("CreateArticleFetcher: Got quit signal")
				return
			default:
				if i > 0 {
					break
				}
				feed, feedErr := fp.ParseURL(feedURL)
				if feedErr != nil {
					log.Println(feedErr)
					break
				}
				for i := len(feed.Items) - 1; i >= 0; i-- {
					item := feed.Items[i]
					article, err := NewArticle(html.UnescapeString(item.Title), item.Link, "", item.Published)
					if err != nil {
						continue
					}

					// Fetch article page content to determine cover image
					article.Image = fetchImageURL(article.Link)

					if err := article.Insert(db); err != nil {
						if !strings.Contains(err.Error(), "(SQLSTATE 23505)") {
							log.Println("CreateArticleFetcher:", err)
						}
						continue
					}

					log.Println("Stored new article:", article.Title)
				}
			}

			i++

			if i >= timeOut {
				i = 0
			}

			time.Sleep(1 * time.Minute)
		}
	}()

	return quit
}

func fetchFeedAsString(feedURL string) (string, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func fetchImageURL(link string) string {
	ar, arErr := http.Get(link)
	if arErr != nil {
		return ""
	}

	defer ar.Body.Close()

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, ar.Body); err != nil {
		return ""
	}

	m, err := regexp.Compile(`meta\sproperty="og:image"\scontent="([^"]+)"`)
	if err != nil {
		return ""
	}

	image := m.FindStringSubmatch(buf.String())
	if len(image) > 1 {
		return image[1]
	}

	return ""
}
