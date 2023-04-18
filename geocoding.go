package mapbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	geocodePath = "geocoding"
)

//////////////////////////////////////////////////////////////////

type ReverseGeocodeRequest struct {
	// required
	Endpoint    Endpoint
	Coordinates Coordinates

	// optional
	Country     string
	Language    string
	Limit       int
	ReverseMode ReverseMode
	Routing     bool
	Types       Types
}

type ReverseGeocodeResponse struct {
	Type        string     `json:"type"`
	Query       []float64  `json:"query"`
	Features    []*Feature `json:"features"`
	Attribution string     `json:"attribution"`
}

//////////////////////////////////////////////////////////////////

type ForwardGeocodeRequest struct {
	// required
	Endpoint   Endpoint
	SearchText string

	// optional
	Autocomplete bool
	BBox         BoundingBox
	Country      string
	FuzzyMatch   bool
	Language     string
	Limit        int
	Proximity    Coordinate
	Routing      bool
	Types        Types
}

type ForwardGeocodeResponse struct {
	Type        string     `json:"type"`
	Query       []string   `json:"query"`
	Features    []*Feature `json:"features"`
	Attribution string     `json:"attribution"`
}

//////////////////////////////////////////////////////////////////

type Feature struct {
	ID                string      `json:"id"`
	Type              string      `json:"type"`
	PlaceType         []string    `json:"place_type"`
	Relevance         float64     `json:"relevance"`
	Address           string      `json:"address,omitempty"`
	Properties        *Properties `json:"properties,omitempty"`
	Text              string      `json:"text"`
	PlaceName         string      `json:"place_name"`
	MatchingText      string      `json:"matching_text,omitempty"`
	MatchingPlaceName string      `json:"matching_place_name,omitempty"`
	Language          string      `json:"language,omitempty"`
	Bbox              []float64   `json:"bbox,omitempty"`
	Center            []float64   `json:"center"`
	Geometry          *Geometry   `json:"geometry"`
	Context           []*Context  `json:"context,omitempty"`
}

// TODO: need to properly unmarshal this data. (In some cases) Mapbox returns {} for properties which creates an empty struct
type Properties struct {
	Accuracy  string `json:"accuracy,omitempty"`
	Address   string `json:"address,omitempty"`
	Category  string `json:"category,omitempty"`
	Maki      string `json:"maki,omitempty"`
	Landmark  bool   `json:"landmark,omitempty"`
	Wikidata  string `json:"wikidata,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
}

type Geometry struct {
	Coordinates  []float64 `json:"coordinates"`
	Type         string    `json:"type"`
	Interpolated bool      `json:"interpolated,omitempty"`
	Omitted      string    `json:"omitted,omitempty"`
}

type Context struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Wikidata  string `json:"wikidata,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
}

//////////////////////////////////////////////////////////////////

// https://docs.mapbox.com/api/search/#forward-geocoding
func forwardGeocode(ctx context.Context, client *Client, req *ForwardGeocodeRequest) (*ForwardGeocodeResponse, error) {
	relPath := fmt.Sprintf("%v/%v/%v/%v.json", geocodePath, v5, req.Endpoint, url.PathEscape(req.SearchText))

	query := url.Values{}
	query.Set("access_token", client.apiKey)
	query.Set("autocomplete", strconv.FormatBool(req.Autocomplete))
	if req.BBox.Min.Lat != 0 && req.BBox.Min.Lng != 0 {
		query.Set("bbox", req.BBox.query())
	}
	if req.Country != "" {
		query.Set("country", req.Country)
	}
	query.Set("fuzzyMatch", strconv.FormatBool(req.FuzzyMatch))
	if req.Language != "" {
		query.Set("language", req.Language)
	}
	if req.Limit != 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.Proximity.Lat != 0 {
		query.Set("proximity", req.Proximity.WGS84Format())
	}
	query.Set("routing", strconv.FormatBool(req.Routing))
	if len(req.Types) != 0 {
		query.Set("types", req.Types.query())
	}

	apiResponse, err := client.get(ctx, relPath, query)
	if err != nil {
		return nil, err
	}

	var response ForwardGeocodeResponse
	if err := client.handleResponse(apiResponse, &response, Geocoding); err != nil {
		return nil, err
	}

	return &response, nil
}

// https://docs.mapbox.com/api/search/#reverse-geocoding
func reverseGeocode(ctx context.Context, client *Client, req *ReverseGeocodeRequest) (*ReverseGeocodeResponse, error) {
	relPath := fmt.Sprintf("%v/%v/%v/%v.json", geocodePath, v5, req.Endpoint, req.Coordinates.WGS84Format())

	query := url.Values{}
	query.Set("access_token", client.apiKey)
	query.Set("country", req.Country)
	query.Set("language", req.Language)
	query.Set("limit", strconv.Itoa(req.Limit))
	query.Set("reverseMode", req.ReverseMode.query())
	query.Set("routing", strconv.FormatBool(req.Routing))
	query.Set("types", req.Types.query())

	apiResponse, err := client.get(ctx, relPath, query)
	if err != nil {
		return nil, err
	}

	var response ReverseGeocodeResponse
	if err := client.handleResponse(apiResponse, &response, Geocoding); err != nil {
		return nil, err
	}

	return &response, nil
}
