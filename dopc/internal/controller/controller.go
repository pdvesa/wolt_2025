package controller

import (
	"dopc/internal/parser"
	"net/http"
)

func ApiController(writer http.ResponseWriter, request *http.Request) {
	var queries parser.Queries
	status := parser.ParseRequest(request, &queries)
	if status != nil {
		http.Error(writer, status.Message, status.Status)
		return
	}
}
