package mapbox

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

const (
	directionsMatrix = "directions-matrix"
)

type DirectionsMatrixRequest struct {
	// required
	Profile Profile
	Coordinates Coordinates

	// optional
	Annotations Annotations
	Approaches Approaches
	Destinations Destinations
	Sources Sources
	FallbackSpeed FallbackSpeed

}

type DirectionsMatrixResponse struct {
	Code string `json:"code"`
	Durations [][]float64 `json:"durations"`
	Distances [][]float64 `json:"distances"`
	Destinations []Waypoint `json:"destinations"`
	Sources []Waypoint `json:"sources"`
}

// https://docs.mapbox.com/api/navigation/#matrix
func (c *Client) DirectionsMatrix(ctx context.Context, req DirectionsMatrixRequest) (*DirectionsMatrixResponse, error) {
	relPath := fmt.Sprintf("%v/%v/%v/%v", directionsMatrix, v1, req.Profile, req.Coordinates.WGS84Format())

	query := url.Values{}
	query.Set("access_token", c.apiKey)
	query.Set("annotations", strings.Join(req.Annotations.strings(), ","))
	query.Set("approaches", strings.Join(req.Approaches.strings(), ";"))
	query.Set("destinations", strings.Join(req.Destinations.strings(), ";"))
	query.Set("sources", strings.Join(req.Sources.strings(), ";"))
	query.Set("fallback_speed", req.FallbackSpeed.String())

	apiResponse, e := c.get(ctx, relPath, query)
	if e != nil {
		return nil, e
	}

	var response DirectionsMatrixResponse
	if err := c.handleResponse(apiResponse, &response); err != nil {
		return nil, err
	}

	return &response, nil
}