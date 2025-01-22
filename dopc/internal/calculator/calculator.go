package calculator

import (
	"dopc/internal/parser"
	"dopc/internal/venueapi"
	"errors"
	"fmt"
	"math"
)

const EarthRadius float64 = 6371000.0

// temp place for structs
type Delivery struct {
	Fee      int `json:"fee"`
	Distance int `json:"distance"`
}

// Struct for the main data
type OrderSummary struct {
	TotalPrice          int      `json:"total_price"`
	SmallOrderSurcharge int      `json:"small_order_surcharge"`
	CartValue           int      `json:"cart_value"`
	Delivery            Delivery `json:"delivery"`
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Differences in coordinates
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the distance
	distance := EarthRadius * c

	return distance
} //stolen from internet

func calculateDistance(clientLat float64, clientLon float64, venueCords []float64) int {
	dist := Haversine(clientLat, clientLon, venueCords[1], venueCords[0]) //lat and lon changed somewhere
	result := int(math.Round(dist))
	fmt.Println(result) //debug
	return result
}

func calculateFee(ranges []venueapi.DistanceRange, baseFee int, distance int) (int, error) {
	for _, bracket := range ranges {
		if distance >= bracket.Min && distance < bracket.Max {

			fee := baseFee + bracket.A
			elementB := bracket.B * float64(distance) / 10
			fee = fee + int(math.Round(elementB))

			fmt.Printf("Distance %d is in the range [%d, %d]\n", distance, bracket.Min, bracket.Max)
			fmt.Println(fee)

			return fee, nil
		} else if distance > bracket.Min && bracket.Max == 0 {
			return 0, errors.New("can't deliver to location")
		}
	}
	return 0, errors.New("something went wrong")
}

func Placeholder(queries *parser.Queries, venue *venueapi.Venue) (*OrderSummary, error) {
	var summary OrderSummary
	var err error
	summary.CartValue = queries.CartValue
	summary.Delivery.Distance = calculateDistance(queries.UserLat, queries.UserLon, venue.Location)
	summary.Delivery.Fee, err = calculateFee(venue.DistanceRanges, venue.BasePrice, summary.Delivery.Distance)
	if err != nil {
		return nil, err
	}
	if summary.CartValue < venue.SurchargeMin {
		summary.SmallOrderSurcharge = venue.SurchargeMin - summary.CartValue
	}
	summary.TotalPrice = summary.CartValue + summary.SmallOrderSurcharge + summary.Delivery.Fee
	return &summary, nil
}
