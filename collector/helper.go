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

// MergeLabelValues merges two maps of labels.
func MergeLabelValues(labels ...map[string]string) map[string]string {
	a := make(map[string]string)
	for _, b := range labels {
		if b != nil {
			for k, v := range b {
				a[k] = v
			}
		}
	}
	return a
}

// RemoveBlankLabelNames removes any labels of "". Blank labels cause cardinality issues.
func RemoveBlankLabelNames(labelNames ...string) []string {
	for i, lab := range labelNames {
		if lab == "" {
			labelNames = append(labelNames[:i], labelNames[i+1:]...)
		}
	}
	return labelNames
}
