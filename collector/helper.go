package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const nginxUp = 1
const nginxDown = 0

func newGlobalMetric(namespace string, metricName string, docString string, constLabels map[string]string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, nil, constLabels)
}

func newUpMetric(namespace string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "up",
		Help:      "Status of the last metric scrape",
	})
}

// MergeLabels merges two maps of labels.
func MergeLabels(a map[string]string, b map[string]string) map[string]string {
	c := make(map[string]string)

	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}

	return c
}

// reservedLabelPrefix is a prefix which is not legal in user-supplied label names.
const reservedLabelPrefix = "__"

// IsValidLabelName does equivalent validation to checkLabelName in prometheus/client_golang
func IsValidLabelName(ln string) bool {
	if len(ln) == 0 || strings.HasPrefix(ln, reservedLabelPrefix) {
		return false
	}
	for i, b := range ln {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || (b >= '0' && b <= '9' && i > 0)) {
			return false
		}
	}
	return true
}
