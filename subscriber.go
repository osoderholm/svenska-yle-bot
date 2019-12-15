package main

import (
	"errors"
	"github.com/osoderholm/svenska-yle-bot/database"
	"time"
)

// Subscriber represents a Telegram user, group or a channel that will receive articles with specified intervals.
// ID is auto generated when stored in the database.
type Subscriber struct {
	ID             int64      `db:"id"`
	ChatID         string     `db:"chat_id"`
	UpdateInterval int        `db:"update_interval"`
	LastArticleID  int64      `db:"last_article_id"`
	LastNotified   *time.Time `db:"last_notified"`
}

// NewSubscriber creates a simple instance of Subscriber with specified data.
func NewSubscriber(chatID string, interval int) *Subscriber {
	return &Subscriber{
		ChatID:         chatID,
		UpdateInterval: interval,
	}
}

// SetLastNotified sets the LastNotified time.
func (s *Subscriber) SetLastNotified(t time.Time) {
	s.LastNotified = &t
}

// GetNotifiableSubscribers returns a list of Subscribers from the database
// that should be notified based on the update interval and current time.
func GetNotifiableSubscribers(db *database.DB) ([]*Subscriber, error) {
	subscribers := make([]*Subscriber, 0)

	q := `select * from bot.subscribers where last_notified is null or extract(minute from now()-last_notified) > update_interval;`

	err := db.Select(&subscribers, q)

	return subscribers, err
}

// GetSubscriberByChatID returns a single Subscriber from the database, based on chat ID.
func GetSubscriberByChatID(db *database.DB, chatID string) (*Subscriber, error) {
	q := `select * from bot.subscribers where chat_id = :chat_id limit 1;`

	stmt, stmtErr := db.PrepareNamed(q)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer stmt.Close()

	s := &Subscriber{}

	getErr := stmt.Get(s, Subscriber{ChatID: chatID})

	return s, getErr
}

// GetChannelSubscribers returns a list of Subscribers from the database that are Telegram channels.
// The filter is based on the fact that the channels chat ID begins with @.
func GetChannelSubscribers(db *database.DB) ([]*Subscriber, error) {
	subscribers := make([]*Subscriber, 0)

	q := `select * from bot.subscribers where chat_id like '@%';`

	err := db.Select(&subscribers, q)

	return subscribers, err
}

// Insert creates a new entry in the database if one with the chat ID does not exist.
func (s *Subscriber) Insert(db *database.DB) error {
	q := `insert into bot.subscribers (chat_id, update_interval, last_article_id) values (:chat_id, :update_interval, :last_article_id);`

	stmt, stmtErr := db.PrepareNamed(q)
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	_, execErr := stmt.Exec(s)

	return execErr
}

// Update updates an entry of Subscriber in the database based on Subscriber's ID.
// The ID is therefore required to be greater than 0.
func (s *Subscriber) Update(db *database.DB) error {
	if s.ID == 0 {
		return errors.New("subscriber does not exist")
	}

	q := `update bot.subscribers set (update_interval, last_article_id, last_notified) 
    	= (:update_interval, :last_article_id, :last_notified) where id = :id;`

	stmt, stmtErr := db.PrepareNamed(q)
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	_, execErr := stmt.Exec(s)

	return execErr
}

// Delete will delete an entry of Subscriber from the database based on Subscriber's ID.
func (s *Subscriber) Delete(db *database.DB) error {
	if s.ID == 0 {
		return errors.New("subscriber does not exist")
	}

	q := `delete from bot.subscribers where id = :id;`

	stmt, stmtErr := db.PrepareNamed(q)
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	_, execErr := stmt.Exec(s)

	return execErr
}
