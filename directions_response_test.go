package mapbox

import (
	"encoding/json"
	"testing"
)

func TestDirectionsResponseUnmarshalAdmins(t *testing.T) {
	payload := `{
		"code": "Ok",
		"routes": [{
			"duration": 3600,
			"distance": 100000,
			"legs": [{
				"distance": 100000,
				"duration": 3600,
				"admins": [
					{"iso_3166_1": "DE", "iso_3166_1_alpha3": "DEU"},
					{"iso_3166_1": "FR", "iso_3166_1_alpha3": "FRA"}
				],
				"steps": [{
					"distance": 50000,
					"duration": 1800,
					"geometry": "",
					"intersections": [
						{"location": [8.2, 48.9], "admin_index": 0},
						{"location": [8.1, 48.95], "admin_index": 1}
					]
				}]
			}]
		}]
	}`

	var resp DirectionsResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	leg := resp.Routes[0].Legs[0]
	if len(leg.Admins) != 2 {
		t.Fatalf("expected 2 admins, got %v", len(leg.Admins))
	}
	if leg.Admins[0].ISO_3166_1 != "DE" || leg.Admins[1].ISO_3166_1 != "FR" {
		t.Errorf("expected admins DE,FR got %v,%v", leg.Admins[0].ISO_3166_1, leg.Admins[1].ISO_3166_1)
	}

	intersections := leg.Steps[0].Intersections
	if len(intersections) != 2 {
		t.Fatalf("expected 2 intersections, got %v", len(intersections))
	}
	if intersections[0].AdminIndex != 0 || intersections[1].AdminIndex != 1 {
		t.Errorf("expected admin indexes 0,1 got %v,%v", intersections[0].AdminIndex, intersections[1].AdminIndex)
	}
}
