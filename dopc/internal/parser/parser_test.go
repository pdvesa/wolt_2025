package parser

import (
	"net/http"
	"net/url"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		queryParams   url.Values
		expectedQuery *Queries
		expectedError *ParserError
	}{
		{
			name:   "Valid request",
			method: http.MethodGet,
			queryParams: url.Values{
				"venue_slug": {"test_venue"},
				"cart_value": {"100"},
				"user_lat":   {"40.7128"},
				"user_lon":   {"-74.0060"},
			},
			expectedQuery: &Queries{
				VenueSlug: "test_venue",
				CartValue: 100,
				UserLat:   40.7128,
				UserLon:   -74.0060,
			},
			expectedError: nil,
		},
		{
			name:   "Invalid HTTP method",
			method: http.MethodPost,
			queryParams: url.Values{
				"venue_slug": {"test_venue"},
				"cart_value": {"100"},
				"user_lat":   {"40.7128"},
				"user_lon":   {"-74.0060"},
			},
			expectedQuery: nil,
			expectedError: &ParserError{
				Status:  http.StatusMethodNotAllowed,
				Message: "",
			},
		},
		{
			name:   "Missing mandatory query 'venue_slug'",
			method: http.MethodGet,
			queryParams: url.Values{
				"cart_value": {"100"},
				"user_lat":   {"40.7128"},
				"user_lon":   {"-74.0060"},
			},
			expectedQuery: nil,
			expectedError: &ParserError{
				Status:  http.StatusBadRequest,
				Message: "missing mandatory query 'venue_slug'",
			},
		},
		{
			name:   "Invalid cart_value",
			method: http.MethodGet,
			queryParams: url.Values{
				"venue_slug": {"test_venue"},
				"cart_value": {"invalid"},
				"user_lat":   {"40.7128"},
				"user_lon":   {"-74.0060"},
			},
			expectedQuery: nil,
			expectedError: &ParserError{
				Status:  http.StatusBadRequest,
				Message: "invalid query in 'cart_value'",
			},
		},
		{
			name:   "Invalid user_lat (out of bounds)",
			method: http.MethodGet,
			queryParams: url.Values{
				"venue_slug": {"test_venue"},
				"cart_value": {"100"},
				"user_lat":   {"100.0"},
				"user_lon":   {"-74.0060"},
			},
			expectedQuery: nil,
			expectedError: &ParserError{
				Status:  http.StatusBadRequest,
				Message: "Coordinates not from Earth",
			},
		},
		{
			name:   "Invalid user_lon (out of bounds)",
			method: http.MethodGet,
			queryParams: url.Values{
				"venue_slug": {"test_venue"},
				"cart_value": {"100"},
				"user_lat":   {"40.7128"},
				"user_lon":   {"200.0"},
			},
			expectedQuery: nil,
			expectedError: &ParserError{
				Status:  http.StatusBadRequest,
				Message: "Coordinates not from Earth",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "http://example.com", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.URL.RawQuery = tt.queryParams.Encode()
			queries, pErr := ParseRequest(req)

			if pErr != nil && (pErr.Status != tt.expectedError.Status || pErr.Message != tt.expectedError.Message) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if queries != nil && *queries != *tt.expectedQuery {
				t.Errorf("expected queries %v, got %v", tt.expectedQuery, queries)
			}
		})
	}
}
