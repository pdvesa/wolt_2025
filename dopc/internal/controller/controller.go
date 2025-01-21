package controller

import (
	"dopc/internal/calculator"
	"dopc/internal/parser"
	"dopc/internal/venueapi"
	"log"
	"net/http"
)

func DopcController(writer http.ResponseWriter, request *http.Request) {
	queries, pErr := parser.ParseRequest(request)
	if pErr != nil {
		log.Println(pErr.Message)
		http.Error(writer, pErr.Message, pErr.Status)
		return
	}

	venueData, aErr := venueapi.ProcessVenue(queries.VenueSlug)
	if aErr != nil {
		log.Println(aErr.Debug)
		http.Error(writer, aErr.Message, aErr.Status)
		return
	}

	calculator.Placeholder(queries, venueData)
}
