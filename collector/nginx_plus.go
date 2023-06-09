package collector

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	plusclient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/prometheus/client_golang/prometheus"
)

// LabelUpdater updates the labels of upstream server and server zone metrics
type LabelUpdater interface {
	UpdateUpstreamServerPeerLabels(upstreamServerPeerLabels map[string][]string)
	DeleteUpstreamServerPeerLabels(peers []string)
	UpdateUpstreamServerLabels(upstreamServerLabelValues map[string][]string)
	DeleteUpstreamServerLabels(upstreamNames []string)
	UpdateStreamUpstreamServerPeerLabels(streamUpstreamServerPeerLabels map[string][]string)
	DeleteStreamUpstreamServerPeerLabels(peers []string)
	UpdateStreamUpstreamServerLabels(streamUpstreamServerPeerLabels map[string][]string)
	DeleteStreamUpstreamServerLabels(peers []string)
	UpdateServerZoneLabels(serverZoneLabelValues map[string][]string)
	DeleteServerZoneLabels(zoneNames []string)
	UpdateStreamServerZoneLabels(streamServerZoneLabelValues map[string][]string)
	DeleteStreamServerZoneLabels(zoneNames []string)
}

// NginxPlusCollector collects NGINX Plus metrics. It implements prometheus.Collector interface.
type NginxPlusCollector struct {
	nginxClient                    *plusclient.NginxClient
	totalMetrics                   map[string]*prometheus.Desc
	serverZoneMetrics              map[string]*prometheus.Desc
	upstreamMetrics                map[string]*prometheus.Desc
	upstreamServerMetrics          map[string]*prometheus.Desc
	streamServerZoneMetrics        map[string]*prometheus.Desc
	streamZoneSyncMetrics          map[string]*prometheus.Desc
	streamUpstreamMetrics          map[string]*prometheus.Desc
	streamUpstreamServerMetrics    map[string]*prometheus.Desc
	locationZoneMetrics            map[string]*prometheus.Desc
	resolverMetrics                map[string]*prometheus.Desc
	limitRequestMetrics            map[string]*prometheus.Desc
	limitConnectionMetrics         map[string]*prometheus.Desc
	streamLimitConnectionMetrics   map[string]*prometheus.Desc
	upMetric                       prometheus.Gauge
	mutex                          sync.Mutex
	variableLabelNames             VariableLabelNames
	upstreamServerLabels           map[string][]string
	streamUpstreamServerLabels     map[string][]string
	serverZoneLabels               map[string][]string
	streamServerZoneLabels         map[string][]string
	upstreamServerPeerLabels       map[string][]string
	streamUpstreamServerPeerLabels map[string][]string
	variableLabelsMutex            sync.RWMutex
	logger                         log.Logger
}

