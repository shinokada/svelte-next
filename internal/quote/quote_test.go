package quote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// testTimeout is a short duration sufficient for in-process httptest servers.
const testTimeout = 500 * time.Millisecond

// Parser references extracted by format, independent of DefaultAPIs ordering.
var (
	parseZenQuotes  = DefaultAPIs[0].Parse // expects [{q,a}] array
	parseQuoteSlate = DefaultAPIs[1].Parse // expects {quote,author} object
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

	api := API{URL: srv.URL, Parse: parseZenQuotes}
	got, err := Fetch([]API{api}, testTimeout)
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

	api := API{URL: srv.URL, Parse: parseQuoteSlate}
	got, err := Fetch([]API{api}, testTimeout)
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
		{URL: bad.URL, Parse: parseZenQuotes},
		{URL: good.URL, Parse: parseQuoteSlate},
	}
	got, err := Fetch(apis, testTimeout)
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
		{URL: srv.URL, Parse: parseZenQuotes},
		{URL: srv.URL, Parse: parseQuoteSlate},
	}
	_, err := Fetch(apis, testTimeout)
	if err == nil {
		t.Error("expected error when all APIs fail")
	}
}

func TestFetch_FallsThrough_OnParseError(t *testing.T) {
	// First API returns 200 but unparseable JSON; second API should be tried.
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `not json`)
	}))
	defer bad.Close()

	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `{"quote":"Try again","author":"Nobody"}`)
	}))
	defer good.Close()

	apis := []API{
		{URL: bad.URL, Parse: parseZenQuotes},
		{URL: good.URL, Parse: parseQuoteSlate},
	}
	got, err := Fetch(apis, testTimeout)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Try again - Nobody" {
		t.Errorf("got %q", got)
	}
}

func TestFetch_NilParser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `{"quote":"hello","author":"world"}`)
	}))
	defer srv.Close()

	apis := []API{{URL: srv.URL, Parse: nil}}
	_, err := Fetch(apis, testTimeout)
	if err == nil {
		t.Error("expected error for nil parser")
	}
}

func TestFetch_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, `not json`)
	}))
	defer srv.Close()

	apis := []API{{URL: srv.URL, Parse: parseZenQuotes}}
	_, err := Fetch(apis, testTimeout)
	if err == nil {
		t.Error("expected error for invalid JSON response")
	}
}

func TestFetch_Timeout(t *testing.T) {
	// Handler delays longer than the client timeout.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * testTimeout)
		writeJSON(t, w, `[{"q":"Too late","a":"Nobody"}]`)
	}))
	defer srv.Close()

	apis := []API{{URL: srv.URL, Parse: parseZenQuotes}}
	_, err := Fetch(apis, testTimeout)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestFetch_OversizedBody(t *testing.T) {
	// Response body exceeds the 1 MiB cap.
	oversized := strings.Repeat("x", (1<<20)+1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, oversized); err != nil {
			t.Errorf("writeOversized: %v", err)
		}
	}))
	defer srv.Close()

	apis := []API{{URL: srv.URL, Parse: parseZenQuotes}}
	_, err := Fetch(apis, testTimeout)
	if err == nil {
		t.Error("expected error for oversized response body")
	}
}
