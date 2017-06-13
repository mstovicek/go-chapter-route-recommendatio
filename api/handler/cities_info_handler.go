package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mstovicek/go-chapter-route-recommendation/api/response"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"io/ioutil"
	"net/http"
)

type cityDetailsCollectionService interface {
	GetPlacesCollectionByPlaceIds(placeIds []string) ([]entity.Place, error)
}

type citiesInfo struct {
	cityDetailService cityDetailsCollectionService
}

func NewCitiesInfo(s cityDetailsCollectionService) *citiesInfo {
	return &citiesInfo{
		cityDetailService: s,
	}
}

func (handler *citiesInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	placesCollection, err := handler.cityDetailService.GetPlacesCollectionByPlaceIds(in)
	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("cannot fetch places, error: %s", err.Error()),
		)
		return
	}

	response.WriteSuccess(w, placesCollection)
}
