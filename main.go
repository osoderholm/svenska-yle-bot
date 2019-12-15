package main

import (
	"fmt"
	"github.com/osoderholm/svenska-yle-bot/database"
	"log"
	"os"
	"strings"
	"time"

	// Telegram Bot api wrapper
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	fmt.Print("Opening DB connection... ")
	db, dbErr := database.Open(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	if dbErr != nil {
		fmt.Println("FAIL")
		log.Fatal(dbErr)
	}
	defer db.Close()
	fmt.Println("DONE")

	fmt.Print("Performing DB migrations... ")
	if err := db.PerformMigrations(); err != nil {
		fmt.Println("FAIL")
		log.Fatal(err)
	}
	fmt.Println("DONE")

	fmt.Print("Subscribing channels to bot... ")
	if err := addChannelSubscribers(db); err != nil {
		fmt.Println("FAIL")
		log.Fatal(err)
	}
	fmt.Println("DONE")

	fmt.Print("Starting news fetcher... ")
	fetcher := CreateArticleFetcher(db)
	fmt.Println("DONE")

	fmt.Print("Connecting to Telegram... ")
	b, botErr := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TG_API_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if botErr != nil {
		fmt.Println("FAIL")
		log.Fatal(botErr)
	}
	fmt.Println("DONE")

	fmt.Print("Setting up endpoint handlers... ")
	CreateHandlers(db, b)
	fmt.Println("DONE")

	fmt.Print("Starting notifier... ")
	notifier := CreateNotifier(db, b)
	fmt.Println("DONE")

	fmt.Println("Bot is running")
	b.Start()

	fmt.Println("Bot is stopping")

	close(notifier)
	close(fetcher)

	fmt.Println("Bot stopped")
}

// addChannelSubscribers reads the environmental variable CHANNEL_SUBSCRIBERS
// and attempts to add the specified, comma separated channels to subscribers.
func addChannelSubscribers(db *database.DB) error {
	// comma separated list
	channels := strings.Split(strings.ReplaceAll(os.Getenv("CHANNEL_SUBSCRIBERS"), "\"", ""), ",")

	// assume list members start with @

	subscribers, err := GetChannelSubscribers(db)
	if err != nil {
		return err
	}

	foundChannels := make([]*Subscriber, 0)

	for _, c := range channels {
		var found bool
		if len(c) == 0 {
			continue
		}
		for _, s := range subscribers {
			if c == s.ChatID {
				foundChannels = append(foundChannels, s)
				found = true
				break
			}
		}
		if found {
			continue
		}
		// Channel is not subscribed, create
		_ = NewSubscriber(c, 1).Insert(db)
	}

	if len(foundChannels) != len(subscribers) {
		for _, s := range subscribers {
			var found bool
			for _, f := range foundChannels {
				if f.ID == s.ID {
					found = true
					break
				}
			}
			if found {
				continue
			}
			// Subscriber was not found, delete
			_ = s.Delete(db)
		}
	}

	return nil
}
