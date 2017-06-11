package handler

import (
	"encoding/json"
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"io/ioutil"
	"net/http"
)

type CityDistancesService interface {
	GetPlacesDistance(placeIDs []string) entity.DistanceMatrix
}

type CitiesDistances struct {
	citiesDistancesService CityDistancesService
}

func NewCitiesDistances(s CityDistancesService) *CitiesDistances {
	return &CitiesDistances{
		citiesDistancesService: s,
	}
}

func (handler *CitiesDistances) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		response.WriteError(w, 400, "invalid request format")
		return
	}

	dMap := handler.citiesDistancesService.GetPlacesDistance(in)

	var slice []entity.Distance

	for _, ddMap := range dMap {
		for _, dist := range ddMap {
			slice = append(slice, dist)
		}
	}

	response.WriteSuccess(w, slice)
}
