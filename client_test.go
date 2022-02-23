package mapbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestClient_rateLimits(t *testing.T) {
	c, _ := NewClient(&MapboxConfig{
		APIKey: "test",
	})
	rlc := &rateLimitingClient{}
	c.httpClient = rlc

	req := ReverseGeocodeRequest{
		Endpoint: EndpointPlaces,
		Coordinates: Coordinates{
			Coordinate{
				Lat: 123.1,
				Lng: 123.2,
			},
		},

		Language: "en",
		Limit:    1,
	}

	// not rate limiting
	_, err := c.ReverseGeocode(
		context.Background(),
		&req,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// rate limiting
	rlc.rateLimiting = true
	_, err = c.ReverseGeocode(
		context.Background(),
		&req,
	)
	t.Logf("error: %v", err)
	if err.Error() != "api error(429): Too Many Requests" {
		t.Fatalf("Expected error, got none")
	}

	// Next request should be auto rate limited
	_, err = c.ReverseGeocode(
		context.Background(),
		&req,
	)
	if err.Error() != "api error(429): Rate limiting geocoding requests" {
		t.Fatalf("Expected error, got none")
	}

	// After reset, should be good to go again
	time.Sleep(3 * time.Second)
	rlc.rateLimiting = false
	_, err = c.ReverseGeocode(
		context.Background(),
		&req,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

////////////////////////////////////////////////////////////////////////////////

type rateLimitingClient struct {
	rateLimiting bool
	reset        time.Time
}

func (rlc *rateLimitingClient) Do(req *http.Request) (*http.Response, error) {
	if !rlc.rateLimiting {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		}, nil
	}

	rateLimitErr := ErrorResponse{
		Message: "Too Many Requests",
		Code:    "too_many_requests",
	}
	resJson, _ := json.Marshal(rateLimitErr)

	headers := http.Header{}
	headers.Add("X-Rate-Limit-Reset", fmt.Sprintf("%v", time.Now().Add(2*time.Second).Unix()))

	return &http.Response{
		StatusCode: 429,
		Body: ioutil.NopCloser(
			bytes.NewBuffer(resJson),
		),
		Header: headers,
	}, nil
}
