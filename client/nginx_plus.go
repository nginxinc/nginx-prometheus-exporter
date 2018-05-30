package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client version - 380a939

// APIVersion is a version of NGINX Plus API.
const APIVersion = 2

// NginxPlusClient lets you add/remove servers to/from NGINX Plus via its API.
type NginxPlusClient struct {
	apiEndpoint string
	httpClient  *http.Client
}

type versions []int

// UpstreamServer lets you configure HTTP upstreams.
type UpstreamServer struct {
	ID          int    `json:"id,omitempty"`
	Server      string `json:"server"`
	MaxFails    int    `json:"max_fails"`
	FailTimeout string `json:"fail_timeout,omitempty"`
	SlowStart   string `json:"slow_start,omitempty"`
}

// StreamUpstreamServer lets you configure Stream upstreams.
type StreamUpstreamServer struct {
	ID          int64  `json:"id,omitempty"`
	Server      string `json:"server"`
	MaxFails    int64  `json:"max_fails"`
	FailTimeout string `json:"fail_timeout,omitempty"`
	SlowStart   string `json:"slow_start,omitempty"`
}

type apiErrorResponse struct {
	Path      string
	Method    string
	Error     apiError
	RequestID string `json:"request_id"`
	Href      string
}

func (resp *apiErrorResponse) toString() string {
	return fmt.Sprintf("path=%v; method=%v; error.status=%v; error.text=%v; error.code=%v; request_id=%v; href=%v",
		resp.Path, resp.Method, resp.Error.Status, resp.Error.Text, resp.Error.Code, resp.RequestID, resp.Href)
}

type apiError struct {
	Status int
	Text   string
	Code   string
}

// Stats represents NGINX Plus stats fetched from the NGINX Plus API.
// https://nginx.org/en/docs/http/ngx_http_api_module.html
type Stats struct {
	Connections  Connections
	HTTPRequests HTTPRequests
	SSL          SSL
	ServerZones  ServerZones
	Upstreams    Upstreams
}

// Connections represents connection related stats.
type Connections struct {
	Accepted uint64
	Dropped  uint64
	Active   uint64
	Idle     uint64
}

// HTTPRequests represents HTTP request related stats.
type HTTPRequests struct {
	Total   uint64
	Current uint64
}

// SSL represents SSL related stats.
type SSL struct {
	Handshakes       uint64
	HandshakesFailed uint64 `json:"handshakes_failed"`
	SessionReuses    uint64 `json:"session_reuses"`
}

// ServerZones is map of server zone stats by zone name
type ServerZones map[string]ServerZone

// ServerZone represents server zone related stats.
type ServerZone struct {
	Processing uint64
	Requests   uint64
	Responses  Responses
	Discarded  uint64
	Received   uint64
	Sent       uint64
}

// Responses represents HTTP reponse related stats.
type Responses struct {
	Responses1xx uint64 `json:"1xx"`
	Responses2xx uint64 `json:"2xx"`
	Responses3xx uint64 `json:"3xx"`
	Responses4xx uint64 `json:"4xx"`
	Responses5xx uint64 `json:"5xx"`
}

// Upstreams is a map of upstream stats by upstream name.
type Upstreams map[string]Upstream

// Upstream represents upstream related stats.
type Upstream struct {
	Peers      []Peer
	Keepalives int
	Zombies    int
	Zone       string
	Queue      Queue
}

// Queue represents queue related stats for an upstream.
type Queue struct {
	Size      int
	MaxSize   int `json:"max_size"`
	Overflows uint64
}

// Peer represents peer (upstream server) related stats.
type Peer struct {
	ID           int
	Server       string
	Service      string
	Name         string
	Backup       bool
	Weight       int
	State        string
	Active       uint64
	MaxConns     int `json:"max_conns"`
	Requests     uint64
	Responses    Responses
	Sent         uint64
	Received     uint64
	Fails        uint64
	Unavail      uint64
	HealthChecks HealthChecks `json:"health_checks"`
	Downtime     uint64
	Downstart    string
	Selected     string
	HeaderTime   uint64 `json:"header_time"`
	ResponseTime uint64 `json:"response_time"`
}

// HealthChecks represents health check related stats for a peer.
type HealthChecks struct {
	Checks     uint64
	Fails      uint64
	Unhealthy  uint64
	LastPassed bool `json:"last_passed"`
}

