package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/tothmate90/news-scraper/newsapi"
)
// Handler Handling all MySQL related methods. Uses Jinzhu's library.
type Handler struct {
	DB *gorm.DB
}
// New Creates connection with the MySQL database. Get's required connection info from the config file.
func New(conn string) (Handler, error) {
	db, err := gorm.Open("mysql", conn)
	return Handler{
		DB: db,
	}, err
}
// Migrate Migrates table based on the model.
func (h *Handler) Migrate() {
	h.DB.AutoMigrate(&Article{})
}
// Save Uploads an array of articles to the MySQL server.
func (h *Handler) Save(articles []newsapi.Article) {
	for _, article := range articles {
		a := translator(article)
		h.DB.Create(&a)
	}
}
