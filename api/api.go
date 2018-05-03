package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	"github.com/tothmate90/news-scraper/utils"

	"github.com/tothmate90/news-scraper/elasticsearch"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	Mux *chi.Mux
	EH  *elasticsearch.Handler
}

func New(port string, eH *elasticsearch.Handler) Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	http.ListenAndServe(":"+port, mux)
	return Handler{
		Mux: mux,
		EH:  eH,
	}
}

func (h *Handler) Get(id string) {
	h.Mux.Get("/news-api/ver1.0/", func(w http.ResponseWriter, r *http.Request) {
		sR, err := h.EH.Get(id)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		articleResponse(sR, w)
	})
}

func (h *Handler) GetAll(from, size int) {
	h.Mux.Get("/news-api/ver1.0/", func(w http.ResponseWriter, r *http.Request) {
		sR, err := h.EH.GetAll(from, size)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		articleResponse(sR, w)
	})
}

func (h *Handler) Post(country, category string) {

}

func articleResponse(sR *elastic.SearchResult, w http.ResponseWriter) {
	articles := utils.Translator(sR)
		message, err := json.Marshal(articles)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		_, err = w.Write(message)
		if err != nil {
			return
		}
}
