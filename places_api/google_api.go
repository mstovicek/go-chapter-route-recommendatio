package places_api

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"googlemaps.github.io/maps"
	"time"
)

type config struct {
	GoogleAPIKey string `env:"API_KEY,required,strict"`
}

type GoogleAPI struct {
	cnf        config
	googleMaps *maps.Client
}

func NewGoogleAPI() *GoogleAPI {
	cnf := config{}

	if err := envdecode.Decode(&cnf); err != nil {
		log.Fatalln(err)
	}

	googleMaps, err := maps.NewClient(maps.WithAPIKey(cnf.GoogleAPIKey))
	if err != nil {
		log.Fatalln(err)
	}

	return &GoogleAPI{
		cnf:        cnf,
		googleMaps: googleMaps,
	}
}

func (googleAPI *GoogleAPI) GetPlaceDetail(placeID string) (*entity.Place, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"placeID": placeID,
	}).Info("Fetching place information")

	req := maps.PlaceDetailsRequest{
		PlaceID:  placeID,
		Language: "en",
	}

	res, err := googleAPI.googleMaps.PlaceDetails(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"placeID": placeID,
			"err":     err.Error(),
		}).Warn("Cannot fetch place detail")

		return nil, fmt.Errorf("Cannot get place information (place id = %s)", placeID)
	}

	log.WithFields(log.Fields{
		"placeID": placeID,
		"timeMs":  time.Since(start),
	}).Info("Fetched place information")

	p := entity.Place{
		PlaceID:          res.PlaceID,
		Name:             res.Name,
		FormattedAddress: res.FormattedAddress,
		Coordinates:      res.Geometry.Location,
	}
	return &p, nil
}

func (googleAPI *GoogleAPI) GetPlaceAutocompleteSuggestions(query string) ([]entity.Suggestion, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"query": query,
	}).Info("Fetching suggestions")

	req := maps.PlaceAutocompleteRequest{
		Input:    query,
		Types:    maps.AutocompletePlaceTypeCities,
		Language: "en",
	}

	res, err := googleAPI.googleMaps.PlaceAutocomplete(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"query": query,
			"err":   err.Error(),
		}).Warn("Suggestions request not successful")
		return nil, err
	}

	log.WithFields(log.Fields{
		"query":  query,
		"timeMs": time.Since(start),
	}).Info("Fetched suggestions")

	suggestionsCollection := make([]entity.Suggestion, len(res.Predictions))
	for i, p := range res.Predictions {
		suggestionsCollection[i] = entity.Suggestion{
			PlaceID:     p.PlaceID,
			Description: p.Description,
		}
	}

	return suggestionsCollection, nil
}
