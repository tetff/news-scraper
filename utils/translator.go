package utils

import (
	"github.com/tothmate90/news-scraper/newsapi"
	"github.com/olivere/elastic"
)

func translater(sr *elastic.SearchResult) []newsapi.Article {
	articles := []newsapi.Article{}
	sr.Each( )
	return articles
}