package api

import (
	"github.com/tothmate90/news-scraper/newsapi"
	"github.com/tothmate90/news-scraper/elasticsearch"
	"net/http"
	"time"
	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	Mux *chi.Mux
	Route *chi.Router
}

func New() Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	return Handler{
		Mux: mux,
	}
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *newsapi.Article
		var err error
		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			article, err = (*elasticsearch.Handler).Get(articleID)
		}
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) Conn() {
	http.ListenAndServe(":3000", h.Mux)
}

func (h *Handler) Route() {
	
}

func (h *Handler) GetAll(from, size int) {

}