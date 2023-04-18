package mapbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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

// RateLimit represents a set of operations that share a rate limit
// see https://docs.mapbox.com/api/overview/#rate-limits
type RateLimit string

const (
	Geocoding = "geocoding"
	Matrix    = "matrix"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	apiKey     string
	// Referer is needed when URL restrictions are enforced, see https://docs.mapbox.com/accounts/guides/tokens/#url-restrictions
	Referer        string
	rateLimits     map[RateLimit]time.Time
	rateLimitMutex sync.RWMutex
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
		rateLimits: make(map[RateLimit]time.Time),
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
	if err := c.checkRateLimit(Matrix); err != nil {
		return nil, err
	}
	return directionsMatrix(ctx, c, req)
}

func (c *Client) ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*ReverseGeocodeResponse, error) {
	if err := c.checkRateLimit(Geocoding); err != nil {
		return nil, err
	}
	return reverseGeocode(ctx, c, req)
}

func (c *Client) ForwardGeocode(ctx context.Context, req *ForwardGeocodeRequest) (*ForwardGeocodeResponse, error) {
	if err := c.checkRateLimit(Geocoding); err != nil {
		return nil, err
	}
	return forwardGeocode(ctx, c, req)
}

//////////////////////////////////////////////////////////////////

func (c *Client) get(ctx context.Context, relPath string, query url.Values) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, relPath, query)
}

func (c *Client) do(ctx context.Context, httpVerb, relPath string, query url.Values) (*http.Response, error) {
	// remove empty entries
	for k := range query {
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
	if c.Referer != "" {
		req.Header.Set("Referer", c.Referer)
	}

	return c.httpClient.Do(req)
}

func (c *Client) handleResponse(apiResponse *http.Response, response interface{}, rateLimit RateLimit) error {
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
			return NewMapboxError(apiResponse.StatusCode, "")
		}

		// If rate limited, hold off till the next X-Rate-Limit-Reset
		if apiResponse.StatusCode == 429 {
			resetUnix, err := strconv.Atoi(apiResponse.Header.Get("X-Rate-Limit-Reset"))
			if err == nil {
				c.rateLimitMutex.Lock()
				defer c.rateLimitMutex.Unlock()
				c.rateLimits[rateLimit] = time.Unix(int64(resetUnix), 0)
			}
		}
		return NewMapboxError(apiResponse.StatusCode, errorResponse.Message)
	}

	// convert to response
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to read body. %w", err)
	}

	return nil
}

func (c *Client) rateLimit(rl RateLimit) time.Time {
	c.rateLimitMutex.RLock()
	defer c.rateLimitMutex.RUnlock()
	return c.rateLimits[rl]
}

func (c *Client) checkRateLimit(rl RateLimit) error {
	reset := c.rateLimit(rl)

	// No reset set
	if reset.IsZero() {
		return nil
	}
	// Reset reached
	if reset.Before(time.Now()) {
		c.rateLimitMutex.Lock()
		defer c.rateLimitMutex.Unlock()

		c.rateLimits[rl] = time.Time{}
		return nil
	}
	// Reset still in future
	return NewMapboxError(429, fmt.Sprintf("Rate limiting %v requests", rl))
}
