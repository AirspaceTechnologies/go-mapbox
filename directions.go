package mapbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	directionsPath = "directions"
)

type DirectionsRequest struct {
	// required
	Profile     Profile
	Coordinates Coordinates

	// optional
	Alternatives        *bool
	Annotations         Annotations
	AvoidManeuverRadius int // Possible values are in the range from 1 to 1000
	ContinueStraight    *bool
	Excludes            Excludes
	Geometries          Geometries
	Includes            Includes
	Overview            Overview
	Approaches          Approaches
	Steps               *bool
	BannerInstructions  *bool
	Language            string
	RoundaboutExits     *bool
	VoiceInstructions   *bool
	VoiceUnits          VoiceUnits
	Waypoints           DirectionWaypoints
	WaypointsPerRoute   *bool
	WaypointNames       WaypointNames
	WaypointTargets     WaypointTargets

	// Optional parameters for the mapbox/walking profile
	WalkingSpeed float32
	WalkwayBias  float32

	// Optional parameters for the mapbox/driving profile
	AlleyBias float32
	ArriveBy  ArriveBy
	DepartAt  DepartAt
	MaxHeight int
	MaxWidth  int
	MaxWeight int

	// Optional parameters for the mapbox/driving-traffic profile
	SnappingIncludeClosures       *bool
	SnappingIncludeStaticClosures *bool
}

// https://docs.mapbox.com/api/navigation/directions/#required-parameters
func directions(ctx context.Context, client *Client, req *DirectionsRequest) (*DirectionsResponse, error) {
	relPath := fmt.Sprintf("%v/%v/%v/%v", directionsPath, v5, req.Profile, req.Coordinates.WGS84Format())

	query := url.Values{}

	query.Set("access_token", client.apiKey)

	if req.Alternatives != nil {
		query.Set("alternatives", strconv.FormatBool(*req.Alternatives))
	}

	if len(req.Annotations) != 0 {
		// Must be used in conjunction with overview=full
		query.Set("annotations", req.Annotations.query())
		req.Overview = OverviewFull
	}

	if req.AvoidManeuverRadius != 0 {
		query.Set("avoid_maneuver_radius", strconv.Itoa(req.AvoidManeuverRadius))
	}

	if req.ContinueStraight != nil {
		query.Set("continue_straight", strconv.FormatBool(*req.ContinueStraight))
	}

	if len(req.Excludes) != 0 {
		query.Set("exclude", req.Excludes.query())
	}

	if req.Geometries != "" {
		query.Set("geometries", string(req.Geometries))
	}

	if len(req.Includes) != 0 {
		query.Set("include", req.Includes.query())
	}

	if req.Overview != "" {
		query.Set("overview", string(req.Overview))
	}

	if len(req.Approaches) != 0 {
		query.Set("approaches", req.Approaches.query())
	}

	if req.Steps != nil {
		query.Set("steps", strconv.FormatBool(*req.Steps))
	}

	if req.BannerInstructions != nil {
		query.Set("banner_instructions", strconv.FormatBool(*req.BannerInstructions))
	}

	if req.Language != "" {
		query.Set("language", req.Language)
	}

	if req.RoundaboutExits != nil {
		query.Set("roundabout_exits", strconv.FormatBool(*req.RoundaboutExits))
	}

	if req.VoiceInstructions != nil {
		query.Set("voice_instructions", strconv.FormatBool(*req.VoiceInstructions))
	}

	if req.VoiceUnits != "" {
		query.Set("voice_units", string(req.VoiceUnits))
	}
	if len(req.Waypoints) != 0 {
		query.Set("waypoints", req.Waypoints.query())
	}

	if req.WaypointsPerRoute != nil {
		query.Set("waypoints_per_route", strconv.FormatBool(*req.WaypointsPerRoute))
	}

	if len(req.WaypointNames) != 0 {
		query.Set("waypoint_names", req.WaypointNames.query())
	}

	if len(req.WaypointTargets) != 0 {
		query.Set("waypoint_targets", req.WaypointTargets.query())
	}
	if req.WalkingSpeed != 0 {
		query.Set("walking_speed", strconv.FormatFloat(float64(req.WalkingSpeed), 'f', 2, 32))
	}

	if req.WalkwayBias != 0 {
		query.Set("walkway_bias", strconv.FormatFloat(float64(req.WalkwayBias), 'f', 2, 32))
	}

	if req.AlleyBias != 0 {
		query.Set("alley_bias", strconv.FormatFloat(float64(req.AlleyBias), 'f', 2, 32))
	}

	if !req.ArriveBy.IsZero() {
		query.Set("arrive_by", req.ArriveBy.query())
	}

	if !req.DepartAt.IsZero() {
		query.Set("depart_at", req.DepartAt.query())
	}

	if req.MaxHeight != 0 {
		query.Set("max_height", strconv.Itoa(req.MaxHeight))
	}

	if req.MaxWidth != 0 {
		query.Set("max_width", strconv.Itoa(req.MaxWidth))
	}

	if req.MaxWeight != 0 {
		query.Set("max_weight", strconv.Itoa(req.MaxWeight))
	}

	if req.SnappingIncludeClosures != nil {
		query.Set("snapping_include_closures", strconv.FormatBool(*req.SnappingIncludeClosures))
	}

	if req.SnappingIncludeStaticClosures != nil {
		query.Set("snapping_include_static_closures", strconv.FormatBool(*req.SnappingIncludeStaticClosures))
	}

	apiResponse, err := client.get(ctx, relPath, query)
	if err != nil {
		return nil, err
	}

	var response DirectionsResponse
	if err := client.handleResponse(apiResponse, &response, DirectionsRateLimit); err != nil {
		return nil, err
	}

	return &response, nil
}
