package main

import (
	"github.com/osoderholm/svenska-yle-bot/database"
	"log"
	"time"

	// Telegram Bot api wrapper
	tb "gopkg.in/tucnak/telebot.v2"
)

// CreateNotifier starts a goroutine that with 1 minute intervals tries to notify subscribers with latest articles.
// If an article ID is less than or equal to last sent ID, skip it. The maximal amount of articles is currently 10.
func CreateNotifier(db *database.DB, b *tb.Bot) chan struct{} {
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				articles, err := GetLatestArticles(db, 10)
				if err != nil {
					log.Println("NotifySubscribers:", err)
				}

				subscribers, err := GetNotifiableSubscribers(db)
				if err != nil {
					log.Println("NotifySubscribers:", err)
					break
				}
				for _, s := range subscribers {

					chat, err := b.ChatByID(s.ChatID)
					if err != nil {
						log.Println("NotifySubscribers:", err)
						continue
					}

					lastArticle := s.LastArticleID

					for _, a := range articles {
						if a.ID <= s.LastArticleID {
							continue
						}
						sendArticleWithImage(b, chat, *a)
						lastArticle = a.ID
					}
					if len(articles) > 0 {
						s.LastArticleID = lastArticle
						s.SetLastNotified(time.Now())
					}
					if err := s.Update(db); err != nil {
						log.Println("NotifySubscribers:", err)
					}
				}
			}

			time.Sleep(1 * time.Minute)
		}
	}()

	return quit
}
