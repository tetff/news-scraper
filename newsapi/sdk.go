package newsapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const apiKey = "9699ebb3ce234b649aedac95e057bdc4"
const host = "https://newsapi.org"

// Wrapper ...
type Wrapper struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
	Code         string    `json:"code,omitempty"`
	Message      string    `json:"message,omitempty"`
}

// Article ...
type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

// Source ...
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Parse ...
func Parse(data []byte) (Wrapper, error) {
	var wrapper Wrapper
	err := json.Unmarshal(data, &wrapper)
	return wrapper, err
}

// GetTopHeadlines ...
func GetTopHeadlines(values url.Values) (Wrapper, error) {
	var wrapper Wrapper
	values.Add("apiKey", apiKey)
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
