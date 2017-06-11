package api

import (
	"github.com/Sirupsen/logrus"
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
	logger *logrus.Logger
	config config
}

func NewApiHttpServer(logger *logrus.Logger) *server {
	var cnf config

	if _, err := flags.NewParser(&cnf, flags.HelpFlag|flags.PassDoubleDash).Parse(); err != nil {
		logger.Fatalln(err)
	}

	return &server{
		logger: logger,
		config: cnf,
	}
}

func (api *server) Run() {
	router := mux.NewRouter()

	googlePlacesAPI, err := places_api.NewGoogleAPI()
	if err != nil {
		api.logger.Fatalln(err)
	}

	router.Handle(
		"/cities-suggestions/",
		handler.NewCitiesSuggestion(
			services.NewPlacesService(
				cache.New(5*time.Second, 10*time.Second),
				googlePlacesAPI,
				api.logger,
			),
		),
	).Methods(http.MethodGet)

	router.Handle(
		"/cities-info/",
		handler.NewCitiesInfo(
			services.NewPlacesService(
				cache.New(5*time.Second, 10*time.Second),
				googlePlacesAPI,
				api.logger,
			),
		),
	).Methods(http.MethodPost)

	router.Handle(
		"/distances-matrix/",
		handler.NewCitiesDistances(
			services.NewPlacesService(
				cache.New(5*time.Second, 10*time.Second),
				googlePlacesAPI,
				api.logger,
			),
		),
	).Methods(http.MethodPost)

	http.Handle("/", router)

	handler := middleware.NewRecovery(
		api.logger,
		middleware.NewLogResponseTime(
			api.logger,
			middleware.NewJsonHeader(
				router,
			),
		),
	)

	api.logger.Infof("Listening on the address %s", api.config.Listen)
	api.logger.Fatal(http.ListenAndServe(api.config.Listen, handler))
}
