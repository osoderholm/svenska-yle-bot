package main

import (
	"database/sql"
	"fmt"
	"github.com/osoderholm/svenska-yle-bot/database"
	"strconv"

	// Telegram Bot api wrapper
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

// handler provides methods for handling different slash-commands from the bot
type handler struct {
	*database.DB
	*tb.Bot
}

// CreateHandlers sets up endpoints and maps handler functions to them
func CreateHandlers(db *database.DB, b *tb.Bot) {
	h := &handler{
		DB:  db,
		Bot: b,
	}

	b.Handle("/start", h.handleStart)

	b.Handle("/latest", h.handleLatest)

	b.Handle("/subscribe", h.handleSubscribe)
	b.Handle("/unsubscribe", h.handleUnsubscribe)

}

func (h *handler) handleStart(m *tb.Message) {
	if !m.Private() {
		return
	}

	_, err := h.Bot.Send(m.Sender, "Hejsan!")

	if err != nil {
		log.Println(err)
	}
}

func (h *handler) handleLatest(m *tb.Message) {
	articles, err := GetLatestArticles(h.DB, 5)
	if err != nil || len(articles) == 0 {
		if err != nil {
			log.Println(err)
		}
		_, _ = h.Bot.Send(m.Sender, "Hittade inga artiklar")
		return
	}
	for _, article := range articles {
		sendArticleWithImage(h.Bot, m.Sender, *article)
	}
}

func (h *handler) handleSubscribe(m *tb.Message) {
	interval := 60 // one hour

	chatID := strconv.Itoa(m.Sender.ID)

	s, err := GetSubscriberByChatID(h.DB, chatID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return
	}
	if s != nil && s.ID > 0 {
		_, _ = h.Bot.Send(m.Sender, "Du prenumererar redan! För att avsluta prenumerationen, använd /unsubscribe.")
		return
	}

	newSubscriber := NewSubscriber(chatID, interval)

	if err := newSubscriber.Insert(h.DB); err != nil {
		_, _ = h.Bot.Send(m.Sender, "Kunde inte prenumerera.")
		return
	}

	_, _ = h.Bot.Send(m.Sender, "Du är nu prenumerant och får uppdateringar med ca 1 timmes intervall.\n"+
		"För att avsluta prenumerationen, använd /unsubscribe.")
}

func (h *handler) handleUnsubscribe(m *tb.Message) {
	chatID := strconv.Itoa(m.Sender.ID)

	s, err := GetSubscriberByChatID(h.DB, chatID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return
	}
	if s == nil || s.ID == 0 {
		_, _ = h.Bot.Send(m.Sender, "Du är inte en prenumerant.")
		return
	}

	if err := s.Delete(h.DB); err != nil {
		log.Println(err)
		_, _ = h.Bot.Send(m.Sender, "Kunde inte avsluta prenumeration.")
		return
	}

	_, _ = h.Bot.Send(m.Sender, "Prenumerationen avslutad.\n"+
		"För att börja prenumerera, använd /subscribe.")
}

// sendArticleWithImage attempts to send an article with an image.
// If an error is returned, the same article is sent without the image.
func sendArticleWithImage(b *tb.Bot, to tb.Recipient, article Article) {
	p := article.Published.Format("02.01.2006 kl. 15:04")
	caption := fmt.Sprintf("%s: [%s](%s)", p, article.Title, article.Link)
	if len(article.Image) > 0 {
		photo := &tb.Photo{File: tb.FromURL(article.Image), Caption: caption}
		_, err := b.Send(to, photo, tb.ModeMarkdown, tb.Silent, tb.NoPreview)
		if err != nil {
			log.Println(err)
			_, _ = b.Send(to, caption, tb.ModeMarkdown, tb.Silent, tb.NoPreview)
		}
	} else {
		_, _ = b.Send(to, caption, tb.ModeMarkdown, tb.Silent, tb.NoPreview)
	}
}
