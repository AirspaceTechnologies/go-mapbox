package mapbox

import (
	"context"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testAPIDelay  = time.Second
	testLocations = map[string]struct {
		Lat     float64
		Lng     float64
		Country string
		Place   string
		POI     string
		Query   string
	}{
		"Eiffel Tower": {
			Lat:     48.858415953144025,
			Lng:     2.2944920264583892,
			Country: "France",
			Place:   "Paris",
			POI:     "Les Boutiques Officielles de la Tour Eiffel",
			Query:   "Av. Gustave Eiffel, 75007 Paris, France",
		},
		"Golden Gate Bridge": {
			Lat:     37.81999562350779,
			Lng:     -122.47855980298934,
			Country: "United States",
			Place:   "Sausalito",
			POI:     "Plaza Park Square",
			Query:   "Golden Gate Bridge, San Francisco, CA, United States",
		},
		"Machu Picchu": {
			Lat:     -13.163104764687816,
			Lng:     -72.54525137460071,
			Country: "Peru",
			Place:   "Machu Picchu",
			Query:   "Santuario Histórico de Machu Picchu, 08680, Peru",
		},
		"Victoria Falls": {
			Lat:     -17.925510375019098,
			Lng:     25.858544325497473,
			Country: "Zimbabwe",
			Place:   "Victoria Falls",
			Query:   "2 Livingstone Way, Victoria Falls, Zimbabwe",
		},
	}
)

func TestIntegration_ReverseGeocode(t *testing.T) {
	// ask for all supported even though some won't exist for the coordinate
	features := Types{
		TypeCountry,
		TypeRegion,
		TypePostcode,
		TypeDistrict,
		TypePlace,
		TypeLocality,
		TypeNeighborhood,
		TypeStreet,
		TypeAddress,
	}

	client, err := NewClient(&MapboxConfig{
		Timeout: 30 * time.Second,
		APIKey:  os.Getenv("API_KEY"),
	})

	if err != nil {
		t.Fatal(err)
	}

	for name, loc := range testLocations {
		t.Run(name, func(t *testing.T) {
			request := &ReverseGeocodeRequest{
				Coordinate: Coordinate{Lat: loc.Lat, Lng: loc.Lng},
				Language:   "en",
				Types:      features,
			}

			resp, err := client.ReverseGeocode(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}

			if resp == nil {
				t.Fatal("response should not be nil")
			}

			// just check the obvious ones
			compared := make(map[Type]struct{})
			for _, feature := range resp.Features {
				//nolint:exhaustive
				switch feature.Properties.FeatureType {
				case TypeCountry:
					compared[TypeCountry] = struct{}{}
					if feature.Properties.Name != loc.Country {
						t.Errorf("expected %v country %v to be %v", name, feature.Properties.Name, loc.Country)
					}
				case TypePlace:
					compared[TypePlace] = struct{}{}
					if feature.Properties.Name != loc.Place {
						t.Errorf("expected %v place %v to be %v", name, feature.Properties.Name, loc.Country)
					}
				}
			}

			if _, seen := compared[TypeCountry]; !seen {
				t.Error("response did not include a country feature")
			}

			if _, seen := compared[TypePlace]; !seen {
				t.Error("response did not include a place feature")
			}
		})

		time.Sleep(testAPIDelay)
	}
}

func TestIntegration_ReverseGeocodeBatch(t *testing.T) {
	// ask for all supported even though some won't exist for the coordinate
	features := Types{
		TypeCountry,
		TypeRegion,
		TypePostcode,
		TypeDistrict,
		TypePlace,
		TypeLocality,
		TypeNeighborhood,
		TypeStreet,
		TypeAddress,
	}

	client, err := NewClient(&MapboxConfig{
		Timeout: 30 * time.Second,
		APIKey:  os.Getenv("API_KEY"),
	})

	if err != nil {
		t.Fatal(err)
	}

	var keys []string // for maintaining order for response
	var requests ReverseGeocodeBatchRequest
	for key, loc := range testLocations {
		keys = append(keys, key)
		requests = append(requests, ReverseGeocodeRequest{
			Coordinate: Coordinate{Lat: loc.Lat, Lng: loc.Lng},
			Language:   "en",
			Types:      features,
		})
	}

	resps, err := client.ReverseGeocodeBatch(context.Background(), requests)
	if err != nil {
		t.Fatal(err)
	}

	if len(resps.Batch) != len(keys) {
		t.Fatalf("expected batch response length %v to be %v", len(resps.Batch), len(keys))
	}

	// just check the obvious ones
	for i, key := range keys {
		expected := testLocations[key]
		compared := make(map[Type]struct{})
		resp := resps.Batch[i]
		for _, feature := range resp.Features {
			//nolint:exhaustive
			switch feature.Properties.FeatureType {
			case TypeCountry:
				compared[TypeCountry] = struct{}{}
				if feature.Properties.Name != expected.Country {
					t.Errorf("expected %v country %v to be %v", key, feature.Properties.Name, expected.Country)
				}
			case TypePlace:
				compared[TypePlace] = struct{}{}
				if feature.Properties.Name != expected.Place {
					t.Errorf("expected %v place %v to be %v", key, feature.Properties.Name, expected.Country)
				}
			}
		}

		if _, seen := compared[TypeCountry]; !seen {
			t.Errorf("response for %v did not include a country feature", key)
		}

		if _, seen := compared[TypePlace]; !seen {
			t.Errorf("response for %v did not include a place feature", key)
		}
	}
}

