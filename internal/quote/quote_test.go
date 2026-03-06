package quote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// writeJSON is a test helper that writes a JSON body, failing the test on error.
func writeJSON(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := fmt.Fprint(w, body); err != nil {
		t.Errorf("writeJSON: %v", err)
	}
}

func TestFetch_ZenQuotes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `[{"q":"Be yourself","a":"Wilde"}]`)
	}))
	defer srv.Close()

	api := API{URL: srv.URL, Parse: DefaultAPIs[0].Parse}
	got, err := Fetch([]API{api}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Be yourself - Wilde" {
		t.Errorf("got %q", got)
	}
}

func TestFetch_QuoteSlate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `{"quote":"Keep going","author":"Unknown"}`)
	}))
	defer srv.Close()

	api := API{URL: srv.URL, Parse: DefaultAPIs[1].Parse}
	got, err := Fetch([]API{api}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Keep going - Unknown" {
		t.Errorf("got %q", got)
	}
}

func TestFetch_FallsThrough(t *testing.T) {
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer bad.Close()

	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `{"quote":"Persist","author":"Someone"}`)
	}))
	defer good.Close()

	apis := []API{
		{URL: bad.URL, Parse: DefaultAPIs[0].Parse},
		{URL: good.URL, Parse: DefaultAPIs[1].Parse},
	}
	got, err := Fetch(apis, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Persist - Someone" {
		t.Errorf("got %q", got)
	}
}

func TestFetch_AllFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	apis := []API{
		{URL: srv.URL, Parse: DefaultAPIs[0].Parse},
		{URL: srv.URL, Parse: DefaultAPIs[1].Parse},
	}
	_, err := Fetch(apis, 5*time.Second)
	if err == nil {
		t.Error("expected error when all APIs fail")
	}
}

func TestFetch_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `not json`)
	}))
	defer srv.Close()

	apis := []API{{URL: srv.URL, Parse: DefaultAPIs[0].Parse}}
	_, err := Fetch(apis, 5*time.Second)
	if err == nil {
		t.Error("expected error for invalid JSON response")
	}
}
