package mapbox

import (
	"context"
	"testing"
)

func checkforwardGeocodeRequestURL(t *testing.T, req *ForwardGeocodeRequest, expectedURL string) {
	t.Helper()
	client, requests := mockClient()
	go client.ForwardGeocode(context.Background(), req)

	httpReq := <-requests
	actualURL := httpReq.URL.RequestURI()
	if expectedURL != actualURL {
		t.Errorf("expected:\n%s, got:\n%s", expectedURL, actualURL)
	}
}

func TestForwardGeocodeURLEncoding(t *testing.T) {
	checkforwardGeocodeRequestURL(t, &ForwardGeocodeRequest{
		SearchText: "query with special chars:/; ",
	}, `/search/geocode/v6/forward?autocomplete=false&q=query+with+special+chars%3A%2F%3B+`)
}
