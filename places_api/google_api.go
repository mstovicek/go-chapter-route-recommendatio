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

type googleAPI struct {
	cnf        config
	googleMaps *maps.Client
}

func NewGoogleAPI() *googleAPI {
	cnf := config{}

	if err := envdecode.Decode(&cnf); err != nil {
		log.Fatalln(err)
	}

	googleMaps, err := maps.NewClient(maps.WithAPIKey(cnf.GoogleAPIKey))
	if err != nil {
		log.Fatalln(err)
	}

	return &googleAPI{
		cnf:        cnf,
		googleMaps: googleMaps,
	}
}

func (api *googleAPI) GetPlaceDetail(placeID string) (*entity.Place, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"placeID": placeID,
	}).Info("Fetching place information")

	req := maps.PlaceDetailsRequest{
		PlaceID:  placeID,
		Language: "en",
	}

	res, err := api.googleMaps.PlaceDetails(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"placeID": placeID,
			"err":     err.Error(),
		}).Warn("Cannot fetch place detail")

		return nil, fmt.Errorf("Cannot get place information (place id = %s)", placeID)
	}

	log.WithFields(log.Fields{
		"placeID": placeID,
		"time":    time.Since(start),
	}).Info("Fetched place information")

	p := entity.Place{
		PlaceID:          res.PlaceID,
		Name:             res.Name,
		FormattedAddress: res.FormattedAddress,
		Coordinates:      res.Geometry.Location,
	}
	return &p, nil
}

func (api *googleAPI) GetPlaceAutocompleteSuggestions(query string) ([]entity.Suggestion, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"query": query,
	}).Info("Fetching suggestions")

	req := maps.PlaceAutocompleteRequest{
		Input:    query,
		Types:    maps.AutocompletePlaceTypeCities,
		Language: "en",
	}

	res, err := api.googleMaps.PlaceAutocomplete(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"query": query,
			"err":   err.Error(),
		}).Warn("Suggestions request not successful")
		return nil, err
	}

	log.WithFields(log.Fields{
		"query": query,
		"time":  time.Since(start),
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

func (api *googleAPI) GetPlacesDistance(placeIDs []string) (entity.DistanceMatrix, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"places": placeIDs,
	}).Info("Fetching distance matrix")

	placeIDsForRequest := make([]string, len(placeIDs))
	for i, placeID := range placeIDs {
		placeIDsForRequest[i] = "place_id:" + placeID
	}

	req := maps.DistanceMatrixRequest{
		Language:     "en",
		Units:        maps.UnitsMetric,
		Origins:      placeIDsForRequest,
		Destinations: placeIDsForRequest,
	}

	distanceMatrix := entity.NewDistanceMatrix()

	res, err := api.googleMaps.DistanceMatrix(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"places": placeIDs,
			"err":    err.Error(),
		}).Warn("Distance matrix request not successful")
		return nil, err
	}

	for i, row := range res.Rows {
		for j, elem := range row.Elements {
			distanceMatrix.Add(
				placeIDs[i],
				placeIDs[j],
				entity.Distance{
					FromPlaceID:     placeIDs[i],
					ToPlaceID:       placeIDs[j],
					DistanceMetres:  elem.Distance.Meters,
					DurationSeconds: elem.Duration.Seconds(),
				},
			)
		}
	}

	log.WithFields(log.Fields{
		"places": placeIDs,
		"time":   time.Since(start),
	}).Info("Fetched distance")

	return distanceMatrix, nil
}
