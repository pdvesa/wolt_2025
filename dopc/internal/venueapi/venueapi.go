package venueapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqError struct {
	Status  int
	Message string
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

func getVenueData(url string, venue *Venue) *ReqError {
	var TempVenue TempVenue
	response, err := http.Get(url)
	if err != nil {
		return &ReqError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("API fetch failed: %v", err),
		}
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Venue: %d %s", response.StatusCode, http.StatusText(response.StatusCode)),
		}
	}

	err = json.NewDecoder(response.Body).Decode(&TempVenue)
	if err != nil {
		return &ReqError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("JSON decoding failed: %v", err),
		}
	}

	if venue.ID == "" {
		venue.ID = TempVenue.VenueRaw.ID
		println(TempVenue.VenueRaw.Location.Coordinates[0], TempVenue.VenueRaw.Location.Coordinates[1])
		venue.Location = TempVenue.VenueRaw.Location.Coordinates
		println(venue.Location[0], venue.Location[1]) // i dont like this
	}
	venue.SurchargeMin = TempVenue.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge
	venue.BasePrice = TempVenue.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice
	venue.DistanceRanges = TempVenue.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges //change maybe

	return nil
}

func ProcessVenue(venue string, venueData *Venue) *ReqError {
	apiUrl := "https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues/"

	url := apiUrl + venue + "/static"
	err := getVenueData(url, venueData)
	if err != nil {
		return err
	}

	url = apiUrl + venue + "/dynamic"
	err = getVenueData(url, venueData)
	if err != nil {
		return err
	}

	//debug
	println(venueData.ID)
	println(venueData.Location)
	println(venueData.SurchargeMin)
	println(venueData.BasePrice)
	println(venueData.DistanceRanges)

	return nil
}
