package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"io/ioutil"
	"net/http"
)

type cityDistancesService interface {
	GetPlacesDistance(placeIDs []string) (entity.DistanceMatrix, error)
}

type citiesDistances struct {
	citiesDistancesService cityDistancesService
}

func NewCitiesDistances(s cityDistancesService) *citiesDistances {
	return &citiesDistances{
		citiesDistancesService: s,
	}
}

func (handler *citiesDistances) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if nil == r.Body {
		response.WriteError(w, http.StatusBadRequest, "body is empty")
		return
	}

	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "cannot read body")
		return
	}

	var in []string
	if err := json.Unmarshal(data, &in); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	distancesMap, err := handler.citiesDistancesService.GetPlacesDistance(in)
	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("cannot fetch distances, error: %s", err.Error()),
		)
		return
	}

	var slice []entity.Distance

	for _, ddMap := range distancesMap {
		for _, dist := range ddMap {
			slice = append(slice, dist)
		}
	}

	response.WriteSuccess(w, slice)
}
