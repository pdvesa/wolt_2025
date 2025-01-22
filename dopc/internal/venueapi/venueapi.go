package venueapi

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

type TempVenue struct {
	VenueRaw VenueRaw `json:"venue_raw"`
}

type Venue struct {
	ID             string
	Location       []float64
	SurchargeMin   int
	BasePrice      int
	DistanceRanges []DistanceRange
}

// refactor structs

type dataFetcher interface {
	getVenueData(venueID string) *ApiError
}

func getVenueData(apiAddress string, venue *Venue) *ApiError {
	var TempVenue TempVenue
	response, err := http.Get(apiAddress)
	if err != nil {
		return &ApiError{
			Status:  http.StatusInternalServerError,
			Message: "",
			Debug:   fmt.Sprintf("api fetch failed: %v", err),
		}
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return &ApiError{
			Status:  http.StatusBadRequest,
			Message: "invalid 'venue_slug'",
			Debug:   fmt.Sprintf("Venue: %d %s", response.StatusCode, http.StatusText(response.StatusCode)),
		}
	} else if response.StatusCode != 200 {
		return &ApiError{
			Status:  http.StatusInternalServerError,
			Message: "",
			Debug:   fmt.Sprintf("venue: %d %s", response.StatusCode, http.StatusText(response.StatusCode)),
		}
	}

	err = json.NewDecoder(response.Body).Decode(&TempVenue)
	if err != nil {
		return &ApiError{
			Status:  http.StatusInternalServerError,
			Message: "",
			Debug:   fmt.Sprintf("json decoding failed: %v", err),
		}
	}

	if venue.ID == "" {
		venue.ID = TempVenue.VenueRaw.ID
		venue.Location = TempVenue.VenueRaw.Location.Coordinates
	}
	venue.SurchargeMin = TempVenue.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge
	venue.BasePrice = TempVenue.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice
	venue.DistanceRanges = TempVenue.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges //change maybe

	return nil
}

func ProcessVenue(venue string) (*Venue, *ApiError) {
	var venueData Venue
	apiUrl := "https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues/"

	url := apiUrl + venue + "/static"
	err := getVenueData(url, &venueData)
	if err != nil {
		return nil, err
	}

	url = apiUrl + venue + "/dynamic"
	err = getVenueData(url, &venueData)
	if err != nil {
		return nil, err
	}

	//debug
	println(venueData.ID)
	println(venueData.Location)
	println(venueData.SurchargeMin)
	println(venueData.BasePrice)
	println(venueData.DistanceRanges)

	return &venueData, nil
}
