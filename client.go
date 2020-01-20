package mapbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	ResponseOK = "Ok"

	baseUrl = "https://api.mapbox.com"
	v1      = "v1"
	v5      = "v5"
)

type MapboxConfig struct {
	Timeout time.Duration
	APIKey  string
}

type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient instantiates a new Mapbox client.
func NewClient(config *MapboxConfig) (*Client, error) {
	// Default timeout
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("missing Mapbox API key")
	}

	return &Client{
		httpClient: &http.Client{Timeout: config.Timeout},
		apiKey:     config.APIKey,
	}, nil
}

//////////////////////////////////////////////////////////////////

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Waypoint struct {
	Distance float64 `json:"distance"`
	Name     string  `json:"name"`
	Location []float64
}

//////////////////////////////////////////////////////////////////

func (c *Client) DirectionsMatrix(ctx context.Context, req *DirectionsMatrixRequest) (*DirectionsMatrixResponse, error) {
	return directionsMatrix(ctx, c, req)
}

func (c *Client) ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*ReverseGeocodeResponse, error) {
	return reverseGeocode(ctx, c, req)
}

func (c *Client) ForwardGeocode(ctx context.Context, req *ForwardGeocodeRequest) (*ForwardGeocodeResponse, error) {
	return forwardGeocode(ctx, c, req)
}

//////////////////////////////////////////////////////////////////

func (c *Client) get(ctx context.Context, relPath string, query url.Values) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, relPath, query)
}

func (c *Client) do(ctx context.Context, httpVerb, relPath string, query url.Values) (*http.Response, error) {
	// remove empty entries
	for k, _ := range query {
		if query.Get(k) == "" {
			query.Del(k)
		}
	}

	// safe to assume '?' as mapbox requires auth token as query param
	uri := fmt.Sprintf("%v/%v?%v", baseUrl, relPath, query.Encode())

	req, err := http.NewRequestWithContext(ctx, httpVerb, uri, nil)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) handleResponse(apiResponse *http.Response, response interface{}) error {
	defer apiResponse.Body.Close()

	// auth checking
	if apiResponse.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized request. Provide Mapbox API key")
	}

	body, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read body. %w", err)
	}

	// check for errors from Mapbox API (non 200 response)
	if apiResponse.StatusCode >= 400 && apiResponse.StatusCode <= 599 {
		var errorResponse ErrorResponse
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return fmt.Errorf("api error(%v): no body", apiResponse.StatusCode)
		}

		return fmt.Errorf("api error(%v): %v", apiResponse.StatusCode, errorResponse.Message)
	}

	//convert to response
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to read body. %w", err)
	}

	return nil
}
