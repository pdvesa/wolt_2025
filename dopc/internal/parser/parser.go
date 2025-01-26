package parser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ParserError struct {
	Status  int
	Message string
}

type Queries struct {
	VenueSlug string
	CartValue int
	UserLat   float64
	UserLon   float64
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
	err := parseElement(queryList, "venue_slug", &queries.VenueSlug, true)
	if err != nil {
		return nil, err
	}

	err = parseElement(queryList, "cart_value", &queries.CartValue, true)
	if err != nil {
		return nil, err
	}

	err = parseElement(queryList, "user_lat", &queries.UserLat, true)
	if err != nil {
		return nil, err
	}

	err = parseElement(queryList, "user_lon", &queries.UserLon, true)
	if err != nil {
		return nil, err
	}

	err = extraChecks(&queries)
	if err != nil {
		return nil, err
	}

	return &queries, nil
}

func parseElement(queryList url.Values, queryName string, query interface{}, mandatory bool) *ParserError {
	conversionStr := queryList.Get(queryName)
	if conversionStr == "" && mandatory {
		return &ParserError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("missing mandatory query '%s'", queryName),
		}
	}

	switch queryType := query.(type) {
	case *string:
		*queryType = conversionStr
	case *int:
		value, err := strconv.Atoi(conversionStr)
		if err != nil {
			return &ParserError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid query in '%s'", queryName),
			}
		}
		*queryType = value

	case *float32:
		value, err := strconv.ParseFloat(conversionStr, 32)
		if err != nil {
			return &ParserError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid query in '%s'", queryName),
			}
		}
		*queryType = float32(value)

	case *float64:
		value, err := strconv.ParseFloat(conversionStr, 64)
		if err != nil {
			return &ParserError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid query in '%s'", queryName),
			}
		}
		*queryType = value

	default:
		log.Println("data conversion not implemented")
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
