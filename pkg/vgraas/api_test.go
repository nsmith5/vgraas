package vgraas

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReviews(t *testing.T) {
	req, err := http.NewRequest("GET", "/reviews/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	api := NewAPI(NewRAMRepo())
	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Review request failed with status %d", rr.Code)
	}

	if rr.Body.String() != "null\n" {
		t.Errorf("HEAD request should have no body found %s", rr.Body.String())
	}
}