// UpdateUpstreamServerPeerLabels updates the Upstream Server Peer Labels
func (c *NginxPlusCollector) UpdateUpstreamServerPeerLabels(upstreamServerPeerLabels map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range upstreamServerPeerLabels {
		c.upstreamServerPeerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteUpstreamServerPeerLabels deletes the Upstream Server Peer Labels
func (c *NginxPlusCollector) DeleteUpstreamServerPeerLabels(peers []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range peers {
		delete(c.upstreamServerPeerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateStreamUpstreamServerPeerLabels updates the Upstream Server Peer Labels
func (c *NginxPlusCollector) UpdateStreamUpstreamServerPeerLabels(streamUpstreamServerPeerLabels map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range streamUpstreamServerPeerLabels {
		c.streamUpstreamServerPeerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteStreamUpstreamServerPeerLabels deletes the Upstream Server Peer Labels
func (c *NginxPlusCollector) DeleteStreamUpstreamServerPeerLabels(peers []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range peers {
		delete(c.streamUpstreamServerPeerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateUpstreamServerLabels updates the Upstream Server Labels
func (c *NginxPlusCollector) UpdateUpstreamServerLabels(upstreamServerLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range upstreamServerLabelValues {
		c.upstreamServerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteUpstreamServerLabels deletes the Upstream Server Labels
func (c *NginxPlusCollector) DeleteUpstreamServerLabels(upstreamNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range upstreamNames {
		delete(c.upstreamServerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateStreamUpstreamServerLabels updates the Upstream Server Labels
func (c *NginxPlusCollector) UpdateStreamUpstreamServerLabels(streamUpstreamServerLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range streamUpstreamServerLabelValues {
		c.streamUpstreamServerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteStreamUpstreamServerLabels deletes the Upstream Server Labels
func (c *NginxPlusCollector) DeleteStreamUpstreamServerLabels(streamUpstreamNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range streamUpstreamNames {
		delete(c.streamUpstreamServerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateServerZoneLabels updates the Server Zone Labels
func (c *NginxPlusCollector) UpdateServerZoneLabels(serverZoneLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range serverZoneLabelValues {
		c.serverZoneLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteServerZoneLabels deletes the Server Zone Labels
func (c *NginxPlusCollector) DeleteServerZoneLabels(zoneNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range zoneNames {
		delete(c.serverZoneLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateStreamServerZoneLabels updates the Stream Server Zone Labels
func (c *NginxPlusCollector) UpdateStreamServerZoneLabels(streamServerZoneLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range streamServerZoneLabelValues {
		c.streamServerZoneLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteStreamServerZoneLabels deletes the Stream Server Zone Labels
func (c *NginxPlusCollector) DeleteStreamServerZoneLabels(zoneNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range zoneNames {
		delete(c.streamServerZoneLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

func (c *NginxPlusCollector) getUpstreamServerLabelValues(upstreamName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.upstreamServerLabels[upstreamName]
}

func (c *NginxPlusCollector) getStreamUpstreamServerLabelValues(upstreamName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.streamUpstreamServerLabels[upstreamName]
}

func (c *NginxPlusCollector) getServerZoneLabelValues(zoneName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.serverZoneLabels[zoneName]
}

func (c *NginxPlusCollector) getStreamServerZoneLabelValues(zoneName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.streamServerZoneLabels[zoneName]
}

func (c *NginxPlusCollector) getUpstreamServerPeerLabelValues(peer string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.upstreamServerPeerLabels[peer]
}

func (c *NginxPlusCollector) getStreamUpstreamServerPeerLabelValues(peer string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.streamUpstreamServerPeerLabels[peer]
}

// VariableLabelNames holds all the variable label names for the different metrics
type VariableLabelNames struct {
	UpstreamServerVariableLabelNames           []string
	ServerZoneVariableLabelNames               []string
	UpstreamServerPeerVariableLabelNames       []string
	StreamUpstreamServerPeerVariableLabelNames []string
	StreamServerZoneVariableLabelNames         []string
	StreamUpstreamServerVariableLabelNames     []string
}

// NewVariableLabels creates a new struct for VariableNames for the collector
func NewVariableLabelNames(upstreamServerVariableLabelNames []string, serverZoneVariableLabelNames []string, upstreamServerPeerVariableLabelNames []string,
	streamUpstreamServerVariableLabelNames []string, streamServerZoneLabels []string, streamUpstreamServerPeerVariableLabelNames []string,
) VariableLabelNames {
	return VariableLabelNames{
		UpstreamServerVariableLabelNames:           upstreamServerVariableLabelNames,
		ServerZoneVariableLabelNames:               serverZoneVariableLabelNames,
		UpstreamServerPeerVariableLabelNames:       upstreamServerPeerVariableLabelNames,
		StreamUpstreamServerVariableLabelNames:     streamUpstreamServerVariableLabelNames,
		StreamServerZoneVariableLabelNames:         streamServerZoneLabels,
		StreamUpstreamServerPeerVariableLabelNames: streamUpstreamServerPeerVariableLabelNames,
	}
}

// NewNginxPlusCollector creates an NginxPlusCollector.
func NewNginxPlusCollector(nginxClient *plusclient.NginxClient, namespace string, variableLabelNames VariableLabelNames, constLabels map[string]string, logger log.Logger) *NginxPlusCollector {
	upstreamServerVariableLabelNames := append(variableLabelNames.UpstreamServerVariableLabelNames, variableLabelNames.UpstreamServerPeerVariableLabelNames...)
	streamUpstreamServerVariableLabelNames := append(variableLabelNames.StreamUpstreamServerVariableLabelNames, variableLabelNames.StreamUpstreamServerPeerVariableLabelNames...)
	return &NginxPlusCollector{
		variableLabelNames:             variableLabelNames,
		upstreamServerLabels:           make(map[string][]string),
		serverZoneLabels:               make(map[string][]string),
		streamServerZoneLabels:         make(map[string][]string),
		upstreamServerPeerLabels:       make(map[string][]string),
		streamUpstreamServerPeerLabels: make(map[string][]string),
		streamUpstreamServerLabels:     make(map[string][]string),
		nginxClient:                    nginxClient,
		logger:                         logger,
		totalMetrics: map[string]*prometheus.Desc{
			"connections_accepted":  newGlobalMetric(namespace, "connections_accepted", "Accepted client connections", constLabels),
			"connections_dropped":   newGlobalMetric(namespace, "connections_dropped", "Dropped client connections", constLabels),
			"connections_active":    newGlobalMetric(namespace, "connections_active", "Active client connections", constLabels),
			"connections_idle":      newGlobalMetric(namespace, "connections_idle", "Idle client connections", constLabels),
			"http_requests_total":   newGlobalMetric(namespace, "http_requests_total", "Total http requests", constLabels),
			"http_requests_current": newGlobalMetric(namespace, "http_requests_current", "Current http requests", constLabels),
			"ssl_handshakes":        newGlobalMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", constLabels),
			"ssl_handshakes_failed": newGlobalMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", constLabels),
			"ssl_session_reuses":    newGlobalMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", constLabels),
		},
		serverZoneMetrics: map[string]*prometheus.Desc{
			"processing":            newServerZoneMetric(namespace, "processing", "Client requests that are currently being processed", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"requests":              newServerZoneMetric(namespace, "requests", "Total client requests", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"responses_1xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "1xx"})),
			"responses_2xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"responses_3xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "3xx"})),
			"responses_4xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"responses_5xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"discarded":             newServerZoneMetric(namespace, "discarded", "Requests completed without sending a response", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"received":              newServerZoneMetric(namespace, "received", "Bytes received from clients", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"sent":                  newServerZoneMetric(namespace, "sent", "Bytes sent to clients", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"codes_100":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "100"})),
			"codes_101":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "101"})),
			"codes_102":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "102"})),
			"codes_200":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "200"})),
			"codes_201":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "201"})),
			"codes_202":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "202"})),
			"codes_204":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "204"})),
			"codes_206":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "206"})),
			"codes_300":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "300"})),
			"codes_301":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "301"})),
			"codes_302":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "302"})),
			"codes_303":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "303"})),
			"codes_304":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "304"})),
			"codes_307":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "307"})),
			"codes_400":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "400"})),
			"codes_401":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "401"})),
			"codes_403":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "403"})),
			"codes_404":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "404"})),
			"codes_405":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "405"})),
			"codes_408":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "408"})),
			"codes_409":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "409"})),
			"codes_411":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "411"})),
			"codes_412":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "412"})),
			"codes_413":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "413"})),
			"codes_414":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "414"})),
			"codes_415":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "415"})),
			"codes_416":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "416"})),
			"codes_429":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "429"})),
			"codes_444":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "444"})),
			"codes_494":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "494"})),
			"codes_495":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "495"})),
			"codes_496":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "496"})),
			"codes_497":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "497"})),
			"codes_499":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "499"})),
			"codes_500":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "500"})),
			"codes_501":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "501"})),
			"codes_502":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "502"})),
			"codes_503":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "503"})),
			"codes_504":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "504"})),
			"codes_507":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "507"})),
			"ssl_handshakes":        newServerZoneMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"ssl_handshakes_failed": newServerZoneMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"ssl_session_reuses":    newServerZoneMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
		},
		streamServerZoneMetrics: map[string]*prometheus.Desc{
			"processing":            newStreamServerZoneMetric(namespace, "processing", "Client connections that are currently being processed", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"connections":           newStreamServerZoneMetric(namespace, "connections", "Total connections", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"sessions_2xx":          newStreamServerZoneMetric(namespace, "sessions", "Total sessions completed", variableLabelNames.StreamServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"sessions_4xx":          newStreamServerZoneMetric(namespace, "sessions", "Total sessions completed", variableLabelNames.StreamServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"sessions_5xx":          newStreamServerZoneMetric(namespace, "sessions", "Total sessions completed", variableLabelNames.StreamServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"discarded":             newStreamServerZoneMetric(namespace, "discarded", "Connections completed without creating a session", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"received":              newStreamServerZoneMetric(namespace, "received", "Bytes received from clients", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"sent":                  newStreamServerZoneMetric(namespace, "sent", "Bytes sent to clients", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"ssl_handshakes":        newStreamServerZoneMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"ssl_handshakes_failed": newStreamServerZoneMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"ssl_session_reuses":    newStreamServerZoneMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
		},
		upstreamMetrics: map[string]*prometheus.Desc{
			"keepalives": newUpstreamMetric(namespace, "keepalives", "Idle keepalive connections", constLabels),
			"zombies":    newUpstreamMetric(namespace, "zombies", "Servers removed from the group but still processing active client requests", constLabels),
		},
		streamUpstreamMetrics: map[string]*prometheus.Desc{
			"zombies": newStreamUpstreamMetric(namespace, "zombies", "Servers removed from the group but still processing active client connections", constLabels),
		},
		upstreamServerMetrics: map[string]*prometheus.Desc{
			"state":                   newUpstreamServerMetric(namespace, "state", "Current state", upstreamServerVariableLabelNames, constLabels),
			"active":                  newUpstreamServerMetric(namespace, "active", "Active connections", upstreamServerVariableLabelNames, constLabels),
			"limit":                   newUpstreamServerMetric(namespace, "limit", "Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit", upstreamServerVariableLabelNames, constLabels),
			"requests":                newUpstreamServerMetric(namespace, "requests", "Total client requests", upstreamServerVariableLabelNames, constLabels),
			"responses_1xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "1xx"})),
			"responses_2xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"responses_3xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "3xx"})),
			"responses_4xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"responses_5xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"sent":                    newUpstreamServerMetric(namespace, "sent", "Bytes sent to this server", upstreamServerVariableLabelNames, constLabels),
			"received":                newUpstreamServerMetric(namespace, "received", "Bytes received to this server", upstreamServerVariableLabelNames, constLabels),
			"fails":                   newUpstreamServerMetric(namespace, "fails", "Number of unsuccessful attempts to communicate with the server", upstreamServerVariableLabelNames, constLabels),
			"unavail":                 newUpstreamServerMetric(namespace, "unavail", "How many times the server became unavailable for client requests (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold", upstreamServerVariableLabelNames, constLabels),
			"header_time":             newUpstreamServerMetric(namespace, "header_time", "Average time to get the response header from the server", upstreamServerVariableLabelNames, constLabels),
			"response_time":           newUpstreamServerMetric(namespace, "response_time", "Average time to get the full response from the server", upstreamServerVariableLabelNames, constLabels),
			"health_checks_checks":    newUpstreamServerMetric(namespace, "health_checks_checks", "Total health check requests", upstreamServerVariableLabelNames, constLabels),
			"health_checks_fails":     newUpstreamServerMetric(namespace, "health_checks_fails", "Failed health checks", upstreamServerVariableLabelNames, constLabels),
			"health_checks_unhealthy": newUpstreamServerMetric(namespace, "health_checks_unhealthy", "How many times the server became unhealthy (state 'unhealthy')", upstreamServerVariableLabelNames, constLabels),
			"codes_100":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "100"})),
			"codes_101":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "101"})),
			"codes_102":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "102"})),
			"codes_200":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "200"})),
			"codes_201":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "201"})),
			"codes_202":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "202"})),
			"codes_204":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "204"})),
			"codes_206":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "206"})),
			"codes_300":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "300"})),
			"codes_301":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "301"})),
			"codes_302":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "302"})),
			"codes_303":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "303"})),
			"codes_304":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "304"})),
			"codes_307":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "307"})),
			"codes_400":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "400"})),
			"codes_401":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "401"})),
			"codes_403":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "403"})),
			"codes_404":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "404"})),
			"codes_405":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "405"})),
			"codes_408":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "408"})),
			"codes_409":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "409"})),
			"codes_411":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "411"})),
			"codes_412":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "412"})),
			"codes_413":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "413"})),
			"codes_414":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "414"})),
			"codes_415":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "415"})),
			"codes_416":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "416"})),
			"codes_429":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "429"})),
			"codes_444":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "444"})),
			"codes_494":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "494"})),
			"codes_495":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "495"})),
			"codes_496":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "496"})),
			"codes_497":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "497"})),
			"codes_499":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "499"})),
			"codes_500":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "500"})),
			"codes_501":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "501"})),
			"codes_502":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "502"})),
			"codes_503":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "503"})),
			"codes_504":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "504"})),
			"codes_507":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "507"})),
			"ssl_handshakes":          newUpstreamServerMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", upstreamServerVariableLabelNames, constLabels),
			"ssl_handshakes_failed":   newUpstreamServerMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", upstreamServerVariableLabelNames, constLabels),
			"ssl_session_reuses":      newUpstreamServerMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", upstreamServerVariableLabelNames, constLabels),
		},
		streamUpstreamServerMetrics: map[string]*prometheus.Desc{
			"state":                   newStreamUpstreamServerMetric(namespace, "state", "Current state", streamUpstreamServerVariableLabelNames, constLabels),
			"active":                  newStreamUpstreamServerMetric(namespace, "active", "Active connections", streamUpstreamServerVariableLabelNames, constLabels),
			"limit":                   newStreamUpstreamServerMetric(namespace, "limit", "Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit", streamUpstreamServerVariableLabelNames, constLabels),
			"sent":                    newStreamUpstreamServerMetric(namespace, "sent", "Bytes sent to this server", streamUpstreamServerVariableLabelNames, constLabels),
			"received":                newStreamUpstreamServerMetric(namespace, "received", "Bytes received from this server", streamUpstreamServerVariableLabelNames, constLabels),
			"fails":                   newStreamUpstreamServerMetric(namespace, "fails", "Number of unsuccessful attempts to communicate with the server", streamUpstreamServerVariableLabelNames, constLabels),
			"unavail":                 newStreamUpstreamServerMetric(namespace, "unavail", "How many times the server became unavailable for client connections (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold", streamUpstreamServerVariableLabelNames, constLabels),
			"connections":             newStreamUpstreamServerMetric(namespace, "connections", "Total number of client connections forwarded to this server", streamUpstreamServerVariableLabelNames, constLabels),
			"connect_time":            newStreamUpstreamServerMetric(namespace, "connect_time", "Average time to connect to the upstream server", streamUpstreamServerVariableLabelNames, constLabels),
			"first_byte_time":         newStreamUpstreamServerMetric(namespace, "first_byte_time", "Average time to receive the first byte of data", streamUpstreamServerVariableLabelNames, constLabels),
			"response_time":           newStreamUpstreamServerMetric(namespace, "response_time", "Average time to receive the last byte of data", streamUpstreamServerVariableLabelNames, constLabels),
			"health_checks_checks":    newStreamUpstreamServerMetric(namespace, "health_checks_checks", "Total health check requests", streamUpstreamServerVariableLabelNames, constLabels),
			"health_checks_fails":     newStreamUpstreamServerMetric(namespace, "health_checks_fails", "Failed health checks", streamUpstreamServerVariableLabelNames, constLabels),
			"health_checks_unhealthy": newStreamUpstreamServerMetric(namespace, "health_checks_unhealthy", "How many times the server became unhealthy (state 'unhealthy')", streamUpstreamServerVariableLabelNames, constLabels),
			"ssl_handshakes":          newStreamUpstreamServerMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", streamUpstreamServerVariableLabelNames, constLabels),
			"ssl_handshakes_failed":   newStreamUpstreamServerMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", streamUpstreamServerVariableLabelNames, constLabels),
			"ssl_session_reuses":      newStreamUpstreamServerMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", streamUpstreamServerVariableLabelNames, constLabels),
		},
		streamZoneSyncMetrics: map[string]*prometheus.Desc{
			"bytes_in":        newStreamZoneSyncMetric(namespace, "bytes_in", "Bytes received by this node", constLabels),
			"bytes_out":       newStreamZoneSyncMetric(namespace, "bytes_out", "Bytes sent by this node", constLabels),
			"msgs_in":         newStreamZoneSyncMetric(namespace, "msgs_in", "Total messages received by this node", constLabels),
			"msgs_out":        newStreamZoneSyncMetric(namespace, "msgs_out", "Total messages sent by this node", constLabels),
			"nodes_online":    newStreamZoneSyncMetric(namespace, "nodes_online", "Number of peers this node is connected to", constLabels),
			"records_pending": newStreamZoneSyncZoneMetric(namespace, "records_pending", "The number of records that need to be sent to the cluster", constLabels),
			"records_total":   newStreamZoneSyncZoneMetric(namespace, "records_total", "The total number of records stored in the shared memory zone", constLabels),
		},
		locationZoneMetrics: map[string]*prometheus.Desc{
			"requests":      newLocationZoneMetric(namespace, "requests", "Total client requests", constLabels),
			"responses_1xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "1xx"})),
			"responses_2xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"responses_3xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "3xx"})),
			"responses_4xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"responses_5xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"discarded":     newLocationZoneMetric(namespace, "discarded", "Requests completed without sending a response", constLabels),
			"received":      newLocationZoneMetric(namespace, "received", "Bytes received from clients", constLabels),
			"sent":          newLocationZoneMetric(namespace, "sent", "Bytes sent to clients", constLabels),
			"codes_100":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "100"})),
			"codes_101":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "101"})),
			"codes_102":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "102"})),
			"codes_200":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "200"})),
			"codes_201":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "201"})),
			"codes_202":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "202"})),
			"codes_204":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "204"})),
			"codes_206":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "206"})),
			"codes_300":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "300"})),
			"codes_301":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "301"})),
			"codes_302":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "302"})),
			"codes_303":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "303"})),
			"codes_304":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "304"})),
			"codes_307":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "307"})),
			"codes_400":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "400"})),
			"codes_401":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "401"})),
			"codes_403":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "403"})),
			"codes_404":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "404"})),
			"codes_405":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "405"})),
			"codes_408":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "408"})),
			"codes_409":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "409"})),
			"codes_411":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "411"})),
			"codes_412":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "412"})),
			"codes_413":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "413"})),
			"codes_414":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "414"})),
			"codes_415":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "415"})),
			"codes_416":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "416"})),
			"codes_429":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "429"})),
			"codes_444":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "444"})),
			"codes_494":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "494"})),
			"codes_495":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "495"})),
			"codes_496":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "496"})),
			"codes_497":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "497"})),
			"codes_499":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "499"})),
			"codes_500":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "500"})),
			"codes_501":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "501"})),
			"codes_502":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "502"})),
			"codes_503":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "503"})),
			"codes_504":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "504"})),
			"codes_507":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "507"})),
		},
		resolverMetrics: map[string]*prometheus.Desc{
			"name":     newResolverMetric(namespace, "name", "Total requests to resolve names to addresses", constLabels),
			"srv":      newResolverMetric(namespace, "srv", "Total requests to resolve SRV records", constLabels),
			"addr":     newResolverMetric(namespace, "addr", "Total requests to resolve addresses to names", constLabels),
			"noerror":  newResolverMetric(namespace, "noerror", "Total number of successful responses", constLabels),
			"formerr":  newResolverMetric(namespace, "formerr", "Total number of FORMERR responses", constLabels),
			"servfail": newResolverMetric(namespace, "servfail", "Total number of SERVFAIL responses", constLabels),
			"nxdomain": newResolverMetric(namespace, "nxdomain", "Total number of NXDOMAIN responses", constLabels),
			"notimp":   newResolverMetric(namespace, "notimp", "Total number of NOTIMP responses", constLabels),
			"refused":  newResolverMetric(namespace, "refused", "Total number of REFUSED responses", constLabels),
			"timedout": newResolverMetric(namespace, "timedout", "Total number of timed out requests", constLabels),
			"unknown":  newResolverMetric(namespace, "unknown", "Total requests completed with an unknown error", constLabels),
		},
		limitRequestMetrics: map[string]*prometheus.Desc{
			"passed":           newLimitRequestMetric(namespace, "passed", "Total number of requests that were neither limited nor accounted as limited", constLabels),
			"delayed":          newLimitRequestMetric(namespace, "delayed", "Total number of requests that were delayed", constLabels),
			"rejected":         newLimitRequestMetric(namespace, "rejected", "Total number of requests that were rejected", constLabels),
			"delayed_dry_run":  newLimitRequestMetric(namespace, "delayed_dry_run", "Total number of requests accounted as delayed in the dry run mode", constLabels),
			"rejected_dry_run": newLimitRequestMetric(namespace, "rejected_dry_run", "Total number of requests accounted as rejected in the dry run mode", constLabels),
		},
		limitConnectionMetrics: map[string]*prometheus.Desc{
			"passed":           newLimitConnectionMetric(namespace, "passed", "Total number of connections that were neither limited nor accounted as limited", constLabels),
			"rejected":         newLimitConnectionMetric(namespace, "rejected", "Total number of connections that were rejected", constLabels),
			"rejected_dry_run": newLimitConnectionMetric(namespace, "rejected_dry_run", "Total number of connections accounted as rejected in the dry run mode", constLabels),
		},
		streamLimitConnectionMetrics: map[string]*prometheus.Desc{
			"passed":           newStreamLimitConnectionMetric(namespace, "passed", "Total number of connections that were neither limited nor accounted as limited", constLabels),
			"rejected":         newStreamLimitConnectionMetric(namespace, "rejected", "Total number of connections that were rejected", constLabels),
			"rejected_dry_run": newStreamLimitConnectionMetric(namespace, "rejected_dry_run", "Total number of connections accounted as rejected in the dry run mode", constLabels),
		},
		upMetric: newUpMetric(namespace, constLabels),
	}
}

