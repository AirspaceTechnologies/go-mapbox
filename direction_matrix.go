package mapbox

import (
	"context"
	"fmt"
	"net/url"
)

const (
	directionsMatrixPath = "directions-matrix"
)

type DirectionsMatrixRequest struct {
	// required
	Profile     Profile
	Coordinates Coordinates

	// optional
	Annotations   Annotations
	Approaches    Approaches
	Destinations  Destinations
	Sources       Sources
	FallbackSpeed FallbackSpeed
}

type DirectionsMatrixResponse struct {
	Code         string       `json:"code"`
	Durations    [][]*float64 `json:"durations"`
	Distances    [][]*float64 `json:"distances"`
	Destinations []Waypoint   `json:"destinations"`
	Sources      []Waypoint   `json:"sources"`
}

// https://docs.mapbox.com/api/navigation/#matrix
func directionsMatrix(ctx context.Context, client *Client, req *DirectionsMatrixRequest) (*DirectionsMatrixResponse, error) {
	relPath := fmt.Sprintf("%v/%v/%v/%v", directionsMatrixPath, v1, req.Profile, req.Coordinates.WGS84Format())

	query := url.Values{}
	query.Set("access_token", client.apiKey)
	query.Set("annotations", req.Annotations.query())
	query.Set("approaches", req.Approaches.query())
	query.Set("destinations", req.Destinations.query())
	query.Set("sources", req.Sources.query())
	query.Set("fallback_speed", req.FallbackSpeed.query())

	apiResponse, err := client.get(ctx, relPath, query)
	if err != nil {
		return nil, err
	}

	var response DirectionsMatrixResponse
	if err := client.handleResponse(apiResponse, &response, Matrix); err != nil {
		return nil, err
	}

	return &response, nil
}
