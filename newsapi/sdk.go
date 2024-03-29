package newsapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/tothmate90/news-scraper/config"
)

const host = "https://newsapi.org"

// Wrapper Handles the response from newsapi.org.
type Wrapper struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
	Code         string    `json:"code,omitempty"`
	Message      string    `json:"message,omitempty"`
}

// Article The basic article structure that the API uses.
type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

// Source Source of the article.
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Parse Parses the response from newsapi.org into the wrapper struct.
func Parse(data []byte) (Wrapper, error) {
	var wrapper Wrapper
	err := json.Unmarshal(data, &wrapper)
	return wrapper, err
}

// GetTopHeadlines Default method to get the data. Url values are provided by the API.
func GetTopHeadlines(values url.Values, config config.Config) (Wrapper, error) {
	var wrapper Wrapper
	values.Add("apiKey", config.APIKey)
	client := new(http.Client)
	req, err := http.NewRequest("GET", urlBuilder("/v2/top-headlines"), nil)
	if err != nil {
		return wrapper, err
	}
	req.URL.RawQuery = values.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return wrapper, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wrapper, err
	}
	return Parse(body)
}

func urlBuilder(endpoint string) string {
	return host + endpoint
}
