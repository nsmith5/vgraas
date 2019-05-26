package vgraas

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Request struct {
	verb string
	path string
	body string
}

func TestReviewsAPI(t *testing.T) {
	requests := []Request{
		Request{"GET", "/reviews/", ""},
		Request{"POST", "/reviews/", `{"author": "me", "body": "this andthat"}`},
		Request{"GET", "/reviews/0", ""},
		Request{"PUT", "/reviews/0", `{"author": "notme", "body": "this andthat"}`},
		Request{"DELETE", "/reviews/0", ""},
	}

	api := NewAPI(NewRAMRepo())

	for _, request := range requests {
		req, err := http.NewRequest(request.verb, request.path, strings.NewReader(request.body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		api.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %v failed with status %d", request, rr.Code)
		}
	}
}

func TestCommentsAPI(t *testing.T) {
	requests := []Request{
		Request{"POST", "/reviews/", `{"author": "me", "body": "this andthat"}`},
		Request{"POST", "/reviews/0/comments", `{"author": "guy", "body": "terrible!"}`},
		Request{"GET", "/reviews/0", ""},
		Request{"GET", "/reviews/0/comments", ""},
		Request{"PUT", "/reviews/0/comments/0", `{"author": "guy", "body": "ok!"}`},
		Request{"DELETE", "/reviews/0/comments/0", ""},
	}

	api := NewAPI(NewRAMRepo())

	for _, request := range requests {
		req, err := http.NewRequest(request.verb, request.path, strings.NewReader(request.body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		api.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %v failed with status %d", request, rr.Code)
		}
	}
}

func TestHealthEndpoint(t *testing.T) {
	api := NewAPI(NewRAMRepo())
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	api.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health endpoint didn't respond with 200")
	}
}