// Describe sends the super-set of all possible descriptors of NGINX Plus metrics
// to the provided channel.
func (c *NginxPlusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upMetric.Desc()

	for _, m := range c.totalMetrics {
		ch <- m
	}
	for _, m := range c.serverZoneMetrics {
		ch <- m
	}
	for _, m := range c.upstreamMetrics {
		ch <- m
	}
	for _, m := range c.upstreamServerMetrics {
		ch <- m
	}
	for _, m := range c.streamServerZoneMetrics {
		ch <- m
	}
	for _, m := range c.streamUpstreamMetrics {
		ch <- m
	}
	for _, m := range c.streamUpstreamServerMetrics {
		ch <- m
	}
	for _, m := range c.streamZoneSyncMetrics {
		ch <- m
	}
	for _, m := range c.locationZoneMetrics {
		ch <- m
	}
	for _, m := range c.resolverMetrics {
		ch <- m
	}
	for _, m := range c.limitRequestMetrics {
		ch <- m
	}
	for _, m := range c.limitConnectionMetrics {
		ch <- m
	}
	for _, m := range c.streamLimitConnectionMetrics {
		ch <- m
	}
}

// Collect fetches metrics from NGINX Plus and sends them to the provided channel.
func (c *NginxPlusCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	stats, err := c.nginxClient.GetStats()
	if err != nil {
		c.upMetric.Set(nginxDown)
		ch <- c.upMetric
		level.Warn(c.logger).Log("msg", "Error getting stats", "error", err.Error())
		return
	}

	c.upMetric.Set(nginxUp)
	ch <- c.upMetric

	ch <- prometheus.MustNewConstMetric(c.totalMetrics["connections_accepted"],
		prometheus.CounterValue, float64(stats.Connections.Accepted))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["connections_dropped"],
		prometheus.CounterValue, float64(stats.Connections.Dropped))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["connections_active"],
		prometheus.GaugeValue, float64(stats.Connections.Active))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["connections_idle"],
		prometheus.GaugeValue, float64(stats.Connections.Idle))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["http_requests_total"],
		prometheus.CounterValue, float64(stats.HTTPRequests.Total))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["http_requests_current"],
		prometheus.GaugeValue, float64(stats.HTTPRequests.Current))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["ssl_handshakes"],
		prometheus.CounterValue, float64(stats.SSL.Handshakes))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["ssl_handshakes_failed"],
		prometheus.CounterValue, float64(stats.SSL.HandshakesFailed))
	ch <- prometheus.MustNewConstMetric(c.totalMetrics["ssl_session_reuses"],
		prometheus.CounterValue, float64(stats.SSL.SessionReuses))

	for name, zone := range stats.ServerZones {
		labelValues := []string{name}
		varLabelValues := c.getServerZoneLabelValues(name)

		if c.variableLabelNames.ServerZoneVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.ServerZoneVariableLabelNames) {
			level.Warn(c.logger).Log("msg", "wrong number of labels for http zone, empty labels will be used instead", "zone", name, "expected", len(c.variableLabelNames.ServerZoneVariableLabelNames), "got", len(varLabelValues))
			for range c.variableLabelNames.ServerZoneVariableLabelNames {
				labelValues = append(labelValues, "")
			}
		} else {
			labelValues = append(labelValues, varLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["processing"],
			prometheus.GaugeValue, float64(zone.Processing), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["requests"],
			prometheus.CounterValue, float64(zone.Requests), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["responses_1xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses1xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["responses_2xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses2xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["responses_3xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses3xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["responses_4xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses4xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["responses_5xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses5xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["discarded"],
			prometheus.CounterValue, float64(zone.Discarded), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["received"],
			prometheus.CounterValue, float64(zone.Received), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["sent"],
			prometheus.CounterValue, float64(zone.Sent), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_100"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPContinue), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_101"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSwitchingProtocols), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_102"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPProcessing), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_200"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPOk), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_201"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPCreated), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_202"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPAccepted), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_204"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNoContent), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_206"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPPartialContent), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_300"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSpecialResponse), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_301"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPMovedPermanently), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_302"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPMovedTemporarily), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_303"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSeeOther), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_304"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotModified), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_307"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPTemporaryRedirect), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_400"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPBadRequest), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_401"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPUnauthorized), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_403"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPForbidden), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_404"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotFound), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_405"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotAllowed), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_408"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestTimeOut), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_409"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPConflict), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_411"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPLengthRequired), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_412"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPPreconditionFailed), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_413"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestEntityTooLarge), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_414"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestURITooLarge), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_415"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPUnsupportedMediaType), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_416"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRangeNotSatisfiable), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_429"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPTooManyRequests), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_444"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPClose), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_494"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestHeaderTooLarge), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_495"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSCertError), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_496"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSNoCert), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_497"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPToHTTPS), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_499"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPClientClosedRequest), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_500"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPInternalServerError), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_501"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotImplemented), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_502"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPBadGateway), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_503"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPServiceUnavailable), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_504"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPGatewayTimeOut), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["codes_507"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPInsufficientStorage), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["ssl_handshakes"],
			prometheus.CounterValue, float64(zone.SSL.Handshakes), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["ssl_handshakes_failed"],
			prometheus.CounterValue, float64(zone.SSL.HandshakesFailed), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.serverZoneMetrics["ssl_session_reuses"],
			prometheus.CounterValue, float64(zone.SSL.SessionReuses), labelValues...)
	}

	for name, zone := range stats.StreamServerZones {
		labelValues := []string{name}
		varLabelValues := c.getStreamServerZoneLabelValues(name)

		if c.variableLabelNames.StreamServerZoneVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.StreamServerZoneVariableLabelNames) {
			level.Warn(c.logger).Log("msg", "wrong number of labels for stream server zone, empty labels will be used instead", "zone", name, "expected", len(c.variableLabelNames.StreamServerZoneVariableLabelNames), "got", len(varLabelValues))
			for range c.variableLabelNames.StreamServerZoneVariableLabelNames {
				labelValues = append(labelValues, "")
			}
		} else {
			labelValues = append(labelValues, varLabelValues...)
		}
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["processing"],
			prometheus.GaugeValue, float64(zone.Processing), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["connections"],
			prometheus.CounterValue, float64(zone.Connections), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["sessions_2xx"],
			prometheus.CounterValue, float64(zone.Sessions.Sessions2xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["sessions_4xx"],
			prometheus.CounterValue, float64(zone.Sessions.Sessions4xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["sessions_5xx"],
			prometheus.CounterValue, float64(zone.Sessions.Sessions5xx), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["discarded"],
			prometheus.CounterValue, float64(zone.Discarded), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["received"],
			prometheus.CounterValue, float64(zone.Received), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["sent"],
			prometheus.CounterValue, float64(zone.Sent), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["ssl_handshakes"],
			prometheus.CounterValue, float64(zone.SSL.Handshakes), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["ssl_handshakes_failed"],
			prometheus.CounterValue, float64(zone.SSL.HandshakesFailed), labelValues...)
		ch <- prometheus.MustNewConstMetric(c.streamServerZoneMetrics["ssl_session_reuses"],
			prometheus.CounterValue, float64(zone.SSL.SessionReuses), labelValues...)
	}

	for name, upstream := range stats.Upstreams {
		for _, peer := range upstream.Peers {
			labelValues := []string{name, peer.Server}
			varLabelValues := c.getUpstreamServerLabelValues(name)

			if c.variableLabelNames.UpstreamServerVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.UpstreamServerVariableLabelNames) {
				level.Warn(c.logger).Log("msg", "wrong number of labels for upstream, empty labels will be used instead", "upstream", name, "expected", len(c.variableLabelNames.UpstreamServerVariableLabelNames), "got", len(varLabelValues))
				for range c.variableLabelNames.UpstreamServerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varLabelValues...)
			}

			upstreamServer := fmt.Sprintf("%v/%v", name, peer.Server)
			varPeerLabelValues := c.getUpstreamServerPeerLabelValues(upstreamServer)
			if c.variableLabelNames.UpstreamServerPeerVariableLabelNames != nil && len(varPeerLabelValues) != len(c.variableLabelNames.UpstreamServerPeerVariableLabelNames) {
				level.Warn(c.logger).Log("msg", "wrong number of labels for upstream peer, empty labels will be used instead", "upstream", name, "peer", peer.Server, "expected", len(c.variableLabelNames.UpstreamServerPeerVariableLabelNames), "got", len(varPeerLabelValues))
				for range c.variableLabelNames.UpstreamServerPeerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varPeerLabelValues...)
			}

			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["state"],
				prometheus.GaugeValue, upstreamServerStates[peer.State], labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["active"],
				prometheus.GaugeValue, float64(peer.Active), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["limit"],
				prometheus.GaugeValue, float64(peer.MaxConns), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["requests"],
				prometheus.CounterValue, float64(peer.Requests), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["responses_1xx"],
				prometheus.CounterValue, float64(peer.Responses.Responses1xx), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["responses_2xx"],
				prometheus.CounterValue, float64(peer.Responses.Responses2xx), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["responses_3xx"],
				prometheus.CounterValue, float64(peer.Responses.Responses3xx), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["responses_4xx"],
				prometheus.CounterValue, float64(peer.Responses.Responses4xx), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["responses_5xx"],
				prometheus.CounterValue, float64(peer.Responses.Responses5xx), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["sent"],
				prometheus.CounterValue, float64(peer.Sent), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["received"],
				prometheus.CounterValue, float64(peer.Received), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["fails"],
				prometheus.CounterValue, float64(peer.Fails), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["unavail"],
				prometheus.CounterValue, float64(peer.Unavail), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["header_time"],
				prometheus.GaugeValue, float64(peer.HeaderTime), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["response_time"],
				prometheus.GaugeValue, float64(peer.ResponseTime), labelValues...)

			if peer.HealthChecks != (plusclient.HealthChecks{}) {
				ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["health_checks_checks"],
					prometheus.CounterValue, float64(peer.HealthChecks.Checks), labelValues...)
				ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["health_checks_fails"],
					prometheus.CounterValue, float64(peer.HealthChecks.Fails), labelValues...)
				ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["health_checks_unhealthy"],
					prometheus.CounterValue, float64(peer.HealthChecks.Unhealthy), labelValues...)
			}
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_100"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPContinue), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_101"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPSwitchingProtocols), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_102"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPProcessing), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_200"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPOk), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_201"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPCreated), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_202"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPAccepted), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_204"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPNoContent), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_206"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPPartialContent), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_300"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPSpecialResponse), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_301"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPMovedPermanently), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_302"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPMovedTemporarily), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_303"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPSeeOther), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_304"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPNotModified), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_307"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPTemporaryRedirect), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_400"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPBadRequest), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_401"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPUnauthorized), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_403"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPForbidden), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_404"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPNotFound), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_405"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPNotAllowed), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_408"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPRequestTimeOut), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_409"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPConflict), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_411"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPLengthRequired), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_412"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPPreconditionFailed), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_413"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPRequestEntityTooLarge), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_414"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPRequestURITooLarge), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_415"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPUnsupportedMediaType), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_416"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPRangeNotSatisfiable), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_429"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPTooManyRequests), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_444"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPClose), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_494"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPRequestHeaderTooLarge), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_495"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPSCertError), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_496"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPSNoCert), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_497"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPToHTTPS), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_499"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPClientClosedRequest), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_500"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPInternalServerError), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_501"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPNotImplemented), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_502"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPBadGateway), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_503"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPServiceUnavailable), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_504"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPGatewayTimeOut), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["codes_507"],
				prometheus.CounterValue, float64(peer.Responses.Codes.HTTPInsufficientStorage), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["ssl_handshakes"],
				prometheus.CounterValue, float64(peer.SSL.Handshakes), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["ssl_handshakes_failed"],
				prometheus.CounterValue, float64(peer.SSL.HandshakesFailed), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.upstreamServerMetrics["ssl_session_reuses"],
				prometheus.CounterValue, float64(peer.SSL.SessionReuses), labelValues...)
		}
		ch <- prometheus.MustNewConstMetric(c.upstreamMetrics["keepalives"],
			prometheus.GaugeValue, float64(upstream.Keepalives), name)
		ch <- prometheus.MustNewConstMetric(c.upstreamMetrics["zombies"],
			prometheus.GaugeValue, float64(upstream.Zombies), name)
	}

	for name, upstream := range stats.StreamUpstreams {
		for _, peer := range upstream.Peers {
			labelValues := []string{name, peer.Server}
			varLabelValues := c.getStreamUpstreamServerLabelValues(name)

			if c.variableLabelNames.StreamUpstreamServerVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.StreamUpstreamServerVariableLabelNames) {
				level.Warn(c.logger).Log("msg", "wrong number of labels for stream server, empty labels will be used instead", "server", name, "labels", c.variableLabelNames.StreamUpstreamServerVariableLabelNames, "values", varLabelValues)
				for range c.variableLabelNames.StreamUpstreamServerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varLabelValues...)
			}

			upstreamServer := fmt.Sprintf("%v/%v", name, peer.Server)
			varPeerLabelValues := c.getStreamUpstreamServerPeerLabelValues(upstreamServer)
			if c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames != nil && len(varPeerLabelValues) != len(c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames) {
				level.Warn(c.logger).Log("msg", "wrong number of labels for stream upstream peer, empty labels will be used instead", "server", upstreamServer, "labels", c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames, "values", varPeerLabelValues)
				for range c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varPeerLabelValues...)
			}

			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["state"],
				prometheus.GaugeValue, upstreamServerStates[peer.State], labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["active"],
				prometheus.GaugeValue, float64(peer.Active), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["limit"],
				prometheus.GaugeValue, float64(peer.MaxConns), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["connections"],
				prometheus.CounterValue, float64(peer.Connections), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["connect_time"],
				prometheus.GaugeValue, float64(peer.ConnectTime), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["first_byte_time"],
				prometheus.GaugeValue, float64(peer.FirstByteTime), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["response_time"],
				prometheus.GaugeValue, float64(peer.ResponseTime), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["sent"],
				prometheus.CounterValue, float64(peer.Sent), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["received"],
				prometheus.CounterValue, float64(peer.Received), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["fails"],
				prometheus.CounterValue, float64(peer.Fails), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["unavail"],
				prometheus.CounterValue, float64(peer.Unavail), labelValues...)
			if peer.HealthChecks != (plusclient.HealthChecks{}) {
				ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["health_checks_checks"],
					prometheus.CounterValue, float64(peer.HealthChecks.Checks), labelValues...)
				ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["health_checks_fails"],
					prometheus.CounterValue, float64(peer.HealthChecks.Fails), labelValues...)
				ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["health_checks_unhealthy"],
					prometheus.CounterValue, float64(peer.HealthChecks.Unhealthy), labelValues...)
			}
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["ssl_handshakes"],
				prometheus.CounterValue, float64(peer.SSL.Handshakes), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["ssl_handshakes_failed"],
				prometheus.CounterValue, float64(peer.SSL.HandshakesFailed), labelValues...)
			ch <- prometheus.MustNewConstMetric(c.streamUpstreamServerMetrics["ssl_session_reuses"],
				prometheus.CounterValue, float64(peer.SSL.SessionReuses), labelValues...)
		}
		ch <- prometheus.MustNewConstMetric(c.streamUpstreamMetrics["zombies"],
			prometheus.GaugeValue, float64(upstream.Zombies), name)
	}

	if stats.StreamZoneSync != nil {
		for name, zone := range stats.StreamZoneSync.Zones {
			ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["records_pending"],
				prometheus.GaugeValue, float64(zone.RecordsPending), name)
			ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["records_total"],
				prometheus.GaugeValue, float64(zone.RecordsTotal), name)
		}

		ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["bytes_in"],
			prometheus.CounterValue, float64(stats.StreamZoneSync.Status.BytesIn))
		ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["bytes_out"],
			prometheus.CounterValue, float64(stats.StreamZoneSync.Status.BytesOut))
		ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["msgs_in"],
			prometheus.CounterValue, float64(stats.StreamZoneSync.Status.MsgsIn))
		ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["msgs_out"],
			prometheus.CounterValue, float64(stats.StreamZoneSync.Status.MsgsOut))
		ch <- prometheus.MustNewConstMetric(c.streamZoneSyncMetrics["nodes_online"],
			prometheus.GaugeValue, float64(stats.StreamZoneSync.Status.NodesOnline))
	}

	for name, zone := range stats.LocationZones {
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["requests"],
			prometheus.CounterValue, float64(zone.Requests), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["responses_1xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses1xx), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["responses_2xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses2xx), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["responses_3xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses3xx), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["responses_4xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses4xx), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["responses_5xx"],
			prometheus.CounterValue, float64(zone.Responses.Responses5xx), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["discarded"],
			prometheus.CounterValue, float64(zone.Discarded), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["received"],
			prometheus.CounterValue, float64(zone.Received), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["sent"],
			prometheus.CounterValue, float64(zone.Sent), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_100"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPContinue), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_101"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSwitchingProtocols), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_102"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPProcessing), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_200"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPOk), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_201"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPCreated), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_202"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPAccepted), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_204"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNoContent), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_206"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPPartialContent), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_300"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSpecialResponse), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_301"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPMovedPermanently), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_302"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPMovedTemporarily), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_303"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSeeOther), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_304"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotModified), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_307"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPTemporaryRedirect), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_400"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPBadRequest), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_401"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPUnauthorized), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_403"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPForbidden), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_404"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotFound), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_405"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotAllowed), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_408"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestTimeOut), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_409"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPConflict), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_411"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPLengthRequired), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_412"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPPreconditionFailed), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_413"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestEntityTooLarge), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_414"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestURITooLarge), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_415"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPUnsupportedMediaType), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_416"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRangeNotSatisfiable), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_429"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPTooManyRequests), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_444"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPClose), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_494"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPRequestHeaderTooLarge), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_495"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSCertError), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_496"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPSNoCert), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_497"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPToHTTPS), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_499"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPClientClosedRequest), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_500"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPInternalServerError), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_501"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPNotImplemented), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_502"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPBadGateway), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_503"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPServiceUnavailable), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_504"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPGatewayTimeOut), name)
		ch <- prometheus.MustNewConstMetric(c.locationZoneMetrics["codes_507"],
			prometheus.CounterValue, float64(zone.Responses.Codes.HTTPInsufficientStorage), name)
	}

	for name, zone := range stats.Resolvers {
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["name"],
			prometheus.CounterValue, float64(zone.Requests.Name), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["srv"],
			prometheus.CounterValue, float64(zone.Requests.Srv), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["addr"],
			prometheus.CounterValue, float64(zone.Requests.Addr), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["noerror"],
			prometheus.CounterValue, float64(zone.Responses.Noerror), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["formerr"],
			prometheus.CounterValue, float64(zone.Responses.Formerr), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["servfail"],
			prometheus.CounterValue, float64(zone.Responses.Servfail), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["nxdomain"],
			prometheus.CounterValue, float64(zone.Responses.Nxdomain), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["notimp"],
			prometheus.CounterValue, float64(zone.Responses.Notimp), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["refused"],
			prometheus.CounterValue, float64(zone.Responses.Refused), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["timedout"],
			prometheus.CounterValue, float64(zone.Responses.Timedout), name)
		ch <- prometheus.MustNewConstMetric(c.resolverMetrics["unknown"],
			prometheus.CounterValue, float64(zone.Responses.Unknown), name)
	}

	for name, zone := range stats.HTTPLimitRequests {
		ch <- prometheus.MustNewConstMetric(c.limitRequestMetrics["passed"], prometheus.CounterValue, float64(zone.Passed), name)
		ch <- prometheus.MustNewConstMetric(c.limitRequestMetrics["rejected"], prometheus.CounterValue, float64(zone.Rejected), name)
		ch <- prometheus.MustNewConstMetric(c.limitRequestMetrics["delayed"], prometheus.CounterValue, float64(zone.Delayed), name)
		ch <- prometheus.MustNewConstMetric(c.limitRequestMetrics["rejected_dry_run"], prometheus.CounterValue, float64(zone.RejectedDryRun), name)
		ch <- prometheus.MustNewConstMetric(c.limitRequestMetrics["delayed_dry_run"], prometheus.CounterValue, float64(zone.DelayedDryRun), name)
	}

	for name, zone := range stats.HTTPLimitConnections {
		ch <- prometheus.MustNewConstMetric(c.limitConnectionMetrics["passed"], prometheus.CounterValue, float64(zone.Passed), name)
		ch <- prometheus.MustNewConstMetric(c.limitConnectionMetrics["rejected"], prometheus.CounterValue, float64(zone.Rejected), name)
		ch <- prometheus.MustNewConstMetric(c.limitConnectionMetrics["rejected_dry_run"], prometheus.CounterValue, float64(zone.RejectedDryRun), name)
	}

	for name, zone := range stats.StreamLimitConnections {
		ch <- prometheus.MustNewConstMetric(c.streamLimitConnectionMetrics["passed"], prometheus.CounterValue, float64(zone.Passed), name)
		ch <- prometheus.MustNewConstMetric(c.streamLimitConnectionMetrics["rejected"], prometheus.CounterValue, float64(zone.Rejected), name)
		ch <- prometheus.MustNewConstMetric(c.streamLimitConnectionMetrics["rejected_dry_run"], prometheus.CounterValue, float64(zone.RejectedDryRun), name)
	}
}

