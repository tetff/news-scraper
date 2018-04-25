package mysql

import (
	"github.com/rs/xid"
	"github.com/tothmate90/news-scraper/newsapi"
)

type Article struct {
	ID          string
	Source      string
	Author      string
	Title       string
	Description string `gorm:"type:longtext"`
	URL         string
	URLToImage  string
	PublishedAt string
}

func translater(article newsapi.Article) Article {
	return Article{
		ID:          xid.New().String(),
		Source:      article.Source.Name,
		Author:      article.Author,
		Title:       article.Title,
		Description: article.Description,
		URL:         article.URL,
		URLToImage:  article.URLToImage,
		PublishedAt: article.PublishedAt.String(),
	}
}
