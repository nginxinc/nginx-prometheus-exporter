package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	plusclient "github.com/nginxinc/nginx-plus-go-client/client"
	"github.com/nginxinc/nginx-prometheus-exporter/client"
	"github.com/nginxinc/nginx-prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func getEnvUint(key string, defaultValue uint) uint {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Fatalf("Environment variable value for %s must be an uint: %v", key, err)
	}
	return uint(i)
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment variable value for %s must be a boolean: %v", key, err)
	}
	return b
}

func getEnvPositiveDuration(key string, defaultValue time.Duration) positiveDuration {
	value, ok := os.LookupEnv(key)
	if !ok {
		return positiveDuration{defaultValue}
	}

	posDur, err := parsePositiveDuration(value)
	if err != nil {
		log.Fatalf("Environment variable value for %s must be a positive duration: %v", key, err)
	}
	return posDur
}

// positiveDuration is a wrapper of time.Duration to ensure only positive values are accepted
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
		return positiveDuration{}, err
	}
	if dur < 0 {
		return positiveDuration{}, fmt.Errorf("negative duration %v is not valid", dur)
	}
	return positiveDuration{dur}, nil
}

func createPositiveDurationFlag(name string, value positiveDuration, usage string) *positiveDuration {
	flag.Var(&value, name, usage)
	return &value
}

func createClientWithRetries(getClient func() (interface{}, error), retries uint, retryInterval time.Duration) (interface{}, error) {
	var err error
	var nginxClient interface{}

	for i := 0; i <= int(retries); i++ {
		nginxClient, err = getClient()
		if err == nil {
			return nginxClient, nil
		}
		if i < int(retries) {
			log.Printf("Could not create Nginx Client. Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}
	return nil, err
}

var (
	// Set during go build
	version   string
	gitCommit string

	// Defaults values
	defaultListenAddress      = getEnv("LISTEN_ADDRESS", ":9113")
	defaultMetricsPath        = getEnv("TELEMETRY_PATH", "/metrics")
	defaultNginxPlus          = getEnvBool("NGINX_PLUS", false)
	defaultScrapeURI          = getEnv("SCRAPE_URI", "http://127.0.0.1:8080/stub_status")
	defaultSslVerify          = getEnvBool("SSL_VERIFY", true)
	defaultTimeout            = getEnvPositiveDuration("TIMEOUT", time.Second*5)
	defaultNginxRetries       = getEnvUint("NGINX_RETRIES", 0)
	defaultNginxRetryInterval = getEnvPositiveDuration("NGINX_RETRY_INTERVAL", time.Second*5)

	// Command-line flags
	listenAddr = flag.String("web.listen-address",
		defaultListenAddress,
		"An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable.")
	metricsPath = flag.String("web.telemetry-path",
		defaultMetricsPath,
		"A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable.")
	nginxPlus = flag.Bool("nginx.plus",
		defaultNginxPlus,
		"Start the exporter for NGINX Plus. By default, the exporter is started for NGINX. The default value can be overwritten by NGINX_PLUS environment variable.")
	scrapeURI = flag.String("nginx.scrape-uri",
		defaultScrapeURI,
		`A URI for scraping NGINX or NGINX Plus metrics.
	For NGINX, the stub_status page must be available through the URI. For NGINX Plus -- the API. The default value can be overwritten by SCRAPE_URI environment variable.`)
	sslVerify = flag.Bool("nginx.ssl-verify",
		defaultSslVerify,
		"Perform SSL certificate verification. The default value can be overwritten by SSL_VERIFY environment variable.")
	nginxRetries = flag.Uint("nginx.retries",
		defaultNginxRetries,
		"A number of retries the exporter will make on start to connect to the NGINX stub_status page/NGINX Plus API before exiting with an error. The default value can be overwritten by NGINX_RETRIES environment variable.")

	// Custom command-line flags
	timeout = createPositiveDurationFlag("nginx.timeout",
		defaultTimeout,
		"A timeout for scraping metrics from NGINX or NGINX Plus. The default value can be overwritten by TIMEOUT environment variable.")

	nginxRetryInterval = createPositiveDurationFlag("nginx.retry-interval",
		defaultNginxRetryInterval,
		"An interval between retries to connect to the NGINX stub_status page/NGINX Plus API on start. The default value can be overwritten by NGINX_RETRY_INTERVAL environment variable.")
)

func main() {
	flag.Parse()

	log.Printf("Starting NGINX Prometheus Exporter Version=%v GitCommit=%v", version, gitCommit)

	registry := prometheus.NewRegistry()

	buildInfoMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nginxexporter_build_info",
			Help: "Exporter build information",
			ConstLabels: prometheus.Labels{
				"version":   version,
				"gitCommit": gitCommit,
			},
		},
	)
	buildInfoMetric.Set(1)

	registry.MustRegister(buildInfoMetric)

	httpClient := &http.Client{
		Timeout: timeout.Duration,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !*sslVerify},
		},
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		log.Printf("SIGTERM received: %v. Exiting...", <-signalChan)
		os.Exit(0)
	}()

	if *nginxPlus {
		plusClient, err := createClientWithRetries(func() (interface{}, error) {
			return plusclient.NewNginxClient(httpClient, *scrapeURI)
		}, *nginxRetries, nginxRetryInterval.Duration)
		if err != nil {
			log.Fatalf("Could not create Nginx Plus Client: %v", err)
		}
		registry.MustRegister(collector.NewNginxPlusCollector(plusClient.(*plusclient.NginxClient), "nginxplus"))
	} else {
		ossClient, err := createClientWithRetries(func() (interface{}, error) {
			return client.NewNginxClient(httpClient, *scrapeURI)
		}, *nginxRetries, nginxRetryInterval.Duration)
		if err != nil {
			log.Fatalf("Could not create Nginx Client: %v", err)
		}
		registry.MustRegister(collector.NewNginxCollector(ossClient.(*client.NginxClient), "nginx"))
	}
	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>NGINX Exporter</title></head>
			<body>
			<h1>NGINX Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Printf("Error while sending a response for the '/' path: %v", err)
		}
	})
	log.Printf("NGINX Prometheus Exporter has successfully started")
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
