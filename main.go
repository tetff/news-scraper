package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/tothmate90/news-scraper/commands"
	"github.com/urfave/cli"
)

var (
	cliApp     *cli.App
	configFile string
)

func init() {
	cliApp = cli.NewApp()
	cliApp.Name = "news-api-go"
	cliApp.Usage = "API server storing articles from newsapi.org in Elasticsearch"
}

func main() {
	cliApp.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "run the actual HTTP server",
			Action: func(c *cli.Context) error {
				return commands.RunServer(configFile)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config-file",
					EnvVar:      "NEWSAPI_CONFIG_FILE",
					Usage:       "Location of the config file",
					Value:       "./dev-config.json",
					Destination: &configFile,
				},
			},
		},
	}
	if err := cliApp.Run(os.Args); err != nil {
		log.Print(err)
	}
}
