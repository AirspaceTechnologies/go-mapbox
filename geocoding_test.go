package mapbox

import (
	"context"
	"net/http"
	"testing"
)

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
