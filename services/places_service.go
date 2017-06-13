package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/mstovicek/go-chapter-route-recommendation/entity"
	"sort"
	"strings"
	"sync"
	"time"
)

const cachePrefixPlace = "place_"
const cachePrefixSuggestion = "suggestion_"
const cachePrefixDistances = "distances_"

type placesService struct {
	cache  Cache
	api    PlacesAPI
	logger *logrus.Logger
}

func NewPlacesService(cache Cache, api PlacesAPI, logger *logrus.Logger) *placesService {
	return &placesService{
		cache:  cache,
		api:    api,
		logger: logger,
	}
}

func (service *placesService) GetPlacesCollectionByPlaceIds(placeIDs []string) ([]entity.Place, error) {
	start := time.Now()

	service.logger.WithFields(logrus.Fields{
		"places": placeIDs,
	}).Info("Fetching places information")

	var waitGroup sync.WaitGroup

	placesChan := make(chan *entity.Place)

	placesCollection := []entity.Place{}

	for _, placeID := range placeIDs {
		waitGroup.Add(1)
		go func(placeID string) {
			place := service.getPlace(placeID)
			placesChan <- place
		}(placeID)
	}

	go func() {
		for place := range placesChan {
			if place != nil {
				placesCollection = append(placesCollection, *place)
			}
			waitGroup.Done()
		}
	}()

	waitGroup.Wait()

	service.logger.WithFields(logrus.Fields{
		"places": placeIDs,
		"time":   time.Since(start),
	}).Info("Fetched places information")

	return placesCollection, nil
}

func (service *placesService) getPlace(placeID string) *entity.Place {
	cacheKey := cachePrefixPlace + placeID

	cachedPlace, found := service.cache.Get(cacheKey)
	if found {
		service.logger.WithFields(logrus.Fields{
			"placeID": placeID,
		}).Info("Returning cached place")

		return cachedPlace.(*entity.Place)
	}

	start := time.Now()

	service.logger.WithFields(logrus.Fields{
		"placeID": placeID,
	}).Info("Fetching place information")

	p, err := service.api.GetPlaceDetail(placeID)
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"placeID": placeID,
			"err":     err.Error(),
		}).Error("Cannot get place")
		return nil
	}

	service.logger.WithFields(logrus.Fields{
		"placeID": placeID,
		"time":    time.Since(start),
	}).Info("Fetched place information")

	service.cache.Set(cacheKey, p, 0)

	return p
}

func (service *placesService) GetPlacesSuggestionsByKeyword(keyword string) ([]entity.Suggestion, error) {
	cacheKey := cachePrefixSuggestion + keyword

	cachedSuggestions, found := service.cache.Get(cacheKey)
	if found {
		service.logger.WithFields(logrus.Fields{
			"keyword": keyword,
		}).Info("Returning cached suggestions")

		return cachedSuggestions.([]entity.Suggestion), nil
	}

	start := time.Now()

	service.logger.WithFields(logrus.Fields{
		"keyword": keyword,
	}).Info("Fetching suggestions")

	suggestions, err := service.api.GetPlaceAutocompleteSuggestions(keyword)
	if err.Error() != "maps: ZERO_RESULTS -" {
		service.logger.WithFields(logrus.Fields{
			"keyword": keyword,
			"err":     err.Error(),
		}).Info("No suggestions found")
		return []entity.Suggestion{}, nil
	}
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"keyword": keyword,
			"err":     err.Error(),
		}).Error("Cannot get suggestions")
		return []entity.Suggestion{}, err
	}

	service.logger.WithFields(logrus.Fields{
		"keyword": keyword,
		"time":    time.Since(start),
	}).Info("Fetched suggestions")

	service.cache.Set(cacheKey, suggestions, 0)

	return suggestions, nil
}

func (service *placesService) GetPlacesDistance(placesIDs []string) (entity.DistanceMatrix, error) {
	sort.Strings(placesIDs)
	cacheKey := cachePrefixDistances + strings.Join(placesIDs, "|")

	cachedDistanceMatrix, found := service.cache.Get(cacheKey)
	if found {
		service.logger.WithFields(logrus.Fields{
			"places": cacheKey,
		}).Info("Returning cached distance matrix")

		return cachedDistanceMatrix.(entity.DistanceMatrix), nil
	}

	start := time.Now()

	service.logger.WithFields(logrus.Fields{
		"places": placesIDs,
	}).Info("Fetching distance matrix")

	distanceMatrix, err := service.api.GetPlacesDistance(placesIDs)
	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"places": placesIDs,
			"err":    err.Error(),
		}).Error("Cannot get distances")
		return entity.NewDistanceMatrix(), err
	}

	service.logger.WithFields(logrus.Fields{
		"places": placesIDs,
		"time":   time.Since(start),
	}).Info("Fetched distance")

	service.cache.Set(cacheKey, distanceMatrix, 0)

	return distanceMatrix, nil
}
