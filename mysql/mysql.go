package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tothmate90/news-scraper/newsapi"
)

type Handler struct {
	DB *gorm.DB
}

func New(conn string) (Handler, error) {
	db, err := gorm.Open("mysql", conn)
	return Handler{
		DB: db,
	}, err
}

func (h *Handler) Migrate() {
	h.DB.AutoMigrate(&Article{})
}

func (h *Handler) Save(articles []newsapi.Article) {
	for _, article := range articles {
		a := translater(article)
		h.DB.Create(&a)
	}
}
