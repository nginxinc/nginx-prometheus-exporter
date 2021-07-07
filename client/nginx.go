package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// NginxClient allows you to fetch NGINX metrics from the stub_status page.
type NginxClient struct {
	apiEndpoint string
	httpClient  *http.Client
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
func NewNginxClient(httpClient *http.Client, apiEndpoint string) (*NginxClient, error) {
	client := &NginxClient{
		apiEndpoint: apiEndpoint,
		httpClient:  httpClient,
	}

	_, err := client.GetStubStats()
	return client, err
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

	var stats StubStats
	err = parseStubStats(body, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body %q: %w", string(body), err)
	}

	return &stats, nil
}

func parseStubStats(data []byte, stats *StubStats) error {
	if bytes.Count(data, []byte("\n")) != 4 {
		return fmt.Errorf("invalid input %s", data)
	}

	cnt, value := 0, int64(0)
	for _, s := range data {
		if s >= '0' && s <= '9' {
			value *= 10
			value += int64(s & 0x0f)
			continue
		}

		if value > 0 {
			switch cnt {
			case 0:
				stats.Connections.Active = value
			case 1:
				stats.Connections.Accepted = value
			case 2:
				stats.Connections.Handled = value
			case 3:
				stats.Requests = value
			case 4:
				stats.Connections.Reading = value
			case 5:
				stats.Connections.Writing = value
			case 6:
				stats.Connections.Waiting = value
			}
			value = 0
			cnt++
		}
	}

	if cnt != 7 {
		return fmt.Errorf("invalid input %s", data)
	}

	return nil
}
