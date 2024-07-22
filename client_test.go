package mapbox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (rt roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

func mockClient(responses ...*http.Response) (*Client, chan *http.Request) {
	ch := make(chan *http.Request)
	i := 0
	client := &Client{
		httpClient: &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				ch <- r
				if i >= len(responses) {
					return nil, errors.New("mockClient: not enough responses")
				}
				resp := responses[i]
				i++
				return resp, nil
			}),
		},
	}
	return client, ch
}

////////////////////////////////////////////////////////////////////////////////

func TestClient_raceCondition(t *testing.T) {

	c, _ := NewClient(&MapboxConfig{
		APIKey: "test",
	})
	rlc := &rateLimitingClient{}
	c.httpClient = rlc

	req := ReverseGeocodeRequest{
		Coordinate: Coordinate{
			Lat: 123.1,
			Lng: 123.2,
		},
		Language: "en",
		Limit:    1,
	}

	// Set limit, then run requests asynchronously until the limit is reset
	n := 50
	var wg sync.WaitGroup
	wg.Add(n * 2)

	rlc.rateLimiting = true
	c.ReverseGeocode(context.Background(), &req)
	rlc.rateLimiting = false

	go func() {
		for i := 0; i < n; i++ {
			time.Sleep(50 * time.Millisecond)
			c.ReverseGeocode(context.Background(), &req)
			wg.Done()
		}
	}()
	go func() {
		for i := 0; i < n; i++ {
			time.Sleep(50 * time.Millisecond)
			c.ReverseGeocode(context.Background(), &req)
			wg.Done()
		}
	}()

	wg.Wait()
}

func TestClient_rateLimits(t *testing.T) {
	c, _ := NewClient(&MapboxConfig{
		APIKey: "test",
	})
	rlc := &rateLimitingClient{}
	c.httpClient = rlc

	req := ReverseGeocodeRequest{
		Coordinate: Coordinate{
			Lat: 123.1,
			Lng: 123.2,
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
	time.Sleep(2 * time.Second)
	rlc.rateLimiting = false
	_, err = c.ReverseGeocode(
		context.Background(),
		&req,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClientReferer(t *testing.T) {
	expectedReferer := "https://example.com/"

	client, requests := mockClient()
	client.Referer = expectedReferer
	go client.do(context.Background(), "GET", "/", url.Values{}, nil)

	httpReq := <-requests
	actualReferer := httpReq.Referer()
	if expectedReferer != actualReferer {
		t.Errorf("expected referer: %q, got: %q", expectedReferer, actualReferer)
	}
}

////////////////////////////////////////////////////////////////////////////////

type rateLimitingClient struct {
	rateLimiting bool
	reset        time.Time //nolint:unused
}

func (rlc *rateLimitingClient) Do(req *http.Request) (*http.Response, error) {
	if !rlc.rateLimiting {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("{}")),
		}, nil
	}

	rateLimitErr := ErrorResponse{
		Message: "Too Many Requests",
		Code:    "too_many_requests",
	}
	resJson, _ := json.Marshal(rateLimitErr)

	headers := http.Header{}
	headers.Add("X-Rate-Limit-Reset", fmt.Sprintf("%v", time.Now().Add(1*time.Second).Unix()))

	return &http.Response{
		StatusCode: 429,
		Body: io.NopCloser(
			bytes.NewBuffer(resJson),
		),
		Header: headers,
	}, nil
}
