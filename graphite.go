package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type GraphiteOutputModule struct {
	prefix   string
	hostname string
	port     int64
	protocol string
	conn     net.Conn
}

func (m *GraphiteOutputModule) Name() string {
	return "graphite"
}

func (m *GraphiteOutputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	var hostname string
	var err error

	if config.Hostname != "" {
		hostname = config.Hostname
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
	graphitePrefix := fmt.Sprintf("sysminerd.%s", hostname)

	// parse graphite settings
	graphiteHostname := moduleConfig.Settings["hostname"]
	if graphiteHostname == "" {
		log.Fatal("hostname must be specified")
	}

	graphitePort, err := strconv.ParseInt(moduleConfig.Settings["port"], 10, 64)
	if err != nil {
		log.Fatalf("Unable to parse port: %v", err)
	} else if graphitePort < 1 || graphitePort > 65535 {
		log.Fatalf("invalid port number: %d", graphitePort)
	}

	protocol := moduleConfig.Settings["protocol"]
	if protocol != "tcp" {
		log.Fatalf("Graphite protocol %s is not supported", protocol)
	}

	// connect to graphite
	address := fmt.Sprintf("%s:%d", graphiteHostname, graphitePort)
	graphiteConnection, err := net.DialTimeout(protocol, address, 5*time.Second)
	if err != nil {
		//eventually we should just log and then retry later
		log.Fatalf("Failed to connect to graphite: %v", err)
	}

	m.prefix = graphitePrefix
	m.hostname = graphiteHostname
	m.port = graphitePort
	m.protocol = protocol
	m.conn = graphiteConnection

	return nil
}

func (m *GraphiteOutputModule) TearDown() error {
	return m.conn.Close()
}

func (m *GraphiteOutputModule) SendMetrics(metrics []Metric) ([]Metric, error) {
	// for now just print the metrics
	for _, metric := range metrics {
		metricName := fmt.Sprintf("%s.%s.%s", m.prefix, metric.module, metric.name)

		graphiteMetric := fmt.Sprintf("%s %f %d\n", metricName, metric.value, metric.timestamp.Unix())
		log.Printf("Graphite: %s", graphiteMetric)

		_, err := m.conn.Write([]byte(graphiteMetric))
		if err != nil {
			log.Printf("Error sending graphite metric: %v", err)
		}
	}

	return nil, nil
}
