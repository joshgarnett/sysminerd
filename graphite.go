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
	prefix        string
	hostname      string
	port          int64
	protocol      string
	conn          net.Conn
	queuedMetrics []Metric
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
	graphiteConnection, err := connectToGraphite(graphiteHostname, graphitePort, protocol)

	// save config data
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
	connectionFailed := false

	allMetrics := metrics

	// add queued metrics also
	if len(m.queuedMetrics) > 0 {
		allMetrics = append(allMetrics, m.queuedMetrics...)
		m.queuedMetrics = make([]Metric, 0, 0)
	}

	// for now just print the metrics
	for _, metric := range allMetrics {
		if connectionFailed {
			m.queuedMetrics = append(m.queuedMetrics, metric)
			continue
		}

		metricName := fmt.Sprintf("%s.%s.%s", m.prefix, metric.module, metric.name)

		graphiteMetric := fmt.Sprintf("%s %f %d\n", metricName, metric.value, metric.timestamp.Unix())
		log.Printf("Graphite: %s", graphiteMetric)

		n, err := m.conn.Write([]byte(graphiteMetric))
		if err != nil || n != len(graphiteMetric) {
			log.Printf("Error sending graphite metric: %v", err)
			m.queuedMetrics = append(m.queuedMetrics, metric)

			// attempt to reconnect
			graphiteConnection, err := connectToGraphite(m.hostname, m.port, m.protocol)
			if err == nil {
				m.conn = graphiteConnection
			}

			// don't send any metrics for the rest of this tick
			connectionFailed = true
		}
	}

	return nil, nil
}

func connectToGraphite(hostname string, port int64, protocol string) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to graphite: %v", err)
	}

	return conn, err
}
