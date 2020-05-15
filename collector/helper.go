package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const nginxUp = 1
const nginxDown = 0

func newGlobalMetric(namespace string, metricName string, docString string, constLabels map[string]string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, nil, constLabels)
}

func newUpMetric(namespace string, constLabels map[string]string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "up",
		Help:        "Status of the last metric scrape",
		ConstLabels: constLabels,
	})
}

// MergeLabels merges two maps of labels.
func MergeLabels(a, b map[string]string) map[string]string {
	c := make(map[string]string)

	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}

	return c
}

// MergeLabelList merges two lists of label keys. Removing blank strings.
// This helper was created to avoid the creation of "" label keys that cause label cardinality issues.
func MergeLabelList(a, b []string) []string {
	for i, lab := range a {
		if lab == "" {
			a = append(a[:i], a[i+1:]...)
		}
	}
	for i, lab := range b {
		if lab == "" {
			b = append(b[:i], b[i+1:]...)
		}
	}

	return append(a, b...)
}
