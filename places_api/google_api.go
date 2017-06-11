package places_api

import (
	"context"
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"googlemaps.github.io/maps"
)

type config struct {
	GoogleAPIKey string `env:"API_KEY,required,strict"`
}

type googleAPI struct {
	cnf        config
	googleMaps *maps.Client
}

const languageEn = "en"

func NewGoogleAPI() (*googleAPI, error) {
	cnf := config{}

	if err := envdecode.Decode(&cnf); err != nil {
		return nil, err
	}

	googleMaps, err := maps.NewClient(maps.WithAPIKey(cnf.GoogleAPIKey))
	if err != nil {
		return nil, err
	}

	return &googleAPI{
		cnf:        cnf,
		googleMaps: googleMaps,
	}, nil
}

func (api *googleAPI) GetPlaceDetail(placeID string) (*entity.Place, error) {
	req := maps.PlaceDetailsRequest{
		PlaceID:  placeID,
		Language: languageEn,
	}

	res, err := api.googleMaps.PlaceDetails(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	place := entity.Place{
		PlaceID:          res.PlaceID,
		Name:             res.Name,
		FormattedAddress: res.FormattedAddress,
		Coordinates:      res.Geometry.Location,
	}
	return &place, nil
}

func (api *googleAPI) GetPlaceAutocompleteSuggestions(query string) ([]entity.Suggestion, error) {
	req := maps.PlaceAutocompleteRequest{
		Input:    query,
		Types:    maps.AutocompletePlaceTypeCities,
		Language: languageEn,
	}

	res, err := api.googleMaps.PlaceAutocomplete(context.Background(), &req)
	if err != nil {
		return nil, err
	}

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
	placeIDsForRequest := make([]string, len(placeIDs))
	for i, placeID := range placeIDs {
		placeIDsForRequest[i] = "place_id:" + placeID
	}

	req := maps.DistanceMatrixRequest{
		Language:     languageEn,
		Units:        maps.UnitsMetric,
		Origins:      placeIDsForRequest,
		Destinations: placeIDsForRequest,
	}

	distanceMatrix := entity.NewDistanceMatrix()

	res, err := api.googleMaps.DistanceMatrix(context.Background(), &req)
	if err != nil {
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

	return distanceMatrix, nil
}
