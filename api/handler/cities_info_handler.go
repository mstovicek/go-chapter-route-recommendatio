package handler

import (
	"encoding/json"
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"io/ioutil"
	"net/http"
)

type CityDetailsCollectionService interface {
	GetPlacesCollectionByPlaceIds(placeIds []string) []entity.Place
}

type CitiesInfo struct {
	cityDetailService CityDetailsCollectionService
}

func NewCitiesInfo(s CityDetailsCollectionService) *CitiesInfo {
	return &CitiesInfo{
		cityDetailService: s,
	}
}

func (handler *CitiesInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	response.WriteSuccess(w, handler.cityDetailService.GetPlacesCollectionByPlaceIds(in))
}
