// Package quote fetches a motivational quote from one of several public APIs,
// falling through to the next API if one fails or is rate-limited.
package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// API describes a single quote provider.
type API struct {
	// URL is the endpoint to GET.
	URL string
	// Parse extracts a "quote - author" string from the response body.
	Parse func(body []byte) (string, error)
}

// DefaultAPIs is the list of quote providers used by the Bash version,
// in the same priority order.
var DefaultAPIs = []API{
	{
		URL: "https://zenquotes.io/api/random",
		Parse: func(body []byte) (string, error) {
			var results []struct {
				Q string `json:"q"`
				A string `json:"a"`
			}
			if err := json.Unmarshal(body, &results); err != nil || len(results) == 0 {
				return "", fmt.Errorf("quote: zenquotes: parse error")
			}
			return results[0].Q + " - " + results[0].A, nil
		},
	},
	{
		URL: "https://quoteslate.vercel.app/api/quotes/random",
		Parse: func(body []byte) (string, error) {
			var result struct {
				Quote  string `json:"quote"`
				Author string `json:"author"`
			}
			if err := json.Unmarshal(body, &result); err != nil || result.Quote == "" {
				return "", fmt.Errorf("quote: quoteslate: parse error")
			}
			return result.Quote + " - " + result.Author, nil
		},
	},
}

// Fetch tries each API in order until one returns a non-empty quote.
// timeout applies to each individual HTTP request.
// Returns an error only if all APIs fail.
func Fetch(apis []API, timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}
	var lastErr error

	for _, api := range apis {
		q, err := fetchOne(client, api)
		if err == nil && q != "" {
			return q, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return "", fmt.Errorf("quote: all APIs failed: %w", lastErr)
	}
	return "", fmt.Errorf("quote: all APIs returned empty quotes")
}

func fetchOne(client *http.Client, api API) (string, error) {
	req, err := http.NewRequest(http.MethodGet, api.URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "svelte-next/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("quote: HTTP %d from %s", resp.StatusCode, api.URL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return api.Parse(body)
}
