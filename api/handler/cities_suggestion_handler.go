package handler

import (
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"net/http"
)

type citySuggestionsCollectionService interface {
	GetPlacesSuggestionsByKeyword(keyword string) []entity.Suggestion
}

type citiesSuggestion struct {
	citySuggestionService citySuggestionsCollectionService
}

func NewCitiesSuggestion(s citySuggestionsCollectionService) *citiesSuggestion {
	return &citiesSuggestion{
		citySuggestionService: s,
	}
}

func (handler *citiesSuggestion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")

	if keyword == "" {
		response.WriteError(w, http.StatusBadRequest, "Query GET parameter (q) is missing")
		return
	}

	response.WriteSuccess(w, handler.citySuggestionService.GetPlacesSuggestionsByKeyword(keyword))
}
