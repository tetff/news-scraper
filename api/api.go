package api

import (
	"github.com/tothmate90/news-scraper/newsapi"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/olivere/elastic"
	"github.com/tothmate90/news-scraper/utils"

	"github.com/tothmate90/news-scraper/elasticsearch"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	Mux            *chi.Mux
	ElasticHandler elasticsearch.Handler
	values         url.Values
}

func New(port string, eH elasticsearch.Handler) Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	http.ListenAndServe(":"+port, mux)
	return Handler{
		Mux:            mux,
		ElasticHandler: eH,
	}
}

func (h *Handler) Get() {
	h.Mux.Get("/news-api/v1/{articleID}", func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		sR, err := h.ElasticHandler.Get(articleID)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		articleResponse(sR, w)
	})
}

func (h *Handler) GetAll() {
	h.Mux.Get("/news-api/v1/", func(w http.ResponseWriter, r *http.Request) {
		from, err := strconv.Atoi(h.values.Get("from"))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		size, err := strconv.Atoi(h.values.Get("size"))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		sR, err := h.ElasticHandler.GetAll(from, size)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		articleResponse(sR, w)
	})
}

func (h *Handler) Post() {
	h.Mux.Post("/news-api/v1/", func(w http.ResponseWriter, r *http.Request) {
		values := url.Values{}
		values.Add("country", h.values.Get("country"))
		values.Add("category", h.values.Get("category"))

		result, err := newsapi.GetTopHeadlines(values)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		err = h.ElasticHandler.Post(result.Articles)
		if err != nil {
			w.WriteHeader(400)
			return			
		}
		message, err := json.Marshal(result.Articles)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		_, err = w.Write(message)
		if err != nil {
			w.WriteHeader(400)
			return			
		}
	})
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
		w.WriteHeader(400)
		return
	}
}
