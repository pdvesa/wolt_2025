package calculator

import (
	"dopc/internal/parser"
	"dopc/internal/venueapi"
	"fmt"
	"math"
)

const EarthRadius float64 = 6371.0

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	println(lat1, lon1, lat2, lon2)
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

func calculateDistance(clientLat float64, clientLon float64, venueCords []float64) error {
	dist := Haversine(clientLat, clientLon, venueCords[1], venueCords[0]) //lat and lon changed somewhere
	fmt.Printf("Distance: %.2f km\n", dist)
	return nil
}

func Placeholder(queries parser.Queries, venue venueapi.Venue) error {
	calculateDistance(queries.UserLat, queries.UserLon, venue.Location)

	return nil
}
