package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	plusclient "github.com/nginxinc/nginx-plus-go-sdk/client"
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

func getEnvInt(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("Environment variable value for %s must be an int", key)
	}
	return int(b)
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment variable value for %s must be a boolean", key)
	}
	return b
}

func createClientWithRetries(httpClient http.Client, scrapeURI, clientName string, nginxPlus bool, retries int, retryInterval time.Duration) (interface{}, error) {
	var err error
	var nginxClient interface{}

	for {
		if nginxPlus {
			nginxClient, err = plusclient.NewNginxClient(&httpClient, scrapeURI)
		} else {
			nginxClient, err = client.NewNginxClient(&httpClient, scrapeURI)
		}
		if err != nil {
			if retries > 0 {
				retries--
				time.Sleep(retryInterval)
				log.Printf("Could not create %v. Retrying in %v...", clientName, *nginxRetryInterval)
				continue
			}
			return nil, err
		}
		return nginxClient, nil
	}
}

var (
	// Set during go build
	version   string
	gitCommit string

	// Defaults values
	defaultListenAddress = getEnv("LISTEN_ADDRESS", ":9113")
	defaultMetricsPath   = getEnv("TELEMETRY_PATH", "/metrics")
	defaultNginxPlus     = getEnvBool("NGINX_PLUS", false)
	defaultScrapeURI     = getEnv("SCRAPE_URI", "http://127.0.0.1:8080/stub_status")
	defaultSslVerify     = getEnvBool("SSL_VERIFY", true)
	defaultNginxRetries  = getEnvInt("NGINX_RETRIES", 0)

	// Command-line flags
	listenAddr = flag.String("web.listen-address", defaultListenAddress,
		"An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable.")
	metricsPath = flag.String("web.telemetry-path", defaultMetricsPath,
		"A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable.")
	nginxPlus = flag.Bool("nginx.plus", defaultNginxPlus,
		"Start the exporter for NGINX Plus. By default, the exporter is started for NGINX. The default value can be overwritten by NGINX_PLUS environment variable.")
	scrapeURI = flag.String("nginx.scrape-uri", defaultScrapeURI,
		`A URI for scraping NGINX or NGINX Plus metrics.
	For NGINX, the stub_status page must be available through the URI. For NGINX Plus -- the API. The default value can be overwritten by SCRAPE_URI environment variable.`)
	sslVerify = flag.Bool("nginx.ssl-verify", defaultSslVerify,
		"Perform SSL certificate verification. The default value can be overwritten by SSL_VERIFY environment variable.")
	timeout      = flag.Duration("nginx.timeout", 5*time.Second, "A timeout for scraping metrics from NGINX or NGINX Plus.")
	nginxRetries = flag.Int("nginx.retries", defaultNginxRetries,
		"A number of retries the exporter will make on start to connect to the NGINX stub_status page/NGINX Plus API before exiting with an error.")
	nginxRetryInterval = flag.Duration("nginx.retry-interval", 5*time.Second, "An interval between retries to connect to the NGINX stub_status page/NGINX Plus API on start.")
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
		Timeout: *timeout,
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
		plusClient, err := createClientWithRetries(*httpClient, *scrapeURI, "Nginx Plus Client", *nginxPlus, *nginxRetries, *nginxRetryInterval)
		if err != nil {
			log.Fatalf("Could not create Nginx Plus Client: %v", err)
		}
		registry.MustRegister(collector.NewNginxPlusCollector(plusClient.(*plusclient.NginxClient), "nginxplus"))
	} else {
		ossClient, err := createClientWithRetries(*httpClient, *scrapeURI, "Nginx Client", *nginxPlus, *nginxRetries, *nginxRetryInterval)
		if err != nil {
			log.Fatalf("Could not create Nginx Client: %v", err)
		}
		registry.MustRegister(collector.NewNginxCollector(ossClient.(*client.NginxClient), "nginx"))
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NGINX Exporter</title></head>
			<body>
			<h1>NGINX Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Printf("NGINX Prometheus Exporter has successfully started")
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
