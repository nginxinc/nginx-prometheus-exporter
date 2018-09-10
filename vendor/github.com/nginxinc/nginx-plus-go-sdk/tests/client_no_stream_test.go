package tests

import (
	"net/http"
	"testing"

	"github.com/nginxinc/nginx-plus-go-sdk/client"
)

// TestStatsNoStream tests the peculiar behavior of getting Stream-related
// stats from the API when there are no stream blocks in the config.
// The API returns a special error code that we can use to determine if the API
// is misconfigured or of the stream block is missing.
func TestStatsNoStream(t *testing.T) {
	httpClient := &http.Client{}
	c, err := client.NewNginxClient(httpClient, "http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("Error connecting to nginx: %v", err)
	}

	stats, err := c.GetStats()
	if err != nil {
		t.Errorf("Error getting stats: %v", err)
	}

	if stats.Connections.Accepted < 1 {
		t.Errorf("Stats should report some connections: %v", stats.Connections)
	}

	if len(stats.StreamServerZones) != 0 {
		t.Error("No stream block should result in no StreamServerZones")
	}

	if len(stats.StreamUpstreams) != 0 {
		t.Error("No stream block should result in no StreamUpstreams")
	}
}