// NewNginxPlusClient creates an NginxPlusClient.
func NewNginxPlusClient(httpClient *http.Client, apiEndpoint string) (*NginxPlusClient, error) {
	versions, err := getAPIVersions(httpClient, apiEndpoint)

	if err != nil {
		return nil, fmt.Errorf("error accessing the API: %v", err)
	}

	found := false
	for _, v := range *versions {
		if v == APIVersion {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("API version %v of the client is not supported by API versions of NGINX Plus: %v", APIVersion, *versions)
	}

	return &NginxPlusClient{
		apiEndpoint: apiEndpoint,
		httpClient:  httpClient,
	}, nil
}

func getAPIVersions(httpClient *http.Client, endpoint string) (*versions, error) {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("%v is not accessible: %v", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v is not accessible: expected %v response, got %v", endpoint, http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body of the response: %v", err)
	}

	var vers versions
	err = json.Unmarshal(body, &vers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling versions, got %q response: %v", string(body), err)
	}

	return &vers, nil
}

func createResponseMismatchError(respBody io.ReadCloser, mainErr error) error {
	apiErr, err := readAPIErrorResponse(respBody)
	if err != nil {
		return fmt.Errorf("%v; failed to read the response body: %v", mainErr, err)
	}

	return fmt.Errorf("%v; error: %v", mainErr, apiErr.toString())
}

func readAPIErrorResponse(respBody io.ReadCloser) (*apiErrorResponse, error) {
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %v", err)
	}

	var apiErr apiErrorResponse
	err = json.Unmarshal(body, &apiErr)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling apiErrorResponse: got %q response: %v", string(body), err)
	}

	return &apiErr, nil
}

// CheckIfUpstreamExists checks if the upstream exists in NGINX. If the upstream doesn't exist, it returns the error.
func (client *NginxPlusClient) CheckIfUpstreamExists(upstream string) error {
	_, err := client.GetHTTPServers(upstream)
	return err
}

// GetHTTPServers returns the servers of the upstream from NGINX.
func (client *NginxPlusClient) GetHTTPServers(upstream string) ([]UpstreamServer, error) {
	path := fmt.Sprintf("http/upstreams/%v/servers", upstream)

	var servers []UpstreamServer
	err := client.get(path, &servers)

	if err != nil {
		return nil, fmt.Errorf("failed to get the HTTP servers of upstream %v: %v", upstream, err)
	}

	return servers, nil
}

// AddHTTPServer adds the server to the upstream.
func (client *NginxPlusClient) AddHTTPServer(upstream string, server UpstreamServer) error {
	id, err := client.getIDOfHTTPServer(upstream, server.Server)

	if err != nil {
		return fmt.Errorf("failed to add %v server to %v upstream: %v", server.Server, upstream, err)
	}
	if id != -1 {
		return fmt.Errorf("failed to add %v server to %v upstream: server already exists", server.Server, upstream)
	}

	path := fmt.Sprintf("http/upstreams/%v/servers/", upstream)
	err = client.post(path, &server)
	if err != nil {
		return fmt.Errorf("failed to add %v server to %v upstream: %v", server.Server, upstream, err)
	}

	return nil
}

// DeleteHTTPServer the server from the upstream.
func (client *NginxPlusClient) DeleteHTTPServer(upstream string, server string) error {
	id, err := client.getIDOfHTTPServer(upstream, server)
	if err != nil {
		return fmt.Errorf("failed to remove %v server from  %v upstream: %v", server, upstream, err)
	}
	if id == -1 {
		return fmt.Errorf("failed to remove %v server from %v upstream: server doesn't exist", server, upstream)
	}

	path := fmt.Sprintf("http/upstreams/%v/servers/%v", upstream, id)
	err = client.delete(path)

	if err != nil {
		return fmt.Errorf("failed to remove %v server from %v upstream: %v", server, upstream, err)
	}

	return nil
}

// UpdateHTTPServers updates the servers of the upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
func (client *NginxPlusClient) UpdateHTTPServers(upstream string, servers []UpstreamServer) ([]UpstreamServer, []UpstreamServer, error) {
	serversInNginx, err := client.GetHTTPServers(upstream)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update servers of %v upstream: %v", upstream, err)
	}

	toAdd, toDelete := determineUpdates(servers, serversInNginx)

	for _, server := range toAdd {
		err := client.AddHTTPServer(upstream, server)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update servers of %v upstream: %v", upstream, err)
		}
	}

	for _, server := range toDelete {
		err := client.DeleteHTTPServer(upstream, server.Server)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update servers of %v upstream: %v", upstream, err)
		}
	}

	return toAdd, toDelete, nil
}

