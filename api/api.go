package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tothmate90/news-scraper/config"
	"github.com/tothmate90/news-scraper/newsapi"

	"github.com/olivere/elastic"
	"github.com/tothmate90/news-scraper/utils"

	"github.com/tothmate90/news-scraper/elasticsearch"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Handler API handler containing the mux from chi's library, the handler managing elasticsearch and the config library.
type Handler struct {
	Mux            *chi.Mux
	ElasticHandler elasticsearch.Handler
	Config         config.Config
}

// New Responsible for setting up the API.
func New(eH elasticsearch.Handler, conf config.Config) Handler {
	mux := chi.NewRouter()
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	return Handler{
		Mux:            mux,
		ElasticHandler: eH,
		Config:         conf,
	}
}

// Listen Opening the port read from the config file.
func (h *Handler) Listen() {
	h.setup()
	http.ListenAndServe(h.Config.Host, h.Mux)
}

func (h *Handler) setup() {
	h.Get()
	log.Info("Get method succesfully set up.")
	h.GetAll()
	log.Info("GetAll method succesfully set up.")
	h.Post()
	log.Info("Post method succesfully set up.")
}

// Get API get method that requires the ID as a param of the article in ES.
func (h *Handler) Get() {
	h.Mux.Get("/news-api/v1/{articleID}", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Get entered.")
		articleID := chi.URLParam(r, "articleID")
		sR, err := h.ElasticHandler.Get(articleID)
		log.Info("Got article.")
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		articleResponse(sR, w)
		return
	})
}

// GetAll API get method that returns all articles in a range. The from and size of the result is required as a query.
func (h *Handler) GetAll() {
	h.Mux.Get("/news-api/v1", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		from, err := strconv.Atoi(values.Get("from"))
		log.Info("From")
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		size, err := strconv.Atoi(values.Get("size"))
		log.Info("Size")
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		sR, err := h.ElasticHandler.GetAll(from, size)
		log.Info("Got entries.")
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		articleResponse(sR, w)
		return
	})
}

// Post API post method that will insert the result from the newsapi.org into ES. Country and category are required as a query.
func (h *Handler) Post() {
	h.Mux.Post("/news-api/v1", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		values := url.Values{}
		values.Add("country", queryValues.Get("country"))
		values.Add("category", queryValues.Get("category"))
		log.Info("Values handled.")
		result, err := newsapi.GetTopHeadlines(values, h.Config)
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		log.Info("Newsapi.org visited.")
		err = h.ElasticHandler.Post(result.Articles)
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		log.Info("Articles posted.")
		message, err := json.Marshal(result.Articles)
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
			return
		}
		_, err = w.Write(message)
		if err != nil {
			w.WriteHeader(400)
			log.Error(err)
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
