package services

import (
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
)

type API interface {
	GetPlaceDetail(placeID string) (*entity.Place, error)
	GetPlaceAutocompleteSuggestions(query string) ([]entity.Suggestion, error)
}
