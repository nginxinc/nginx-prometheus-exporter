package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	plusclient "github.com/nginx/nginx-plus-go-client/v2/client"
	"github.com/nginx/nginx-prometheus-exporter/client"
	"github.com/nginx/nginx-prometheus-exporter/collector"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	common_version "github.com/prometheus/common/version"

	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"

	proxyproto "github.com/pires/go-proxyproto"
)

// positiveDuration is a wrapper of time.Duration to ensure only positive values are accepted.
type positiveDuration struct{ time.Duration }

func (pd *positiveDuration) Set(s string) error {
	dur, err := parsePositiveDuration(s)
	if err != nil {
		return err
	}

	pd.Duration = dur.Duration
	return nil
}

func parsePositiveDuration(s string) (positiveDuration, error) {
	dur, err := time.ParseDuration(s)
	if err != nil {
		return positiveDuration{}, fmt.Errorf("failed to parse duration %q: %w", s, err)
	}
	if dur < 0 {
		return positiveDuration{}, fmt.Errorf("negative duration %v is not valid", dur)
	}
	return positiveDuration{dur}, nil
}

func createPositiveDurationFlag(s kingpin.Settings) (target *time.Duration) {
	target = new(time.Duration)
	s.SetValue(&positiveDuration{Duration: *target})
	return
}

func parseUnixSocketAddress(address string) (string, string, error) {
	addressParts := strings.Split(address, ":")
	addressPartsLength := len(addressParts)

	if addressPartsLength > 3 || addressPartsLength < 1 {
		return "", "", errors.New("address for unix domain socket has wrong format")
	}

	unixSocketPath := addressParts[1]
	requestPath := ""
	if addressPartsLength == 3 {
		requestPath = addressParts[2]
	}
	return unixSocketPath, requestPath, nil
}

var (
	constLabels = map[string]string{}

	// Command-line flags.
	webConfig     = kingpinflag.AddFlags(kingpin.CommandLine, ":9113")
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").Envar("TELEMETRY_PATH").String()
	nginxPlus     = kingpin.Flag("nginx.plus", "Start the exporter for NGINX Plus. By default, the exporter is started for NGINX.").Default("false").Envar("NGINX_PLUS").Bool()
	scrapeURIs    = kingpin.Flag("nginx.scrape-uri", "A URI or unix domain socket path for scraping NGINX or NGINX Plus metrics. For NGINX, the stub_status page must be available through the URI. For NGINX Plus -- the API. Repeatable for multiple URIs.").Default("http://127.0.0.1:8080/stub_status").Envar("SCRAPE_URI").HintOptions("http://127.0.0.1:8080/stub_status", "http://127.0.0.1:8080/api").Strings()
	sslVerify     = kingpin.Flag("nginx.ssl-verify", "Perform SSL certificate verification.").Default("false").Envar("SSL_VERIFY").Bool()
	sslCaCert     = kingpin.Flag("nginx.ssl-ca-cert", "Path to the PEM encoded CA certificate file used to validate the servers SSL certificate.").Default("").Envar("SSL_CA_CERT").String()
	sslClientCert = kingpin.Flag("nginx.ssl-client-cert", "Path to the PEM encoded client certificate file to use when connecting to the server.").Default("").Envar("SSL_CLIENT_CERT").String()
	sslClientKey  = kingpin.Flag("nginx.ssl-client-key", "Path to the PEM encoded client certificate key file to use when connecting to the server.").Default("").Envar("SSL_CLIENT_KEY").String()
	useProxyProto = kingpin.Flag("nginx.proxy-protocol", "Pass proxy protocol payload to nginx listeners.").Default("false").Envar("PROXY_PROTOCOL").Bool()

	// Custom command-line flags.
	timeout = createPositiveDurationFlag(kingpin.Flag("nginx.timeout", "A timeout for scraping metrics from NGINX or NGINX Plus.").Default("5s").Envar("TIMEOUT").HintOptions("5s", "10s", "30s", "1m", "5m"))
)

const exporterName = "nginx_exporter"

func main() {
	kingpin.Flag("prometheus.const-label", "Label that will be used in every metric. Format is label=value. It can be repeated multiple times.").Envar("CONST_LABELS").StringMapVar(&constLabels)

	// convert deprecated flags to new format
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			newArg := "-" + arg
			fmt.Printf("the flag format is deprecated and will be removed in a future release, please use the new format: %s\n", newArg)
			os.Args[i] = newArg
		}
	}

	config := &promslog.Config{}

	flag.AddFlags(kingpin.CommandLine, config)
	kingpin.Version(common_version.Print(exporterName))
	kingpin.HelpFlag.Short('h')

	addMissingEnvironmentFlags(kingpin.CommandLine)

	kingpin.Parse()
	logger := promslog.New(config)

	logger.Info("nginx-prometheus-exporter", "version", common_version.Info())
	logger.Info("build context", "build_context", common_version.BuildContext())

	prometheus.MustRegister(version.NewCollector(exporterName))

	if len(*scrapeURIs) == 0 {
		logger.Error("no scrape addresses provided")
		os.Exit(1)
	}

	// #nosec G402
	sslConfig := &tls.Config{InsecureSkipVerify: !*sslVerify}
	if *sslCaCert != "" {
		caCert, err := os.ReadFile(*sslCaCert)
		if err != nil {
			logger.Error("loading CA cert failed", "err", err.Error())
			os.Exit(1)
		}
		sslCaCertPool := x509.NewCertPool()
		ok := sslCaCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			logger.Error("parsing CA cert file failed.")
			os.Exit(1)
		}
		sslConfig.RootCAs = sslCaCertPool
	}

	if *sslClientCert != "" && *sslClientKey != "" {
		clientCert, err := tls.LoadX509KeyPair(*sslClientCert, *sslClientKey)
		if err != nil {
			logger.Error("loading client certificate failed", "error", err.Error())
			os.Exit(1)
		}
		sslConfig.Certificates = []tls.Certificate{clientCert}
	}

	transport := &http.Transport{
		TLSClientConfig: sslConfig,
	}

	if len(*scrapeURIs) == 1 {
		registerCollector(logger, transport, (*scrapeURIs)[0], constLabels)
	} else {
		for _, addr := range *scrapeURIs {
			// add scrape URI to const labels
			labels := maps.Clone(constLabels)
			labels["addr"] = addr

			registerCollector(logger, transport, addr, labels)
		}
	}

	http.Handle(*metricsPath, promhttp.Handler())

	if *metricsPath != "/" && *metricsPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "NGINX Prometheus Exporter",
			Description: "Prometheus Exporter for NGINX and NGINX Plus",
			HeaderColor: "#039900",
			Version:     common_version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error("failed to create landing page", "error", err.Error())
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	srv := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("HTTP server closed", "error", err.Error())
				os.Exit(0)
			}
			logger.Error("HTTP server failed", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	srvCtx, srvCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer srvCancel()
	_ = srv.Shutdown(srvCtx)
}

