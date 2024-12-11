package mapbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	GeocodingBatchEndpoint   = "/search/geocode/v6/batch"
	GeocodingReverseEndpoint = "/search/geocode/v6/reverse"
	GeocodingForwardEndpoint = "/search/geocode/v6/forward"
)

//////////////////////////////////////////////////////////////////

type ReverseGeocodeRequest struct {
	Coordinate

	// optional
	Country  string `json:"country,omitempty"`
	Language string `json:"language,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Types    Types  `json:"types,omitempty"`
}

type ForwardGeocodeBatchRequest []ForwardGeocodeRequest

type ForwardGeocodeBatchResponse struct {
	Batch []GeocodeResponse `json:"batch"`
}

type ReverseGeocodeBatchRequest []ReverseGeocodeRequest

type ForwardGeocodeRequest struct {
	SearchText   string
	AddressLine1 string
	Postcode     string
	Place        string
	Autocomplete bool
	BBox         BoundingBox
	Country      string
	Language     string
	Limit        int
	Proximity    Coordinate
	Types        Types
}

func (r ForwardGeocodeRequest) MarshalJSON() ([]byte, error) {
	type forwardGeocodeRequest struct {
		SearchText   string       `json:"q,omitempty"`
		AddressLine1 string       `json:"address_line1,omitempty"`
		Postcode     string       `json:"postcode,omitempty"`
		Place        string       `json:"place,omitempty"`
		Autocomplete bool         `json:"autocomplete,omitempty"`
		BBox         *BoundingBox `json:"bbox,omitempty"`
		Country      string       `json:"country,omitempty"`
		Language     string       `json:"language,omitempty"`
		Limit        int          `json:"limit,omitempty"`
		Proximity    *Coordinate  `json:"proximity,omitempty"`
		Types        []string     `json:"types,omitempty"`
	}

	var resp = forwardGeocodeRequest{
		SearchText:   r.SearchText,
		AddressLine1: r.AddressLine1,
		Postcode:     r.Postcode,
		Place:        r.Place,
		Autocomplete: r.Autocomplete,
		Country:      r.Country,
		Language:     r.Language,
		Limit:        r.Limit,
	}

	if !r.BBox.Min.IsZero() && !r.BBox.Max.IsZero() {
		resp.BBox = &r.BBox
	}

	types := r.Types.strings()
	if len(types) > 0 {
		resp.Types = types
	}

	return json.Marshal(resp)
}

type GeocodeResponse struct {
	Type        string     `json:"type"`
	Features    []*Feature `json:"features"`
	Attribution string     `json:"attribution"`
}

type GeocodeBatchResponse struct {
	Batch []GeocodeResponse `json:"batch"`
}

//////////////////////////////////////////////////////////////////

type Feature struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Geometry   *Geometry   `json:"geometry"` // The center of Properties.BoundingBox
	Properties *Properties `json:"properties,omitempty"`
}

type Properties struct {
	MapboxID       string             `json:"mapbox_id"`
	FeatureType    Type               `json:"feature_type"`
	Name           string             `json:"name"`
	NamePreferred  string             `json:"name_preferred"`
	PlaceFormatted string             `json:"place_formatted"`
	FullAddress    string             `json:"full_address"`
	Coordinates    ExtendedCoordinate `json:"coordinates"`
	Context        map[Type]Context   `json:"context,omitempty"`
	BoundingBox    []float64          `json:"bbox,omitempty"`
	MatchCode      *MatchCode         `json:"match_code,omitempty"`
}

// There are many different types of context objects, which are all mashed together
// here.
// https://docs.mapbox.com/api/search/geocoding/#the-context-object
type Context struct {
	// Always present
	MapboxID string `json:"mapbox_id"`
	Name     string `json:"name"`

	// Optional but shared between many context types
	WikidataID string `json:"wikidata_id,omitempty"`

	// Region fields
	RegionCode     string `json:"region_code,omitempty"`
	RegionCodeFull string `json:"region_code_full,omitempty"`

	// Address fields
	AddressNumber string `json:"address_number,omitempty"`
	StreetName    string `json:"street_name,omitempty"`

	// Country fields
	CountryCode       string `json:"country_code,omitempty"`
	CountryCodeAlpha3 string `json:"country_code_alpha_3,omitempty"`
}

type MatchCode struct {
	AddressNumber MatchCodeValue      `json:"address_number"`
	Street        MatchCodeValue      `json:"street"`
	Postcode      MatchCodeValue      `json:"postcode"`
	Place         MatchCodeValue      `json:"place"`
	Region        MatchCodeValue      `json:"region"`
	Locality      MatchCodeValue      `json:"locality"`
	Country       MatchCodeValue      `json:"country"`
	Confidence    MatchCodeConfidence `json:"confidence"`
}

type MatchCodeConfidence string
type MatchCodeValue string

const (
	MatchCodeConfidenceExact  MatchCodeConfidence = "exact"
	MatchCodeConfidenceHigh   MatchCodeConfidence = "high"
	MatchCodeConfidenceMedium MatchCodeConfidence = "medium"
	MatchCodeConfidenceLow    MatchCodeConfidence = "low"

	MatchCodeValueMatched       MatchCodeValue = "matched"
	MatchCodeValueUnmatched     MatchCodeValue = "unmatched"
	MatchCodeValueNotApplicable MatchCodeValue = "not_applicable"
	MatchCodeValueInferred      MatchCodeValue = "inferred"
	MatchCodeValuePlausible     MatchCodeValue = "plausible"
)

//////////////////////////////////////////////////////////////////

// https://docs.mapbox.com/api/search/geocoding/#forward-geocoding-with-search-text-input
func forwardGeocode(ctx context.Context, client *Client, req *ForwardGeocodeRequest) (*GeocodeResponse, error) {
	query := url.Values{}
	query.Set("access_token", client.apiKey)
	query.Set("autocomplete", strconv.FormatBool(req.Autocomplete))

	if req.SearchText != "" {
		query.Set("q", req.SearchText)
	}
	if req.AddressLine1 != "" {
		query.Set("address_line1", req.AddressLine1)
	}
	if req.Postcode != "" {
		query.Set("postcode", req.Postcode)
	}
	if req.Place != "" {
		query.Set("place", req.Place)
	}
	if !req.BBox.Min.IsZero() {
		query.Set("bbox", req.BBox.query())
	}

	if req.Country != "" {
		query.Set("country", req.Country)
	}

	if req.Language != "" {
		query.Set("language", req.Language)
	}

	if req.Limit != 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}

	if !req.Proximity.IsZero() {
		query.Set("proximity", req.Proximity.WGS84Format())
	}

	if len(req.Types) != 0 {
		query.Set("types", req.Types.query())
	}

	apiResponse, err := client.get(ctx, GeocodingForwardEndpoint, query)
	if err != nil {
		return nil, err
	}

	var response GeocodeResponse
	if err := client.handleResponse(apiResponse, &response, GeocodingRateLimit); err != nil {
		return nil, err
	}

	return &response, nil
}

// https://docs.mapbox.com/api/search/geocoding/#batch-geocoding
func forwardGeocodeBatch(ctx context.Context, client *Client, req ForwardGeocodeBatchRequest) (*GeocodeBatchResponse, error) {
	query := url.Values{}
	query.Set("access_token", client.apiKey)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiResponse, err := client.post(ctx, GeocodingBatchEndpoint, query, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var response *GeocodeBatchResponse
	if err := client.handleResponse(apiResponse, &response, GeocodingRateLimit); err != nil {
		return nil, err
	}

	return response, nil
}

// https://docs.mapbox.com/api/search/geocoding/#reverse-geocoding
func reverseGeocode(ctx context.Context, client *Client, req *ReverseGeocodeRequest) (*GeocodeResponse, error) {
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

	apiResponse, err := client.get(ctx, GeocodingReverseEndpoint, query)
	if err != nil {
		return nil, err
	}

	var response GeocodeResponse
	if err := client.handleResponse(apiResponse, &response, GeocodingRateLimit); err != nil {
		return nil, err
	}

	return &response, nil
}

// https://docs.mapbox.com/api/search/geocoding/#batch-geocoding, but only supports reverse
func reverseGeocodeBatch(ctx context.Context, client *Client, req ReverseGeocodeBatchRequest) (*GeocodeBatchResponse, error) {
	query := url.Values{}
	query.Set("access_token", client.apiKey)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiResponse, err := client.post(ctx, GeocodingBatchEndpoint, query, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var response *GeocodeBatchResponse
	if err := client.handleResponse(apiResponse, &response, GeocodingRateLimit); err != nil {
		return nil, err
	}

	return response, nil
}
