package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const NetworkModuleName = "network"

type NetworkInputModule struct {
	previousIfaces map[string]map[string]float64
}

func (m *NetworkInputModule) Name() string {
	return NetworkModuleName
}

func (m *NetworkInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *NetworkInputModule) TearDown() error {
	return nil
}

func (m *NetworkInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 48)

	ifaces, err := ParseNetworkDev("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	if m.previousIfaces != nil {
		for iface, fields := range ifaces {
			previous, ok := m.previousIfaces[iface]

			if ok {
				for name, value := range fields {
					metric := NewMetric(fmt.Sprintf("%s.%s", iface, name), value-previous[name])
					metrics = append(metrics, metric)
				}
			}
		}
	}

	m.previousIfaces = ifaces

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}

func ParseNetworkDev(path string) (map[string]map[string]float64, error) {
	ifaces := make(map[string]map[string]float64)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")

	headers := []string{"face", "rx_bytes", "rx_packets", "rx_errs", "rx_drop", "rx_fifo", "rx_frame",
		"rx_compressed", "rx_multicast", "tx_bytes", "tx_packets", "tx_errs", "tx_drop", "tx_fifo",
		"tx_colls", "tx_carrier", "tx_compressed"}

	for _, line := range lines {
		line = strings.Replace(line, "|", " ", -1)
		line = strings.Replace(line, ":", " ", -1)

		fields := strings.Fields(line)
		if len(fields) < 6 {
			//throwaway Inter-| Receive | Transmit row
			continue
		}

		if fields[0] == "face" {
			//TODO: parse header row?
		} else {
			iface := make(map[string]float64)
			name := fields[0]

			for i, field := range fields {
				if i == 0 {
					continue
				}
				header := headers[i]
				value, err := strconv.ParseFloat(field, 64)
				if err == nil {
					iface[header] = value
				}
			}

			ifaces[name] = iface
		}
	}

	return ifaces, nil
}
