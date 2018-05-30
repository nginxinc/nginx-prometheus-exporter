package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/nginxinc/nginx-prometheus-exporter/client"
	"github.com/nginxinc/nginx-prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Set during go build
	version   string
	gitCommit string

	// Command-line flags
	listenAddr  = flag.String("web.listen-address", ":9113", "An address to listen on for web interface and telemetry.")
	metricsPath = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	nginxPlus   = flag.Bool("nginx.plus", false, "Start the exporter for NGINX Plus. By default, the exporter is started for NGINX.")
	scrapeURI   = flag.String("nginx.scrape-uri", "http://127.0.0.1:8080/stub_status",
		`A URI for scraping NGINX or NGINX Plus metrics.
	For NGINX, the stub_status page must be available through the URI. For NGINX Plus -- the API.`)
)

func main() {
	flag.Parse()

	log.Printf("Starting NGINX Prometheus Exporter Version=%v GitCommit=%v", version, gitCommit)

	registry := prometheus.NewRegistry()

	if *nginxPlus {
		client, err := client.NewNginxPlusClient(&http.Client{}, *scrapeURI)
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
