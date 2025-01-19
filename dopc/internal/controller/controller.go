package controller

import (
	"dopc/internal/calculator"
	"dopc/internal/parser"
	"dopc/internal/venueapi"
	"log"
	"net/http"
)

func DopcController(writer http.ResponseWriter, request *http.Request) {
	var queries parser.Queries
	parseStatus := parser.ParseRequest(request, &queries)
	if parseStatus != nil {
		log.Println(parseStatus.Message)
		http.Error(writer, parseStatus.Message, parseStatus.Status)
		return
	}

	var venueData venueapi.Venue
	venueStatus := venueapi.ProcessVenue(queries.VenueSlug, &venueData)
	if venueStatus != nil {
		log.Println(venueStatus.Message)
		http.Error(writer, venueStatus.Message, venueStatus.Status)
		return
	}

	calculator.Placeholder(queries, venueData)
}
