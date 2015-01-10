package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const GraphiteModuleName = "graphite"

type GraphiteOutputModule struct {
	Prefix        string
	Hostname      string
	Port          int
	Protocol      string
	MaxQueueSize  int
	conn          net.Conn
	queuedMetrics []string
}

func (m *GraphiteOutputModule) Name() string {
	return GraphiteModuleName
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
	graphiteHostname, err := moduleConfig.SettingsString("hostname")
	if err != nil || graphiteHostname == "" {
		log.Fatalf("hostname must be specified: %v", err)
	}

	graphitePort, err := moduleConfig.SettingsInt("port")
	if err != nil {
		log.Fatalf("Unable to parse port: %v", err)
	} else if graphitePort < 1 || graphitePort > 65535 {
		log.Fatalf("invalid port number: %d", graphitePort)
	}

	protocol, err := moduleConfig.SettingsString("protocol")
	if protocol != "tcp" && protocol != "udp" {
		log.Fatalf("Graphite protocol %s is not supported", protocol)
	}

	maxQueueSize, err := moduleConfig.SettingsInt("max_queue_size")
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

func (m *GraphiteOutputModule) SendMetrics(moduleMetrics []*ModuleMetrics) error {
	var err error

	metrics := make([]string, 0, len(moduleMetrics)*5)

	// convert metrics to graphite metrics
	for _, module := range moduleMetrics {
		moduleName := module.Module
		for _, metric := range module.Metrics {
			metricName := fmt.Sprintf("%s.%s.%s", m.Prefix, moduleName, metric.Name)
			graphiteMetric := fmt.Sprintf("%s %f %d\n", metricName, metric.Value, metric.Timestamp.Unix())
			metrics = append(metrics, graphiteMetric)
		}
	}

	// add queued metrics also
	if len(m.queuedMetrics) > 0 {
		metrics = append(m.queuedMetrics, metrics...)
		m.queuedMetrics = make([]string, 0, len(metrics))
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

		log.Printf("Graphite: %s", metric)

		n, err := m.conn.Write([]byte(metric))
		if err != nil || n != len(metric) {
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
