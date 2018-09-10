package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

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

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment Variable value for %s must be a boolean", key)
	}
	return b
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
)

func main() {
	flag.Parse()

	log.Printf("Starting NGINX Prometheus Exporter Version=%v GitCommit=%v", version, gitCommit)

	registry := prometheus.NewRegistry()

	if *nginxPlus {
		client, err := plusclient.NewNginxClient(&http.Client{}, *scrapeURI)
		if err != nil {
			log.Fatalf("Could not create Nginx Plus Client: %v", err)
		}

		registry.MustRegister(collector.NewNginxPlusCollector(client, "nginxplus"))
	} else {
		client, err := client.NewNginxClient(&http.Client{}, *scrapeURI)
		if err != nil {
			log.Fatalf("Could not create Nginx Client: %v", err)
		}

		registry.MustRegister(collector.NewNginxCollector(client, "nginx"))
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
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
