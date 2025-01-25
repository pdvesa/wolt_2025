package handler

import (
	"dopc/internal/api"
	"dopc/internal/calculator"
	"dopc/internal/parser"
	"encoding/json"
	"log"
	"net/http"
)

func DopcHandler(writer http.ResponseWriter, request *http.Request) {
	queries, pErr := parser.ParseRequest(request)
	if pErr != nil {
		log.Println(pErr.Message)
		sendError(writer, pErr.Message, pErr.Status)
		return
	}

	venueData, aErr := api.ProcessVenue(queries.VenueSlug)
	if aErr != nil {
		log.Println(aErr.Debug)
		sendError(writer, aErr.Message, aErr.Status)
		return
	}

	summary, err := calculator.Calculator(queries, venueData)
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
