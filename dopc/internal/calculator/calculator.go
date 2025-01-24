package calculator

import (
	"dopc/internal/api"
	"dopc/internal/parser"
	"errors"
	"fmt"
	"math"
)

type Delivery struct {
	Fee      int `json:"fee"`
	Distance int `json:"distance"`
}

type OrderSummary struct {
	TotalPrice          int      `json:"total_price"`
	SmallOrderSurcharge int      `json:"small_order_surcharge"`
	CartValue           int      `json:"cart_value"`
	Delivery            Delivery `json:"delivery"`
}

const EarthRadius float64 = 6371000.0

func Calculator(queries *parser.Queries, venue *api.Venue) (*OrderSummary, error) {
	var summary OrderSummary
	var err error

	summary.CartValue = queries.CartValue
	summary.Delivery.Distance = calculateDistance(queries.UserLat, queries.UserLon, venue.Lat, venue.Lon)

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

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := EarthRadius * c

	return distance
} //stolen from internet

func calculateDistance(clientLat float64, clientLon float64, venueLat float64, venueLon float64) int {
	dist := haversine(clientLat, clientLon, venueLat, venueLon) //lat and lon changed somewhere
	result := int(math.Round(dist))
	return result
}

func calculateFee(ranges []api.DistanceRange, baseFee int, distance int) (int, error) {
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
