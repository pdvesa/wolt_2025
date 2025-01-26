package calculator

import (
	"dopc/internal/api"
	"dopc/internal/parser"
	"testing"
)

type mockDistanceCalculator struct {
	mockDistance int
}

func (m *mockDistanceCalculator) CalculateDistance(clientLat, clientLon, venueLat, venueLon float64) int {
	return m.mockDistance
}
func TestCalculator(t *testing.T) {
	venue := &api.Venue{
		ID:           "1",
		Lat:          0,
		Lon:          0,
		SurchargeMin: 50,
		BasePrice:    10,
		DistanceRanges: []api.DistanceRange{
			{
				Min: 0,
				Max: 1000,
				A:   5,
				B:   0,
			},
			{
				Min: 1000,
				Max: 2000,
				A:   10,
				B:   2,
			},
			{
				Min: 2000,
				Max: 0,
				A:   0,
				B:   0,
			},
		},
	}

	queries := &parser.Queries{
		VenueSlug: "test-venue",
		CartValue: 0,
		UserLat:   0,
		UserLon:   0,
	}

	testCases := []struct {
		name                string
		mockDistance        int
		cartValue           int
		expectedTotalPrice  int
		expectedSurcharge   int
		expectedDeliveryFee int
		expectedErrorMsg    string
	}{
		{
			name:                "Distance 800 meters, Cart value 40",
			mockDistance:        800,
			cartValue:           40,
			expectedTotalPrice:  40 + 10 + 10 + 5 + int(0*float64(800)/10), // cart_value + surcharge + base_price + distance_a + distance_b
			expectedSurcharge:   10,
			expectedDeliveryFee: 10 + 5 + int(0*float64(800)/10),
			expectedErrorMsg:    "",
		},
		{
			name:                "Distance 800 meters, Cart value 60",
			mockDistance:        800,
			cartValue:           60,
			expectedTotalPrice:  60 + 0 + 10 + 5 + int(0*float64(800)/10),
			expectedSurcharge:   0,
			expectedDeliveryFee: 10 + 5 + int(0*float64(800)/10),
			expectedErrorMsg:    "",
		},
		{
			name:                "Distance 1400 meters, Cart value 50",
			mockDistance:        1400,
			cartValue:           50,
			expectedTotalPrice:  50 + 0 + 10 + 10 + int(2*float64(1400)/10),
			expectedSurcharge:   0,
			expectedDeliveryFee: 10 + 10 + int(2*float64(1400)/10),
			expectedErrorMsg:    "",
		},
		{
			name:                "Distance 1400 meters, Cart value 80",
			mockDistance:        1400,
			cartValue:           80,
			expectedTotalPrice:  80 + 0 + 10 + 10 + int(2*float64(1400)/10),
			expectedSurcharge:   0,
			expectedDeliveryFee: 10 + 10 + int(2*float64(1400)/10),
			expectedErrorMsg:    "",
		},
		{
			name:                "Distance 1800 meters, Cart value 40",
			mockDistance:        1800,
			cartValue:           40,
			expectedTotalPrice:  40 + (50 - 40) + 10 + 10 + int(2*float64(1800)/10),
			expectedSurcharge:   10,
			expectedDeliveryFee: 10 + 10 + int(2*float64(1800)/10),
			expectedErrorMsg:    "",
		},
		{
			name:                "Distance 2000 meters, Cart value 40",
			mockDistance:        2000,
			cartValue:           0,
			expectedTotalPrice:  0,
			expectedSurcharge:   0,
			expectedDeliveryFee: 0,
			expectedErrorMsg:    "location out of delivery range",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			queries.CartValue = tc.cartValue

			mockDistanceCalculator := &mockDistanceCalculator{mockDistance: tc.mockDistance}

			result, err := Calculator(queries, venue, mockDistanceCalculator)

			if err != nil {
				if err.Error() != tc.expectedErrorMsg {
					t.Errorf("expected error message %s, got %s", tc.expectedErrorMsg, err.Error())
				}
			} else {
				if tc.expectedErrorMsg != "" {
					t.Errorf("expected error message %s, but got no error", tc.expectedErrorMsg)
				}

				if result.TotalPrice != tc.expectedTotalPrice {
					t.Errorf("expected total price %d, got %d", tc.expectedTotalPrice, result.TotalPrice)
				}

				if result.SmallOrderSurcharge != tc.expectedSurcharge {
					t.Errorf("expected small order surcharge %d, got %d", tc.expectedSurcharge, result.SmallOrderSurcharge)
				}

				if result.Delivery.Fee != tc.expectedDeliveryFee {
					t.Errorf("expected delivery fee %d, got %d", tc.expectedDeliveryFee, result.Delivery.Fee)
				}
			}
		})
	}
}
