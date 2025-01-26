package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessVenue_Success(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/home-assignment-api/v1/venues/test-venue/static", func(w http.ResponseWriter, r *http.Request) {
		venue := DecodeVenue{
			VenueRaw: VenueRaw{
				ID: "test-venue",
				Location: Location{
					Coordinates: []float64{12.3456, 65.4321},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(venue)
	})
	handler.HandleFunc("/home-assignment-api/v1/venues/test-venue/dynamic", func(w http.ResponseWriter, r *http.Request) {
		venue := DecodeVenue{
			VenueRaw: VenueRaw{
				DeliverySpecs: DeliverySpecs{
					OrderMinimumNoSurcharge: 50,
					DeliveryPricing: DeliveryPricing{
						BasePrice: 100,
						DistanceRanges: []DistanceRange{
							{Min: 0, Max: 5, A: 10, B: 1.5},
							{Min: 6, Max: 10, A: 15, B: 2.0},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(venue)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	apiAddress := server.URL + "/home-assignment-api/v1/venues/"

	venue, aErr := ProcessVenue("test-venue", apiAddress)
	if aErr != nil {
		t.Fatalf("expected no error, got %v", aErr.Debug)
	}

	if venue.ID != "test-venue" {
		t.Errorf("expected venue ID 'test-venue', got '%s'", venue.ID)
	}
	if venue.Lat != 65.4321 {
		t.Errorf("expected latitude 65.4321, got %f", venue.Lat)
	}
	if venue.Lon != 12.3456 {
		t.Errorf("expected longitude 12.3456, got %f", venue.Lon)
	}
	if venue.SurchargeMin != 50 {
		t.Errorf("expected surcharge minimum 50, got %d", venue.SurchargeMin)
	}
	if venue.BasePrice != 100 {
		t.Errorf("expected base price 100, got %d", venue.BasePrice)
	}
	if len(venue.DistanceRanges) != 2 {
		t.Errorf("expected 2 distance ranges, got %d", len(venue.DistanceRanges))
	} else {
		if venue.DistanceRanges[0].Min != 0 || venue.DistanceRanges[0].Max != 5 {
			t.Errorf("expected distance range [{0, 5}], got %v", venue.DistanceRanges[0])
		}
		if venue.DistanceRanges[1].Min != 6 || venue.DistanceRanges[1].Max != 10 {
			t.Errorf("expected distance range [{6, 10}], got %v", venue.DistanceRanges[1])
		}
	}
}
