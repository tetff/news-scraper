package mysql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tothmate90/news-scraper/newsapi"
)

func TestTranslator(t *testing.T) {
	testArticle := newsapi.Article{
		Source: newsapi.Source{
			Name: "CNBC",
		},
		Author:      "Sara Salinas",
		Title:       "Test title",
		Description: "Test description",
		URL:         "Test URL",
		URLToImage:  "Test URLToImage",
		PublishedAt: time.Now(),
	}
	article := translator(testArticle)
	assert.Equal(t, testArticle.Source.Name, article.Source)
	assert.Equal(t, testArticle.URL, article.URL)
	assert.Equal(t, testArticle.Description, article.Description)
}
