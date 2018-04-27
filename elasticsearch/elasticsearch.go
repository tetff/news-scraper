package elasticsearch

import (
	"context"
	"fmt"

	"github.com/tothmate90/news-scraper/newsapi"

	"github.com/olivere/elastic"
)

type Handler struct {
	conn   string
	Client *elastic.Client
}

func New(conn string) (Handler, error) {
	var handler Handler
	var err error
	handler.conn = conn
	handler.Client, err = elastic.NewClient()
	return handler, err
}

func (h *Handler) Create(name string) error {
	_, err := h.Client.CreateIndex(name).Do(context.Background())
	return err
}

func (h *Handler) Save(articles []newsapi.Article, name string) error {
	for id, article := range articles {
		_, err := h.Client.Index().
			Index(name).
			Type("doc").
			Id(fmt.Sprintf("%d", id+1)).
			BodyJson(article).
			Refresh("wait_for").
			Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) GetAll(from, size int, name string) (error, *elastic.SearchResult) {
	termQuery := elastic.NewMatchAllQuery()
	searchResult, err := h.Client.Search().
		Index(name).             // search in index "tweets"
		Query(termQuery).        // specify the query
		From(from).Size(size).   // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	fmt.Println(searchResult)
	return err, searchResult
}

func (h *Handler) Get(id string) (error, *elastic.SearchResult) {
	termQuery := elastic.NewTermQuery("id", id)
	searchResult, err := h.Client.Search().
		Index("tweets").         // search in index "tweets"
		Query(termQuery).        // specify the query
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	return err, searchResult
}

func (h *Handler) Post(article newsapi.Article, name string) error {
	id := elastic.NewMaxAggregation().Field("id")
	_, err := h.Client.Index().
		Index(name).
		Type("doc").
		// Id(fmt.Sprintf("%d", id+1)).
		BodyJson(article).
		Refresh("wait_for").
		Do(context.Background())
	fmt.Println(id)
	return err
}
