package mapbox

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestIntegration_DirectionsAdminCrossings exercises real Mapbox responses to
// verify admins, admin_index, geometry_index, and annotations arrive as
// expected and that RouteLeg.AdminCrossings locates a plausible crossing.
func TestIntegration_DirectionsAdminCrossings(t *testing.T) {
	key := os.Getenv("API_KEY")
	if key == "" {
		t.Fatalf("missing Mapbox API key")
	}

	client, err := NewClient(&MapboxConfig{
		Timeout: 15 * time.Second,
		APIKey:  key,
	})
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	cases := map[string]struct {
		origin, destination Coordinate
		fromCountry         string
		toCountry           string
	}{
		"san diego to tijuana": {
			origin:      Coordinate{Lat: 32.7157, Lng: -117.1611},
			destination: Coordinate{Lat: 32.5149, Lng: -117.0382},
			fromCountry: "US",
			toCountry:   "MX",
		},
		"freiburg to mulhouse": {
			origin:      Coordinate{Lat: 47.9990, Lng: 7.8421},
			destination: Coordinate{Lat: 47.7508, Lng: 7.3359},
			fromCountry: "DE",
			toCountry:   "FR",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			trueVal := true
			resp, err := client.Directions(context.Background(), &DirectionsRequest{
				Profile:     ProfileDrivingTraffic,
				Coordinates: Coordinates{tc.origin, tc.destination},
				Geometries:  GeometriesPolyline,
				Overview:    OverviewFull,
				Steps:       &trueVal,
				Annotations: Annotations{AnnotationDistance, AnnotationDuration},
			})
			if err != nil {
				t.Fatalf("directions failed: %v", err)
			}
			if len(resp.Routes) == 0 || len(resp.Routes[0].Legs) == 0 {
				t.Fatalf("expected at least one route with one leg")
			}

			leg := resp.Routes[0].Legs[0]
			t.Logf("leg: %.0fs %.0fm admins=%v annotations=%v/%v",
				leg.Duration, leg.Distance, len(leg.Admins), len(leg.Annotation.Duration), len(leg.Annotation.Distance))

			if len(leg.Admins) < 2 {
				t.Fatalf("expected at least 2 admins, got %v: %+v", len(leg.Admins), leg.Admins)
			}
			if len(leg.Annotation.Duration) == 0 || len(leg.Annotation.Distance) == 0 {
				t.Errorf("expected duration/distance annotations, got %v/%v", len(leg.Annotation.Duration), len(leg.Annotation.Distance))
			}

			crossings := leg.AdminCrossings()
			if len(crossings) == 0 {
				t.Fatalf("expected at least one admin crossing")
			}

			c := crossings[0]
			t.Logf("crossing: %v -> %v at [lng %v, lat %v], %.0fs / %.0fm from leg start",
				c.From.ISO_3166_1, c.To.ISO_3166_1, c.Location[0], c.Location[1], c.Duration, c.Distance)

			if c.From.ISO_3166_1 != tc.fromCountry || c.To.ISO_3166_1 != tc.toCountry {
				t.Errorf("expected crossing %v->%v, got %v->%v", tc.fromCountry, tc.toCountry, c.From.ISO_3166_1, c.To.ISO_3166_1)
			}
			if c.Duration <= 0 || c.Duration >= leg.Duration {
				t.Errorf("expected crossing duration within (0, %.0f), got %.0f", leg.Duration, c.Duration)
			}
			if c.Distance <= 0 || c.Distance >= leg.Distance {
				t.Errorf("expected crossing distance within (0, %.0f), got %.0f", leg.Distance, c.Distance)
			}
		})
	}
}
