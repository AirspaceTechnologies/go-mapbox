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
						{"location": [8.2, 48.9], "admin_index": 0, "geometry_index": 0},
						{"location": [8.1, 48.95], "admin_index": 1, "geometry_index": 2}
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
	if intersections[1].GeometryIndex != 2 {
		t.Errorf("expected geometry index 2, got %v", intersections[1].GeometryIndex)
	}
}

func TestRouteLegBorderCrossings(t *testing.T) {
	de := Admin{ISO_3166_1: "DE", ISO_3166_1_alpha3: "DEU"}
	fr := Admin{ISO_3166_1: "FR", ISO_3166_1_alpha3: "FRA"}

	cases := map[string]struct {
		leg      RouteLeg
		expected []BorderCrossing
	}{
		"single admin produces none": {
			leg: RouteLeg{
				Admins: []Admin{de},
				Steps: []Step{{
					Duration: 600, Distance: 10000,
					Intersections: []Intersection{{Location: []float64{8.2, 48.9}, AdminIndex: 0}},
				}},
			},
			expected: nil,
		},
		"crossing with annotations uses exact prefix sums": {
			leg: RouteLeg{
				Admins: []Admin{de, fr},
				Annotation: DirectionsAnnotation{
					Duration: []float64{100, 200, 300, 400},
					Distance: []float64{1000, 2000, 3000, 4000},
				},
				Steps: []Step{
					{
						Duration: 300, Distance: 3000,
						Intersections: []Intersection{
							{Location: []float64{8.3, 48.8}, AdminIndex: 0, GeometryIndex: 0},
						},
					},
					{
						Duration: 700, Distance: 7000,
						Intersections: []Intersection{
							{Location: []float64{8.2, 48.9}, AdminIndex: 0, GeometryIndex: 1},
							{Location: []float64{8.1, 48.95}, AdminIndex: 1, GeometryIndex: 3},
						},
					},
				},
			},
			expected: []BorderCrossing{{
				Location: []float64{8.1, 48.95},
				From:     de,
				To:       fr,
				Duration: 600,  // 100+200+300
				Distance: 6000, // 1000+2000+3000
			}},
		},
		"crossing without annotations falls back to completed steps": {
			leg: RouteLeg{
				Admins: []Admin{de, fr},
				Steps: []Step{
					{
						Duration: 300, Distance: 3000,
						Intersections: []Intersection{
							{Location: []float64{8.3, 48.8}, AdminIndex: 0},
						},
					},
					{
						Duration: 700, Distance: 7000,
						Intersections: []Intersection{
							{Location: []float64{8.1, 48.95}, AdminIndex: 1},
						},
					},
				},
			},
			expected: []BorderCrossing{{
				Location: []float64{8.1, 48.95},
				From:     de,
				To:       fr,
				Duration: 300,
				Distance: 3000,
			}},
		},
		"missing admin indexes produce none": {
			leg: RouteLeg{
				Admins: []Admin{de, fr},
				Steps: []Step{{
					Duration: 600, Distance: 10000,
					Intersections: []Intersection{
						{Location: []float64{8.3, 48.8}},
						{Location: []float64{8.1, 48.95}},
					},
				}},
			},
			expected: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tc.leg.BorderCrossings()

			if len(got) != len(tc.expected) {
				t.Fatalf("expected %v crossings, got %v", len(tc.expected), len(got))
			}
			for i, want := range tc.expected {
				if got[i].From != want.From || got[i].To != want.To {
					t.Errorf("crossing %v: expected %v->%v, got %v->%v", i, want.From.ISO_3166_1, want.To.ISO_3166_1, got[i].From.ISO_3166_1, got[i].To.ISO_3166_1)
				}
				if got[i].Duration != want.Duration || got[i].Distance != want.Distance {
					t.Errorf("crossing %v: expected duration/distance %v/%v, got %v/%v", i, want.Duration, want.Distance, got[i].Duration, got[i].Distance)
				}
				if got[i].Location[0] != want.Location[0] || got[i].Location[1] != want.Location[1] {
					t.Errorf("crossing %v: expected location %v, got %v", i, want.Location, got[i].Location)
				}
			}
		})
	}
}
