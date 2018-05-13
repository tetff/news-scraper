package utils

import (
	"reflect"

	"github.com/olivere/elastic"
	"github.com/tothmate90/news-scraper/newsapi"
)
// Translator Converts ES searchresults into our Article type.
func Translator(sr *elastic.SearchResult) []newsapi.Article {
	articles := []newsapi.Article{}
	for _, item := range sr.Each(reflect.TypeOf(newsapi.Article{})) {
		article := item.(newsapi.Article)
		articles = append(articles, article)
	}
	return articles
}
