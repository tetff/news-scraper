package elasticsearch

import (
	"context"
	"fmt"

	"github.com/tothmate90/news-scraper/newsapi"

	"github.com/olivere/elastic"
)
// Handler Interface managing all the functions related to ES, based on Olivere's elastic library.
type Handler interface {
	Create() error
	GetAll(from, size int) (*elastic.SearchResult, error)
	Get(id string) (*elastic.SearchResult, error)
	Post(articles []newsapi.Article) error
}

type handler struct {
	conn   string
	name   string
	Client *elastic.Client
}
// New Creates the connection with the ES server. Requires the name of the index, and the connection information (stored in the config file).
func New(conn, name string) (Handler, error) {
	var handler handler
	var err error
	handler.conn = conn
	handler.name = name
	handler.Client, err = elastic.NewClient()
	return &handler, err
}
// Create Creates the index.
func (h *handler) Create() error {
	_, err := h.Client.CreateIndex(h.name).Do(context.Background())
	return err
}
// GetAll Returns results. Requires from and size parameters.
func (h *handler) GetAll(from, size int) (*elastic.SearchResult, error) {
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
// Get Returns one result specified by ID.
func (h *handler) Get(id string) (*elastic.SearchResult, error) {
	termQuery := elastic.NewTermQuery("id", id)
	searchResult, err := h.Client.Search().
		Index(h.name).
		Query(termQuery).
		Pretty(true).
		Do(context.Background())
	return searchResult, err
}
// Post Uploads an array of articles to the ES server.
func (h *handler) Post(articles []newsapi.Article) error {
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
