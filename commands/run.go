package commands

import (
	"github.com/tothmate90/news-scraper/api"
	"github.com/tothmate90/news-scraper/config"
	"github.com/tothmate90/news-scraper/elasticsearch"
	"github.com/tothmate90/news-scraper/mysql"
)

func RunServer(configFile string) error {
	config, err := config.ReadJson(configFile)
	// MySQL section
	_, err = mysql.New(config.Conn)
	if err != nil {
		return err
	}
	// Elastic section
	handler, err := elasticsearch.New(config.Conn, "")
	if err != nil {
		return err
	}
	// Api section
	api.New("8080", handler, config)
	return err
}
