package main

import (
	"github.com/osoderholm/svenska-yle-bot/database"
	"time"
)

// Article represents a news article from the feed.
// The ID is generated as the article is stored in the database.
type Article struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Link      string    `db:"link"`
	Image     string    `db:"image"`
	Published time.Time `db:"published"`
}

// NewArticle creates a new Article and converts the time stamp.
// Primarily used for converting feed item into Article.
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

// GetLatestArticles fetches specified amount of newest articles from the database.
// The returned order is oldest first.
func GetLatestArticles(db *database.DB, limit int) ([]*Article, error) {
	q := `select * from (select * from bot.articles order by published desc limit ?) t order by published;`

	articles := make([]*Article, 0)

	err := db.Select(&articles, db.Rebind(q), limit)

	return articles, err
}

// Insert creates a database entry of Article if there is none with the link.
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
