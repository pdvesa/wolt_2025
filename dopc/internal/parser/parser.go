package parser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ReqError struct {
	Status  int
	Message string
}

type Queries struct {
	VenueSlug string
	CartValue int
	UserLat   float64
	UserLon   float64
}

func parseNumericElement(queryList url.Values, queryName string, query interface{}) *ReqError {
	conversionStr := queryList.Get(queryName)
	if conversionStr == "" {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Missing mandatory query '%s'", queryName),
		}
	}

	switch target := query.(type) {
	case *int:
		value, err := strconv.Atoi(conversionStr)
		if err != nil {
			return &ReqError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid query in '%s': %v", queryName, err),
			}
		}
		*target = value

	case *float32:
		value, err := strconv.ParseFloat(conversionStr, 32)
		if err != nil {
			return &ReqError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid query in '%s': %v", queryName, err),
			}
		}
		*target = float32(value)

	case *float64:
		value, err := strconv.ParseFloat(conversionStr, 64)
		if err != nil {
			return &ReqError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid query in '%s': %v", queryName, err),
			} //maybe print to terminal for debug
		}
		*target = value

	default:
		log.Println("Error: User error :)")
	}
	return nil
}

func extraChecks(queries *Queries) *ReqError {
	if queries.CartValue < 0 {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: "We can't pay for your order",
		}
	}

	if queries.UserLat <= -90 || queries.UserLat >= 90 {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: "Not from this planet",
		}
	}

	if queries.UserLon <= -180 || queries.UserLon >= 180 {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: "Not from this planet",
		}
	}

	return nil
}

func ParseRequest(request *http.Request, queries *Queries) *ReqError {
	if request.Method != http.MethodGet {
		return &ReqError{
			Status:  http.StatusMethodNotAllowed,
			Message: "",
		}
	}

	queryList := request.URL.Query()
	if len(queryList) == 0 {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: "Missing mandatory queries",
		}
	}
	/*	// debug
		for key, values := range queryList {
			fmt.Println("Key:", key)
			for _, value := range values {
				fmt.Println("Value:", value)
			}
		}
		//debug*/
	queries.VenueSlug = queryList.Get("venue_slug")
	if queries.VenueSlug == "" {
		return &ReqError{
			Status:  http.StatusBadRequest,
			Message: "Missing mandatory query 'venue_slug'",
		}
	}

	err := parseNumericElement(queryList, "cart_value", &queries.CartValue)
	if err != nil {
		return err
	}

	err = parseNumericElement(queryList, "user_lat", &queries.UserLat)
	if err != nil {
		return err
	}

	err = parseNumericElement(queryList, "user_lon", &queries.UserLon)
	if err != nil {
		return err
	}

	err = extraChecks(queries)
	if err != nil {
		return err
	}

	return nil
}
