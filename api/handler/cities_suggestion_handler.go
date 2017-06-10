package handler

import (
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"net/http"
)

type CitySuggestionsCollectionService interface {
	GetPlacesSuggestionsByKeyword(keyword string) []entity.Suggestion
}

type CitiesSuggestion struct {
	citySuggestionService CitySuggestionsCollectionService
}

func NewCitiesSuggestion(s CitySuggestionsCollectionService) *CitiesSuggestion {
	return &CitiesSuggestion{
		citySuggestionService: s,
	}
}

func (handler *CitiesSuggestion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")

	if keyword == "" {
		response.WriteError(w, http.StatusBadRequest, "Query GET parameter (q) is missing")
		return
	}

	response.WriteSuccess(w, handler.citySuggestionService.GetPlacesSuggestionsByKeyword(keyword))
}
