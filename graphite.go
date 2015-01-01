package main

import (
	"fmt"
	"log"
)

type GraphiteOutputModule struct{}

func (m GraphiteOutputModule) Name() string {
	return "graphite"
}

func (m GraphiteOutputModule) Init(config map[interface{}]interface{}) error {
	return nil
}

func (m GraphiteOutputModule) TearDown() error {
	return nil
}

func (m GraphiteOutputModule) SendMetrics(metrics []Metric) ([]Metric, error) {
	// for now just print the metrics
	for _, metric := range metrics {
		metricName := fmt.Sprintf("%s.%s", metric.module, metric.name)
		log.Printf("Graphite: %-20s = %f", metricName, metric.value)
	}

	return nil, nil
}