func registerCollector(logger *slog.Logger, transport *http.Transport,
	addr string, labels map[string]string,
) {
	var socketPath string

	if strings.HasPrefix(addr, "unix:") {
		var err error
		var requestPath string
		socketPath, requestPath, err = parseUnixSocketAddress(addr)
		if err != nil {
			logger.Error("parsing unix domain socket scrape address failed", "uri", addr, "error", err.Error())
			os.Exit(1)
		}
		addr = "http://unix" + requestPath
	}

	if !*useProxyProto && socketPath != "" {
		transport.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		}
	}

	if *useProxyProto {
		transport.DialContext = func(_ context.Context, network, addr string) (net.Conn, error) {
			if socketPath != "" {
				network = "unix"
				addr = socketPath
			}

			conn, err := (&net.Dialer{}).Dial(network, addr)
			if err != nil {
				return nil, fmt.Errorf("dialing %s %s: %w", network, addr, err)
			}

			localAddr := conn.LocalAddr()
			remoteAddr := conn.RemoteAddr()
			transportProtocol := proxyproto.TCPv4

			switch addr := remoteAddr.(type) {
			case *net.TCPAddr:
				if addr.IP.To4() == nil {
					transportProtocol = proxyproto.TCPv6
				}
			case *net.UnixAddr:
				transportProtocol = proxyproto.UnixStream
			}

			header := &proxyproto.Header{
				Version:           2,
				Command:           proxyproto.PROXY,
				TransportProtocol: transportProtocol,
				SourceAddr:        localAddr,
				DestinationAddr:   remoteAddr,
			}

			_, err = header.WriteTo(conn)
			if err != nil {
				return nil, fmt.Errorf("writing proxyproto header: %w", err)
			}

			return conn, nil
		}
	}

	userAgent := fmt.Sprintf("NGINX-Prometheus-Exporter/v%v", common_version.Version)

	httpClient := &http.Client{
		Timeout: *timeout,
		Transport: &userAgentRoundTripper{
			agent: userAgent,
			rt:    transport,
		},
	}

	if *nginxPlus {
		plusClient, err := plusclient.NewNginxClient(addr, plusclient.WithHTTPClient(httpClient))
		if err != nil {
			logger.Error("could not create Nginx Plus Client", "error", err.Error())
			os.Exit(1)
		}
		variableLabelNames := collector.NewVariableLabelNames(nil, nil, nil, nil, nil, nil, nil)
		prometheus.MustRegister(collector.NewNginxPlusCollector(plusClient, "nginxplus", variableLabelNames, labels, logger))
	} else {
		ossClient := client.NewNginxClient(httpClient, addr)
		prometheus.MustRegister(collector.NewNginxCollector(ossClient, "nginx", labels, logger))
	}
}

type userAgentRoundTripper struct {
	rt    http.RoundTripper
	agent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.Header.Set("User-Agent", rt.agent)
	roundTrip, err := rt.rt.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("round trip failed: %w", err)
	}
	return roundTrip, nil
}

func cloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)
	*r = *req // shallow clone

	// deep copy headers
	r.Header = make(http.Header, len(req.Header))
	for key, values := range req.Header {
		newValues := make([]string, len(values))
		copy(newValues, values)
		r.Header[key] = newValues
	}
	return r
}

// addMissingEnvironmentFlags sets Envar on any flag which has
// the "web." prefix which doesn't already have an Envar set.
func addMissingEnvironmentFlags(ka *kingpin.Application) {
	for _, f := range ka.Model().FlagGroupModel.Flags {
		if strings.HasPrefix(f.Name, "web.") && f.Envar == "" {
			retrievedFlag := ka.GetFlag(f.Name)
			if retrievedFlag != nil {
				retrievedFlag.Envar(convertFlagToEnvar(strings.TrimPrefix(f.Name, "web.")))
			}
		}
	}
}

func convertFlagToEnvar(f string) string {
	env := strings.ToUpper(f)
	for _, s := range []string{"-", "."} {
		env = strings.ReplaceAll(env, s, "_")
	}
	return env
}
