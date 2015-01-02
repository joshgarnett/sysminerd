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
	Prefix        string
	Hostname      string
	Port          int
	Protocol      string
	MaxQueueSize  int
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
	if protocol != "tcp" && protocol != "udp" {
		log.Fatalf("Graphite protocol %s is not supported", protocol)
	}

	maxQueueSize, err := strconv.ParseInt(moduleConfig.Settings["max_queue_size"], 10, 64)
	if err != nil {
		maxQueueSize = 0
	}

	// save config data
	m.Prefix = graphitePrefix
	m.Hostname = graphiteHostname
	m.Port = int(graphitePort)
	m.Protocol = protocol
	m.MaxQueueSize = int(maxQueueSize)

	// connect to graphite
	m.conn, err = connectToGraphite(m.Hostname, m.Port, m.Protocol)

	return err
}

func (m *GraphiteOutputModule) TearDown() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

func (m *GraphiteOutputModule) SendMetrics(metrics []Metric) error {
	var err error

	// add queued metrics also
	if len(m.queuedMetrics) > 0 {
		metrics = append(m.queuedMetrics, metrics...)
		m.queuedMetrics = make([]Metric, 0, len(metrics))
	}

	// attempt to reconnect to graphite
	if m.conn == nil {
		graphiteConnection, err := connectToGraphite(m.Hostname, m.Port, m.Protocol)
		if err == nil {
			log.Print("Reconnected to graphite")
			m.conn = graphiteConnection
		}
	}

	// for now just print the metrics
	for _, metric := range metrics {
		if m.conn == nil {
			m.queuedMetrics = append(m.queuedMetrics, metric)
			continue
		}

		metricName := fmt.Sprintf("%s.%s.%s", m.Prefix, metric.module, metric.name)

		graphiteMetric := fmt.Sprintf("%s %f %d\n", metricName, metric.value, metric.timestamp.Unix())
		log.Printf("Graphite: %s", graphiteMetric)

		n, err := m.conn.Write([]byte(graphiteMetric))
		if err != nil || n != len(graphiteMetric) {
			log.Printf("Error sending graphite metric: %v", err)
			m.queuedMetrics = append(m.queuedMetrics, metric)

			// close the existing connection
			m.conn.Close()
			m.conn = nil
		}
	}

	//see if we need to trim the queued metrics
	if m.MaxQueueSize > 0 && len(m.queuedMetrics) > m.MaxQueueSize {
		log.Printf("Graphite metric queue overflow, throwing away %d metrics", len(m.queuedMetrics)-m.MaxQueueSize)
		m.queuedMetrics = m.queuedMetrics[len(m.queuedMetrics)-m.MaxQueueSize:]
	}

	return err
}

func connectToGraphite(hostname string, port int, protocol string) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to graphite: %v", err)
	}

	return conn, err
}
