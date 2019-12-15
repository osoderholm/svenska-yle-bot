package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/osoderholm/svenska-yle-bot/database"
	"log"
	"regexp"
	"strings"
	"time"
)

const timeOut = 5

func CreateArticleFetcher(db *database.DB) chan struct{} {
	quit := make(chan struct{})

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
				fp := gofeed.NewParser()
				feed, _ := fp.ParseURL("http://svenska.yle.fi/nyheter/senaste-nytt.rss")
				for i := len(feed.Items) - 1; i >= 0; i-- {
					item := feed.Items[i]
					article, err := NewArticle(item.Title, item.Link, "", item.Published)
					if err != nil {
						continue
					}

					m, err := regexp.Compile(`img src="([^"]+)"`)
					if err == nil {
						image := m.FindStringSubmatch(item.Description)
						if len(image) > 1 {
							article.Image = image[1]
						}
					}

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
