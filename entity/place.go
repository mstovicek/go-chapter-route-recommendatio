package entity

type Place struct {
	PlaceID          string `json:"place_id"`
	Name             string `json:"name"`
	FormattedAddress string `json:"formatted_address"`
	Coordinates      struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinates"`
}
