package controller

import (
	"dopc/internal/calculator"
	"dopc/internal/parser"
	"dopc/internal/venueapi"
	"encoding/json"
	"log"
	"net/http"
)

func DopcController(writer http.ResponseWriter, request *http.Request) {
	queries, pErr := parser.ParseRequest(request)
	if pErr != nil {
		log.Println(pErr.Message)
		sendError(writer, pErr.Message, pErr.Status)
		return
	}

	venueData, aErr := venueapi.ProcessVenue(queries.VenueSlug)
	if aErr != nil {
		log.Println(aErr.Debug)
		sendError(writer, aErr.Message, aErr.Status)
		return
	}

	summary, err := calculator.Placeholder(queries, venueData)
	if err != nil {
		log.Println(err)
		sendError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(summary)
	if err != nil {
		log.Println(err)
		sendError(writer, "", http.StatusInternalServerError)
	}
}

func sendError(writer http.ResponseWriter, message string, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	payload := map[string]string{}
	payload["error"] = message

	err := json.NewEncoder(writer).Encode(payload)
	if err != nil {
		log.Println(err)
		http.Error(writer, "", http.StatusInternalServerError)
	}
}
