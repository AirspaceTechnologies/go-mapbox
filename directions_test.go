package mapbox

import (
	"context"
	"testing"
)

func checkforwardDirectionsRequestURL(t *testing.T, req *DirectionsRequest, expectedURL string) {
	t.Helper()
	client, requests := mockClient()
	go client.Directions(context.Background(), req)

	httpReq := <-requests
	actualURL := httpReq.URL.RequestURI()
	if expectedURL != actualURL {
		t.Errorf("expected:\n%s, got:\n%s", expectedURL, actualURL)
	}
}

func TestForwardDirectionsURLEncoding(t *testing.T) {
	trueVal := true

	checkforwardDirectionsRequestURL(t, &DirectionsRequest{
		Profile: ProfileDrivingTraffic,
		Coordinates: Coordinates{
			Coordinate{Lat: 33.122508, Lng: -117.306786},
			Coordinate{Lat: 32.733810, Lng: -117.193443},
		},
		Alternatives:                  &trueVal,
		AvoidManeuverRadius:           1,
		ContinueStraight:              &trueVal,
		Steps:                         &trueVal,
		BannerInstructions:            &trueVal,
		Language:                      "en",
		RoundaboutExits:               &trueVal,
		VoiceInstructions:             &trueVal,
		VoiceUnits:                    "metric",
		WaypointsPerRoute:             &trueVal,
		Excludes:                      Excludes{ExcludeUnpaved, ExcludeCashOnlyTolls},
		Geometries:                    GeometriesGeoJSON,
		Includes:                      Includes{IncludeHov2, IncludeHot},
		Overview:                      OverviewSimplified,
		Approaches:                    Approaches{ApproachUnrestricted},
		WaypointNames:                 WaypointNames{"wp1", "wp2"},
		WaypointTargets:               WaypointTargets{"wpt1", "wpt2"},
		Annotations:                   Annotations{AnnotationDistance, AnnotationDuration},
		SnappingIncludeClosures:       &trueVal,
		SnappingIncludeStaticClosures: &trueVal,
	}, `/directions/v5/mapbox/driving-traffic/-117.306786,33.122508;-117.193443,32.73381?alternatives=true&annotations=distance%2Cduration&approaches=unrestricted&avoid_maneuver_radius=1&banner_instructions=true&continue_straight=true&exclude=unpaved%2Ccash_only_tolls&geometries=geojson&include=hov2%2Chot&language=en&overview=full&roundabout_exits=true&snapping_include_closures=true&snapping_include_static_closures=true&steps=true&voice_instructions=true&voice_units=metric&waypoint_names=wp1%3Bwp2&waypoint_targets=wpt1%3Bwpt2&waypoints_per_route=true`)
}
