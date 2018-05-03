package elasticsearch

import (
	"context"
	"fmt"

	"github.com/tothmate90/news-scraper/newsapi"

	"github.com/olivere/elastic"
)

type Handler struct {
	conn   string
	name string
	Client *elastic.Client
}

func New(conn, name string) (Handler, error) {
	var handler Handler
	var err error
	handler.conn = conn
	handler.name = name
	handler.Client, err = elastic.NewClient()
	return handler, err
}

func (h *Handler) Create() error {
	_, err := h.Client.CreateIndex(h.name).Do(context.Background())
	return err
}

func (h *Handler) GetAll(from, size int) (*elastic.SearchResult, error) {
	termQuery := elastic.NewMatchAllQuery()
	searchResult, err := h.Client.Search().
		Index(h.name).
		Query(termQuery).
		From(from).Size(size).
		Pretty(true).
		Do(context.Background())
	fmt.Println(searchResult)
	return searchResult, err
}

func (h *Handler) Get(id string) (*elastic.SearchResult, error) {
	termQuery := elastic.NewTermQuery("id", id)
	searchResult, err := h.Client.Search().
		Index(h.name).
		Query(termQuery).
		Pretty(true).
		Do(context.Background())
	return searchResult, err
}

func (h *Handler) Post(articles []newsapi.Article) error {
	idIF, err := elastic.NewMaxAggregation().Field("id").Source()
	if err != nil {
		return err
	}
	id := idIF.(int)
	if idIF == nil {
		id = 0
	}
	for _, article := range articles {
		_, err := h.Client.Index().
			Index(h.name).
			Type("doc").
			Id(fmt.Sprintf("%d", id + 1)).
			BodyJson(article).
			Refresh("wait_for").
			Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
