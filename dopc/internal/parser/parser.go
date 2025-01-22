package parser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Parser interface {
	ParseRequest(request *http.Request) *ParserError
}

type ParserError struct {
	Status  int
	Message string
	Debug   string
}

type Queries struct {
	VenueSlug string
	CartValue int
	UserLat   float64
	UserLon   float64
}

func parseNumericElement(queryList url.Values, queryName string, query interface{}) *ParserError {
	conversionStr := queryList.Get(queryName)
	if conversionStr == "" {
		return &ParserError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("missing mandatory query '%s'", queryName),
		}
	}

	switch target := query.(type) {
	case *int:
		value, err := strconv.Atoi(conversionStr)
		if err != nil {
			return &ParserError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid query in '%s': %v", queryName, err),
			}
		}
		*target = value

	case *float32:
		value, err := strconv.ParseFloat(conversionStr, 32)
		if err != nil {
			return &ParserError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid query in '%s': %v", queryName, err),
			}
		}
		*target = float32(value)

	case *float64:
		value, err := strconv.ParseFloat(conversionStr, 64)
		if err != nil {
			return &ParserError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid query in '%s': %v", queryName, err),
			}
		}
		*target = value

	default:
		log.Println("error: User error :)")
	}
	return nil
}

func extraChecks(queries *Queries) *ParserError {
	if queries.CartValue < 0 {
		return &ParserError{
			Status:  http.StatusBadRequest,
			Message: "query 'cart_value' negative",
		}
	}

	if queries.UserLat <= -90 || queries.UserLat >= 90 {
		return &ParserError{
			Status:  http.StatusBadRequest,
			Message: "Coordinates not from Earth",
		}
	}

	if queries.UserLon <= -180 || queries.UserLon >= 180 {
		return &ParserError{
			Status:  http.StatusBadRequest,
			Message: "Coordinates not from Earth",
		}
	}

	return nil
}

func ParseRequest(request *http.Request) (*Queries, *ParserError) {
	var queries Queries
	if request.Method != http.MethodGet {
		return nil, &ParserError{
			Status:  http.StatusMethodNotAllowed,
			Message: "",
		}
	}

	queryList := request.URL.Query()
	if len(queryList) == 0 {
		return nil, &ParserError{
			Status:  http.StatusBadRequest,
			Message: "missing mandatory queries",
		}
	}

	queries.VenueSlug = queryList.Get("venue_slug")
	if queries.VenueSlug == "" {
		return nil, &ParserError{
			Status:  http.StatusBadRequest,
			Message: "missing mandatory query 'venue_slug'",
		}
	}

	err := parseNumericElement(queryList, "cart_value", &queries.CartValue)
	if err != nil {
		return nil, err
	}

	err = parseNumericElement(queryList, "user_lat", &queries.UserLat)
	if err != nil {
		return nil, err
	}

	err = parseNumericElement(queryList, "user_lon", &queries.UserLon)
	if err != nil {
		return nil, err
	}

	err = extraChecks(&queries)
	if err != nil {
		return nil, err
	}

	return &queries, nil
}
