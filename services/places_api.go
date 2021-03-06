package services

import (
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
)

type PlacesAPI interface {
	GetPlaceDetail(placeID string) (*entity.Place, error)
	GetPlaceAutocompleteSuggestions(query string) ([]entity.Suggestion, error)
	GetPlacesDistance(placeIDs []string) (entity.DistanceMatrix, error)
}
