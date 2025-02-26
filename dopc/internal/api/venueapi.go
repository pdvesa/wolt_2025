package api

/*

Package responsible for fetching and saving data from the home-assignment API

*/

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiError struct {
	Status  int
	Message string
	Debug   string
}

type DistanceRange struct {
	Min int     `json:"min"`
	Max int     `json:"max"`
	A   int     `json:"a"`
	B   float64 `json:"b"`
}

type DeliveryPricing struct {
	BasePrice      int             `json:"base_price"`
	DistanceRanges []DistanceRange `json:"distance_ranges"`
}

type DeliverySpecs struct {
	OrderMinimumNoSurcharge int             `json:"order_minimum_no_surcharge"`
	DeliveryPricing         DeliveryPricing `json:"delivery_pricing"`
}

type Location struct {
	Coordinates []float64 `json:"coordinates"`
}

type VenueRaw struct {
	ID            string        `json:"id"`
	Location      Location      `json:"location"`
	DeliverySpecs DeliverySpecs `json:"delivery_specs"`
}

type DecodeVenue struct {
	VenueRaw VenueRaw `json:"venue_raw"`
}

type Venue struct {
	ID             string
	Lat            float64
	Lon            float64
	SurchargeMin   int
	BasePrice      int
	DistanceRanges []DistanceRange
}

func ProcessVenue(venue string, apiAddress string) (*Venue, *ApiError) {
	venueData, err := getVenueData(venue, apiAddress)
	if err != nil {
		return nil, err
	}

	return venueData, nil
}

func getVenueData(venueID string, apiAddress string) (*Venue, *ApiError) {
	var tempVenue DecodeVenue
	var venueData Venue

	response, err := getApiResponse(venueID, apiAddress, "/static")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	err = decodeFields(response, &tempVenue)
	if err != nil {
		return nil, err
	}

	venueData.ID = tempVenue.VenueRaw.ID
	venueData.Lon = tempVenue.VenueRaw.Location.Coordinates[0]
	venueData.Lat = tempVenue.VenueRaw.Location.Coordinates[1]

	response, err = getApiResponse(venueID, apiAddress, "/dynamic")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	err = decodeFields(response, &tempVenue)
	if err != nil {
		return nil, err
	}

	venueData.SurchargeMin = tempVenue.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge
	venueData.BasePrice = tempVenue.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice
	venueData.DistanceRanges = tempVenue.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges

	return &venueData, nil
}

func getApiResponse(venue string, apiAddress string, endpoint string) (*http.Response, *ApiError) {
	url := apiAddress + venue + endpoint
	response, err := http.Get(url)

	if err != nil {
		return nil, &ApiError{
			Status:  http.StatusInternalServerError,
			Message: "",
			Debug:   fmt.Sprintf("api fetch failed: %v", err),
		}
	}

	if response.StatusCode == http.StatusNotFound {
		response.Body.Close()
		return nil, &ApiError{
			Status:  http.StatusBadRequest,
			Message: "invalid 'venue_slug'",
			Debug:   fmt.Sprintf("venue: %d %s", response.StatusCode, http.StatusText(response.StatusCode)),
		}
	} else if response.StatusCode != 200 {
		response.Body.Close()
		return nil, &ApiError{
			Status:  http.StatusInternalServerError,
			Message: "",
			Debug:   fmt.Sprintf("venue: %d %s", response.StatusCode, http.StatusText(response.StatusCode)),
		}
	}

	return response, nil
}

func decodeFields(response *http.Response, decodeStruct interface{}) *ApiError {
	err := json.NewDecoder(response.Body).Decode(decodeStruct)
	if err != nil {
		return &ApiError{
			Status:  http.StatusInternalServerError,
			Message: "",
			Debug:   fmt.Sprintf("json decoding failed: %v", err),
		}
	}

	return nil
}
