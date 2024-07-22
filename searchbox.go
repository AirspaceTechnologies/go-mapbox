package mapbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	SearchboxReverseEndpoint = "/search/searchbox/v1/reverse"
)

type SearchboxReverseRequest struct {
	Coordinate

	// optional
	Country  string `json:"country,omitempty"`
	Language string `json:"language,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Types    Types  `json:"types,omitempty"`
}

type SearchboxReverseResponse struct {
	Type        string     `json:"type"`
	Features    []*Feature `json:"features"`
	Attribution string     `json:"attribution"`
}

type SearchboxReverseFeature struct {
	ID         string                      `json:"id"`
	Type       string                      `json:"type"`
	Geometry   *Geometry                   `json:"geometry"`
	Properties *SearchboxReverseProperties `json:"properties,omitempty"`
}

type SearchboxReverseProperties struct {
	MapboxID       string             `json:"mapbox_id"`
	FeatureType    Type               `json:"feature_type"`
	Name           string             `json:"name"`
	NamePreferred  string             `json:"name_preferred"`
	PlaceFormatted string             `json:"place_formatted"`
	FullAddress    string             `json:"full_address"`
	Coordinates    ExtendedCoordinate `json:"coordinates"`
	Context        map[Type]Context   `json:"context,omitempty"`
	BoundingBox    []float64          `json:"bbox,omitempty"`
	Language       string             `json:"language"`
	Maki           string             `json:"maki"`
	POICategory    []string           `json:"poi_category"`
	POICategoryIDs []string           `json:"poi_category_ids"`
	Brand          string             `json:"brand"`
	BrandID        string             `json:"brand_id"`
	ExternalIDs    map[string]string  `json:"external_ids,omitempty"`
	Metadata       map[string]string  `json:"metadata,omitempty"`
}

// https://docs.mapbox.com/api/search/search-box/#reverse-lookup
func searchboxReverse(ctx context.Context, client *Client, req *SearchboxReverseRequest) (*SearchboxReverseResponse, error) {
	query := url.Values{}
	query.Set("access_token", client.apiKey)
	query.Set("latitude", strconv.FormatFloat(req.Lat, 'f', -1, 64))
	query.Set("longitude", strconv.FormatFloat(req.Lng, 'f', -1, 64))

	if req.Country != "" {
		query.Set("country", req.Country)
	}

	if req.Language != "" {
		query.Set("language", req.Language)
	}

	if req.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%v", req.Limit))
	}

	if len(req.Types) > 0 {
		query.Set("types", req.Types.query())
	}

	apiResponse, err := client.get(ctx, SearchboxReverseEndpoint, query)
	if err != nil {
		return nil, err
	}

	var response SearchboxReverseResponse
	if err := client.handleResponse(apiResponse, &response, SearchboxRateLimit); err != nil {
		return nil, err
	}

	return &response, nil
}
