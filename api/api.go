package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/mstovicek/go-chapter-route-recommendation/api/handler"
	"github.com/mstovicek/go-chapter-route-recommendation/api/middleware"
	"github.com/mstovicek/go-chapter-route-recommendation/places_api"
	"github.com/mstovicek/go-chapter-route-recommendation/services"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

type config struct {
	Listen string `long:"listen" description:"Listen address" default:":8080" required:"true"`
}

type server struct {
	logger *log.Logger
	config config
}

func NewApiHttpServer(l *log.Logger) *server {
	var cnf config

	if _, err := flags.NewParser(&cnf, flags.HelpFlag|flags.PassDoubleDash).Parse(); err != nil {
		log.Fatalln(err)
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
			services.NewPlacesService(
				cache.New(5*time.Second, 10*time.Second),
				places_api.NewGoogleAPI(),
			),
		),
	).Methods(http.MethodGet)

	r.Handle(
		"/cities-info/",
		handler.NewCitiesInfo(
			services.NewPlacesService(
				cache.New(5*time.Second, 10*time.Second),
				places_api.NewGoogleAPI(),
			),
		),
	).Methods(http.MethodPost)

	r.Handle(
		"/distances-matrix/",
		handler.NewCitiesDistances(
			services.NewPlacesService(
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

	log.Infof("Listening on the address %s", api.config.Listen)
	log.Fatal(http.ListenAndServe(api.config.Listen, router))
}
