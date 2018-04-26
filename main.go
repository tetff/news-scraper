package main

import (
	"net/url"

	"github.com/tothmate90/news-scraper/elasticsearch"
	"github.com/tothmate90/news-scraper/mysql"
	"github.com/tothmate90/news-scraper/newsapi"
)

const conn = "root:toor@tcp(127.0.0.1:3306)/newsapitest?charset=utf8&parseTime=True&loc=Local"

func main() {
	var values = url.Values{}
	values.Add("country", "us")
	values.Add("category", "business")
	result, err := newsapi.GetTopHeadlines(values)
	if err != nil {
		panic(err)
	}

	// MySQL section
	handler, err := mysql.New(conn)
	if err != nil {
		panic(err)
	}
	handler.Save(result.Articles)

	// Elastic section
	esHandler, err := elasticsearch.New("")
	if err != nil {
		panic(err)
	}
	err = esHandler.Create("articles")
	if err != nil {
		panic(err)
	}
	err = esHandler.Save(result.Articles, "articles")
	if err != nil {
		panic(err)
	}
}
