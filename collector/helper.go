package collector

import "github.com/prometheus/client_golang/prometheus"

func newGlobalMetric(namespace string, metricName string, docString string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, nil, nil)
}