func determineUpdates(updatedServers []UpstreamServer, nginxServers []UpstreamServer) (toAdd []UpstreamServer, toRemove []UpstreamServer) {
	for _, server := range updatedServers {
		found := false
		for _, serverNGX := range nginxServers {
			if server.Server == serverNGX.Server {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, server)
		}
	}

	for _, serverNGX := range nginxServers {
		found := false
		for _, server := range updatedServers {
			if serverNGX.Server == server.Server {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, serverNGX)
		}
	}

	return
}

func (client *NginxPlusClient) getIDOfHTTPServer(upstream string, name string) (int, error) {
	servers, err := client.GetHTTPServers(upstream)
	if err != nil {
		return -1, fmt.Errorf("error getting id of server %v of upstream %v: %v", name, upstream, err)
	}

	for _, s := range servers {
		if s.Server == name {
			return s.ID, nil
		}
	}

	return -1, nil
}

func (client *NginxPlusClient) get(path string, data interface{}) error {
	url := fmt.Sprintf("%v/%v/%v", client.apiEndpoint, APIVersion, path)
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get %v: %v", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		mainErr := fmt.Errorf("expected %v response, got %v", http.StatusOK, resp.StatusCode)
		return createResponseMismatchError(resp.Body, mainErr)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %v", err)
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return fmt.Errorf("error unmarshaling response %q: %v", string(body), err)
	}
	return nil
}

func (client *NginxPlusClient) post(path string, input interface{}) error {
	url := fmt.Sprintf("%v/%v/%v", client.apiEndpoint, APIVersion, path)

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshall input: %v", err)
	}

	resp, err := client.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonInput))
	if err != nil {
		return fmt.Errorf("failed to post %v: %v", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		mainErr := fmt.Errorf("expected %v response, got %v", http.StatusCreated, resp.StatusCode)
		return createResponseMismatchError(resp.Body, mainErr)
	}

	return nil
}

func (client *NginxPlusClient) delete(path string) error {
	path = fmt.Sprintf("%v/%v/%v/", client.apiEndpoint, APIVersion, path)

	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to create a delete request: %v", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		mainErr := fmt.Errorf("failed to complete delete request: expected %v response, got %v",
			http.StatusOK, resp.StatusCode)
		return createResponseMismatchError(resp.Body, mainErr)
	}
	return nil
}

// CheckIfStreamUpstreamExists checks if the stream upstream exists in NGINX. If the upstream doesn't exist, it returns the error.
func (client *NginxPlusClient) CheckIfStreamUpstreamExists(upstream string) error {
	_, err := client.GetStreamServers(upstream)
	return err
}

// GetStreamServers returns the stream servers of the upstream from NGINX.
func (client *NginxPlusClient) GetStreamServers(upstream string) ([]StreamUpstreamServer, error) {
	path := fmt.Sprintf("stream/upstreams/%v/servers", upstream)

	var servers []StreamUpstreamServer
	err := client.get(path, &servers)

	if err != nil {
		return nil, fmt.Errorf("failed to get stream servers of upstream server %v: %v", upstream, err)
	}

	return servers, nil
}

// AddStreamServer adds the server to the upstream.
func (client *NginxPlusClient) AddStreamServer(upstream string, server StreamUpstreamServer) error {
	id, err := client.getIDOfStreamServer(upstream, server.Server)

	if err != nil {
		return fmt.Errorf("failed to add %v stream server to %v upstream: %v", server.Server, upstream, err)
	}
	if id != -1 {
		return fmt.Errorf("failed to add %v stream server to %v upstream: server already exists", server.Server, upstream)
	}

	path := fmt.Sprintf("stream/upstreams/%v/servers/", upstream)
	err = client.post(path, &server)

	if err != nil {
		return fmt.Errorf("failed to add %v stream server to %v upstream: %v", server.Server, upstream, err)
	}

	return nil
}

