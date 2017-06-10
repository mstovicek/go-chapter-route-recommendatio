package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/go-chapter-route-recommendation/api/handler"
	"github.com/mstovicek/go-chapter-route-recommendation/api/middleware"
	"github.com/mstovicek/go-chapter-route-recommendation/places_api"
	"github.com/mstovicek/go-chapter-route-recommendation/services"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

type config struct {
	Port string `env:"PORT,default=8080"`
}

type server struct {
	logger *log.Logger
	config config
}

func NewApiHttpServer(l *log.Logger) *server {
	var cnf config
	if err := envdecode.Decode(&cnf); err != nil {
		panic(err)
	}

	return &server{
		logger: l,
		config: cnf,
	}
}

func (api *server) Run() {
	r := mux.NewRouter()

	r.Handle(
		"/cities-suggestions/",
		handler.NewCitiesSuggestion(
			services.NewService(
				cache.New(5*time.Second, 10*time.Second),
				places_api.NewGoogleAPI(),
			),
		),
	).Methods(http.MethodGet)

	r.Handle(
		"/cities-info/",
		handler.NewCitiesInfo(
			services.NewService(
				cache.New(5*time.Second, 10*time.Second),
				places_api.NewGoogleAPI(),
			),
		),
	).Methods(http.MethodPost)

	http.Handle("/", r)

	router := middleware.NewRecovery(
		api.logger,
		middleware.NewLogResponseTime(
			api.logger,
			middleware.NewJsonHeader(
				r,
			),
		),
	)

	log.Fatal(http.ListenAndServe(":"+api.config.Port, router))
}
