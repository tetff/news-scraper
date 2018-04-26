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
