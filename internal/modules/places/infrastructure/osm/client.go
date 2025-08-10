package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OverpassClient is a minimal HTTP client for the Overpass API.
type OverpassClient struct {
	httpClient  *http.Client
	endpointURL string
}

func NewOverpassClient(overpassURL string) *OverpassClient {
	if overpassURL == "" {
		overpassURL = "https://overpass-api.de/api/interpreter"
	}
	return &OverpassClient{
		httpClient:  &http.Client{Timeout: 3 * time.Second},
		endpointURL: overpassURL,
	}
}

// Execute posts an Overpass QL query and returns the decoded response.
func (c *OverpassClient) Execute(ctx context.Context, ql string) (overpassResponse, error) {
	form := url.Values{}
	form.Set("data", ql)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpointURL, strings.NewReader(form.Encode()))
	if err != nil {
		return overpassResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return overpassResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return overpassResponse{}, fmt.Errorf("osm: status %d", resp.StatusCode)
	}

	var payload overpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return overpassResponse{}, err
	}
	return payload, nil
}
