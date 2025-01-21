package calculator

import (
	"dopc/internal/parser"
	"dopc/internal/venueapi"
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
type Order struct {
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

func calculateFee(ranges []venueapi.DistanceRange, baseFee int, delivery Delivery) int {
	for _, bracket := range ranges {
		if delivery.Distance >= bracket.Min && delivery.Distance < bracket.Max {

			fee := baseFee + bracket.A
			elementB := bracket.B * float64(delivery.Distance) / 10
			fee = fee + int(math.Round(elementB))
			delivery.Fee = fee
			fmt.Printf("Distance %d is in the range [%d, %d]\n", delivery.Distance, bracket.Min, bracket.Max)
			fmt.Println(delivery.Fee)
			return (0)
		} else if delivery.Distance > bracket.Min && bracket.Max == 0 {
			return (1)
		}
	}
	return (0)
}

func Placeholder(queries *parser.Queries, venue *venueapi.Venue) error {
	var Order Order
	Order.Delivery.Distance = calculateDistance(queries.UserLat, queries.UserLon, venue.Location)
	if queries.CartValue < venue.SurchargeMin {
		Order.SmallOrderSurcharge = venue.SurchargeMin - queries.CartValue
	}
	Order.Delivery.Distance = 600
	calculateFee(venue.DistanceRanges, venue.BasePrice, Order.Delivery)
	return nil
}
