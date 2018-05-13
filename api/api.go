package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tothmate90/news-scraper/config"
	"github.com/tothmate90/news-scraper/newsapi"

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
	Config         config.Config
}

func New(eH elasticsearch.Handler, conf config.Config) Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	return Handler{
		Mux:            mux,
		ElasticHandler: eH,
		Config:         conf,
	}
}

func (h *Handler) Listen() {
	http.ListenAndServe(":"+h.Config.Port, h.Mux)
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
		values := url.Values{}
		from, err := strconv.Atoi(values.Get("from"))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		size, err := strconv.Atoi(values.Get("size"))
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
		queryValues := url.Values{}
		values := url.Values{}
		values.Add("country", queryValues.Get("country"))
		values.Add("category", queryValues.Get("category"))

		result, err := newsapi.GetTopHeadlines(values, h.Config)
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
