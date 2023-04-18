package mapbox

import (
	"context"
	"errors"
	"net/http"
	"testing"
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
				if i > len(responses) {
					return nil, errors.New("mockClient: not enough responses")
				}
				resp := responses[i]
				i++
				if i > len(responses) {
					close(ch)
				}
				return resp, nil
			}),
		},
	}
	return client, ch
}

func checkforwardGeocodeRequestURL(t *testing.T, req *ForwardGeocodeRequest, expectedURL string) {
	t.Helper()
	client, requests := mockClient(&http.Response{StatusCode: 200})
	go client.ForwardGeocode(context.Background(), req)

	httpReq := <-requests
	actualURL := httpReq.URL.RequestURI()
	if expectedURL != actualURL {
		t.Errorf("expected:\n%s, got:\n%s", expectedURL, actualURL)
	}
}

func TestForwardGeocodeURLEncoding(t *testing.T) {
	checkforwardGeocodeRequestURL(t, &ForwardGeocodeRequest{
		Endpoint:   EndpointPlaces,
		SearchText: "query with special chars:/; ",
	}, `/geocoding/v5/mapbox.places/query%20with%20special%20chars:%2F%3B%20.json?autocomplete=false&fuzzyMatch=false&routing=false`)
}
