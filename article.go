package main

import (
	"github.com/osoderholm/svenska-yle-bot/database"
	"time"
)

type Article struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Link      string    `db:"link"`
	Image     string    `db:"image"`
	Published time.Time `db:"published"`
}

func NewArticle(title, link, image, published string) (*Article, error) {
	t, err := time.Parse(time.RFC1123Z, published)
	if err != nil {
		return nil, err
	}

	return &Article{
		Title:     title,
		Link:      link,
		Image:     image,
		Published: t,
	}, nil
}

func GetLatestArticles(db *database.DB, limit int) ([]*Article, error) {
	q := `select * from (select * from bot.articles order by published desc limit ?) t order by published;`

	articles := make([]*Article, 0)

	err := db.Select(&articles, db.Rebind(q), limit)

	return articles, err
}

func (a *Article) Insert(db *database.DB) error {
	q := `insert into bot.articles (title, link, image, published) values (:title, :link, :image, :published);`

	stmt, stmtErr := db.PrepareNamed(q)
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	_, execErr := stmt.Exec(a)

	return execErr
}