var upstreamServerStates = map[string]float64{
	"up":        1.0,
	"draining":  2.0,
	"down":      3.0,
	"unavail":   4.0,
	"checking":  5.0,
	"unhealthy": 6.0,
}

func newServerZoneMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"server_zone"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "server_zone", metricName), docString, labels, constLabels)
}

func newStreamServerZoneMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"server_zone"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_server_zone", metricName), docString, labels, constLabels)
}

func newUpstreamMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "upstream", metricName), docString, []string{"upstream"}, constLabels)
}

func newStreamUpstreamMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_upstream", metricName), docString, []string{"upstream"}, constLabels)
}

func newUpstreamServerMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"upstream", "server"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "upstream_server", metricName), docString, labels, constLabels)
}

func newStreamUpstreamServerMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"upstream", "server"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_upstream_server", metricName), docString, labels, constLabels)
}

func newStreamZoneSyncMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_zone_sync_status", metricName), docString, nil, constLabels)
}

func newStreamZoneSyncZoneMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_zone_sync_zone", metricName), docString, []string{"zone"}, constLabels)
}

func newLocationZoneMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "location_zone", metricName), docString, []string{"location_zone"}, constLabels)
}

func newResolverMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "resolver", metricName), docString, []string{"resolver"}, constLabels)
}

func newLimitRequestMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "limit_request", metricName), docString, []string{"zone"}, constLabels)
}

func newLimitConnectionMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "limit_connection", metricName), docString, []string{"zone"}, constLabels)
}

func newStreamLimitConnectionMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_limit_connection", metricName), docString, []string{"zone"}, constLabels)
}
