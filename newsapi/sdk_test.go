package newsapi

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tothmate90/news-scraper/config"
)

const testResult = `{
	"status": "ok",
	"totalResults": 20,
	"articles": [
			{
					"source": {
							"id": "cnbc",
							"name": "CNBC"
					},
					"author": "Sara Salinas",
					"title": "LendingClub plunges after the FTC charges the online lender with 'deceiving customers'",
					"description": "LendingClub, a pioneer in online marketplace lending, is facing yet another challenge, this time from the FTC.",
					"url": "https://www.cnbc.com/2018/04/25/lendingclub-lc-tanks-on-ftc-charges-of-deceiving-customers.html",
					"urlToImage": "https://fm.cnbc.com/applications/cnbc.com/resources/img/editorial/2016/06/28/103750569-GettyImages-182563630.1910x1000.jpg",
					"publishedAt": "2018-04-25T18:41:09Z"
			},
			{
					"source": {
							"id": null,
							"name": "Presstelegram.com"
					},
					"author": null,
					"title": "JetBlue reducing Long Beach flights, blames city for halting international plan",
					"description": null,
					"url": "http://www.presstelegram.com/jetblue-reducing-long-beach-flights-blames-city-for-halting-international-plan",
					"urlToImage": null,
					"publishedAt": "2018-04-25T18:35:00Z"
			}
	]

}`

func TestParse(t *testing.T) {
	result, err := Parse([]byte(testResult))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Sara Salinas", result.Articles[0].Author)
	assert.Equal(t, "", result.Articles[1].Source.ID)
	assert.Equal(t, "JetBlue reducing Long Beach flights, blames city for halting international plan", result.Articles[1].Title)
}

func TestGetTopHeadlines(t *testing.T) {
	testConfig, err := config.ReadJson("./../dev-config.json")
	if err!= nil {
		t.Error(err)
	}
	var values = url.Values{}
	values.Add("country", "us")
	values.Add("category", "business")
	result, err := GetTopHeadlines(values, testConfig)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "ok", result.Status)
}
