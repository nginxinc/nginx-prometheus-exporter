package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const templateMetrics string = `Active connections: %d
server accepts handled requests
%d %d %d
Reading: %d Writing: %d Waiting: %d
`

// NginxClient allows you to fetch NGINX metrics from the stub_status page.
type NginxClient struct {
	httpClient  *http.Client
	apiEndpoint string
}

// StubStats represents NGINX stub_status metrics.
type StubStats struct {
	Connections StubConnections
	Requests    int64
}

// StubConnections represents connections related metrics.
type StubConnections struct {
	Active   int64
	Accepted int64
	Handled  int64
	Reading  int64
	Writing  int64
	Waiting  int64
}

// NewNginxClient creates an NginxClient.
func NewNginxClient(httpClient *http.Client, apiEndpoint string) *NginxClient {
	client := &NginxClient{
		apiEndpoint: apiEndpoint,
		httpClient:  httpClient,
	}

	return client
}

// GetStubStats fetches the stub_status metrics.
func (client *NginxClient) GetStubStats() (*StubStats, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a get request: %w", err)
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get %v: %w", client.apiEndpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v response, got %v", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}

	r := bytes.NewReader(body)
	stats, err := parseStubStats(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body %q: %w", string(body), err)
	}

	return stats, nil
}

func parseStubStats(r io.Reader) (*StubStats, error) {
	var s StubStats
	if _, err := fmt.Fscanf(r, templateMetrics,
		&s.Connections.Active,
		&s.Connections.Accepted,
		&s.Connections.Handled,
		&s.Requests,
		&s.Connections.Reading,
		&s.Connections.Writing,
		&s.Connections.Waiting); err != nil {
		return nil, fmt.Errorf("failed to scan template metrics: %w", err)
	}
	return &s, nil
}