// DeleteStreamServer the server from the upstream.
func (client *NginxPlusClient) DeleteStreamServer(upstream string, server string) error {
	id, err := client.getIDOfStreamServer(upstream, server)
	if err != nil {
		return fmt.Errorf("failed to remove %v stream server from  %v upstream: %v", server, upstream, err)
	}
	if id == -1 {
		return fmt.Errorf("failed to remove %v stream server from %v upstream: server doesn't exist", server, upstream)
	}

	path := fmt.Sprintf("stream/upstreams/%v/servers/%v", upstream, id)
	err = client.delete(path)

	if err != nil {
		return fmt.Errorf("failed to remove %v stream server from %v upstream: %v", server, upstream, err)
	}

	return nil
}

// UpdateStreamServers updates the servers of the upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
func (client *NginxPlusClient) UpdateStreamServers(upstream string, servers []StreamUpstreamServer) ([]StreamUpstreamServer, []StreamUpstreamServer, error) {
	serversInNginx, err := client.GetStreamServers(upstream)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update stream servers of %v upstream: %v", upstream, err)
	}

	toAdd, toDelete := determineStreamUpdates(servers, serversInNginx)

	for _, server := range toAdd {
		err := client.AddStreamServer(upstream, server)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update stream servers of %v upstream: %v", upstream, err)
		}
	}

	for _, server := range toDelete {
		err := client.DeleteStreamServer(upstream, server.Server)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update stream servers of %v upstream: %v", upstream, err)
		}
	}

	return toAdd, toDelete, nil
}

func (client *NginxPlusClient) getIDOfStreamServer(upstream string, name string) (int64, error) {
	servers, err := client.GetStreamServers(upstream)
	if err != nil {
		return -1, fmt.Errorf("error getting id of stream server %v of upstream %v: %v", name, upstream, err)
	}

	for _, s := range servers {
		if s.Server == name {
			return s.ID, nil
		}
	}

	return -1, nil
}

func determineStreamUpdates(updatedServers []StreamUpstreamServer, nginxServers []StreamUpstreamServer) (toAdd []StreamUpstreamServer, toRemove []StreamUpstreamServer) {
	for _, server := range updatedServers {
		found := false
		for _, serverNGX := range nginxServers {
			if server.Server == serverNGX.Server {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, server)
		}
	}

	for _, serverNGX := range nginxServers {
		found := false
		for _, server := range updatedServers {
			if serverNGX.Server == server.Server {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, serverNGX)
		}
	}

	return
}

// GetStats gets connection, request, ssl, zone, and upstream related stats from the NGINX Plus API.
func (client *NginxPlusClient) GetStats() (*Stats, error) {
	cons, err := client.getConnections()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}

	requests, err := client.getHTTPRequests()
	if err != nil {
		return nil, fmt.Errorf("Failed to get stats: %v", err)
	}

	ssl, err := client.getSSL()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}

	zones, err := client.getServerZones()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}

	upstreams, err := client.getUpstreams()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}

	return &Stats{
		Connections:  *cons,
		HTTPRequests: *requests,
		SSL:          *ssl,
		ServerZones:  *zones,
		Upstreams:    *upstreams,
	}, nil
}

func (client *NginxPlusClient) getConnections() (*Connections, error) {
	var cons Connections
	err := client.get("connections", &cons)
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %v", err)
	}
	return &cons, nil
}

func (client *NginxPlusClient) getHTTPRequests() (*HTTPRequests, error) {
	var requests HTTPRequests

	err := client.get("http/requests", &requests)
	if err != nil {
		return nil, fmt.Errorf("failed to get http requests: %v", err)
	}

	return &requests, nil
}

func (client *NginxPlusClient) getSSL() (*SSL, error) {
	var ssl SSL
	err := client.get("ssl", &ssl)
	if err != nil {
		return nil, fmt.Errorf("failed to get ssl: %v", err)
	}
	return &ssl, nil
}

func (client *NginxPlusClient) getServerZones() (*ServerZones, error) {
	var zones ServerZones
	err := client.get("http/server_zones", &zones)
	if err != nil {
		return nil, fmt.Errorf("failed to get server zones: %v", err)
	}
	return &zones, err
}

func (client *NginxPlusClient) getUpstreams() (*Upstreams, error) {
	var upstreams Upstreams
	err := client.get("http/upstreams", &upstreams)
	if err != nil {
		return nil, fmt.Errorf("failed to get upstreams: %v", err)
	}
	return &upstreams, nil
}
