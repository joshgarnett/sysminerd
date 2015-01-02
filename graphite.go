package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var graphitePrefix string

type GraphiteOutputModule struct{}

func (m GraphiteOutputModule) Name() string {
	return "graphite"
}

func (m GraphiteOutputModule) Init(config Config, moduleConfig map[interface{}]interface{}) error {
	var hostname string
	var err error

	if config.hostname != "" {
		hostname = config.hostname
	} else {
		hostname, err = os.Hostname()
		if err != nil {
			addrs, err := net.InterfaceAddrs()
			if err != nil || len(addrs) == 0 {
				log.Printf("Unable to get the system hostname: %v", err)
				hostname = "unknown"
			} else {
				hostname = addrs[0].String()
			}
		}
	}

	// replace periods in the fqdn with underscores
	hostname = strings.Replace(hostname, ".", "_", -1)

	graphitePrefix = fmt.Sprintf("sysminerd.%s", hostname)

	return nil
}

func (m GraphiteOutputModule) TearDown() error {
	return nil
}

func (m GraphiteOutputModule) SendMetrics(metrics []Metric) ([]Metric, error) {

	// for now just print the metrics
	for _, metric := range metrics {
		metricName := fmt.Sprintf("%s.%s.%s", graphitePrefix, metric.module, metric.name)
		log.Printf("Graphite: %-30s = %f", metricName, metric.value)
	}

	return nil, nil
}
