package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/tothmate90/news-scraper/config"
	"github.com/tothmate90/news-scraper/newsapi"
)

var log = logrus.New()

// Save saves the result from newsapi.org in the given category and country
// to the file in the results folder.
func Save(fileName, category, country, configFile string) error {
	config, err := config.ReadJSON(configFile)
	if err != nil {
		return err
	}
	values := url.Values{}
	values.Add("country", country)
	values.Add("category", category)
	result, err := newsapi.GetTopHeadlines(values, config)
	if err != nil {
		log.Error(err)
		return err
	}
	err = writeFile(fileName, result.Articles)
	if err != nil {
		log.Error(err)
		return err
	}
	return err
}

func writeFile(fileName string, articles []newsapi.Article) error {
	articleJSON, err := json.Marshal(articles)
	if err != nil {
		log.Error(err)
		return err
	}
	err = ioutil.WriteFile(fileName+".json", articleJSON, 0644)
	if err != nil {
		log.Error(err)
		return err
	}
	return err
}