func TestIntegration_ForwardGeocode(t *testing.T) {
	// ask for all supported even though some won't exist for the coordinate
	features := Types{
		TypeCountry,
		TypeRegion,
		TypePostcode,
		TypeDistrict,
		TypePlace,
		TypeLocality,
		TypeNeighborhood,
		TypeStreet,
		TypeAddress,
	}

	client, err := NewClient(&MapboxConfig{
		Timeout: 30 * time.Second,
		APIKey:  os.Getenv("API_KEY"),
	})

	if err != nil {
		t.Fatal(err)
	}

	for name, loc := range testLocations {
		t.Run(name, func(t *testing.T) {
			request := ForwardGeocodeRequest{SearchText: loc.Query, Types: features}
			resp, err := client.ForwardGeocode(context.Background(), &request)
			if err != nil {
				t.Fatal(err)
			}

			if resp == nil {
				t.Fatal("resp should not be nil")
			}

			// Check if response includes a nearby place
			compared := make(map[Type]struct{})
			for _, feature := range resp.Features {
				if feature.Properties.FeatureType != TypePlace {
					continue
				}

				// skip if already validated
				if _, seen := compared[TypePlace]; seen {
					continue
				}

				// not the match we're looking for: incorrect country
				if !strings.Contains(feature.Properties.FullAddress, loc.Country) {
					continue
				}

				// not the match we're looking for: latitude out of bounds
				if math.Abs(feature.Properties.Coordinates.Latitude-loc.Lat) > 0.5 {
					continue
				}

				// not the match we're looking for: longitude out of bounds
				if math.Abs(feature.Properties.Coordinates.Longitude-loc.Lng) > 0.5 {
					continue
				}

				// match found!
				compared[TypePlace] = struct{}{}
			}

			if _, seen := compared[TypePlace]; !seen {
				t.Errorf("response did not include a place with country %v and near coordinates [%v, %v]", loc.Country, loc.Lat, loc.Lng)
			}
		})
	}
}

func TestIntegration_SearchboxReverse(t *testing.T) {
	// ask for all supported even though some won't exist for the coordinate
	features := Types{
		TypePOI,
	}

	client, err := NewClient(&MapboxConfig{
		Timeout: 30 * time.Second,
		APIKey:  os.Getenv("API_KEY"),
	})

	if err != nil {
		t.Fatal(err)
	}

	for name, loc := range testLocations {
		t.Run(name, func(t *testing.T) {
			if loc.POI == "" {
				t.Skip(name, "does not currently return any point of interests -- skipping...")
				return
			}

			request := &SearchboxReverseRequest{
				Coordinate: Coordinate{Lat: loc.Lat, Lng: loc.Lng},
				Language:   "en",
				Types:      features,
			}

			resp, err := client.SearchboxReverse(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}

			if resp == nil {
				t.Fatal("response should not be nil")
			}

			// just check the obvious ones
			compared := make(map[Type]struct{})
			for _, feature := range resp.Features {
				if feature.Properties.FeatureType != TypePOI {
					continue
				}

				// skip if already validated
				if _, seen := compared[TypePOI]; seen {
					continue
				}

				// not the match we're looking for: incorrect POI
				if feature.Properties.Name != loc.POI {
					continue
				}

				// match found!
				compared[TypePOI] = struct{}{}
			}

			if _, seen := compared[TypePOI]; !seen {
				t.Error("response did not include the specified point of interest")
			}
		})

		time.Sleep(testAPIDelay)
	}
}
