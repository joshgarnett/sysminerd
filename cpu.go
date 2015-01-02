package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var previousCPUStats map[string][]float64

// CPU Fields for /proc/stat
var cpuFields = map[string]int{
	"User":      1,
	"Nice":      2,
	"System":    3,
	"Idle":      4,
	"IOWait":    5,
	"IRQ":       6,
	"SoftIRQ":   7,
	"Steal":     8,
	"Guest":     9,
	"GuestNice": 10,
}

type CPUInputModule struct{}

func (m *CPUInputModule) Name() string {
	return "cpu"
}

func (m *CPUInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *CPUInputModule) TearDown() error {
	return nil
}

func (m *CPUInputModule) GetMetrics() ([]Metric, error) {
	metrics := make([]Metric, 0, 10)

	b, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")

	cpus := make(map[string][]float64)

	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0][:3] == "cpu" {
			if i != 0 {
				cpuStats := make([]float64, 11)

				for j, field := range fields {
					if j != 0 {
						value, err := strconv.ParseFloat(field, 64)
						if err != nil {
							log.Printf("Error parsing %s as float64: %v", field, err)
						} else {
							cpuStats[0] += value
							cpuStats[j] = value
						}
					}
				}
				cpus[fields[0]] = cpuStats
			}
		}
	}

	now := time.Now()

	if previousCPUStats != nil {
		for cpu, values := range cpus {
			totalDiff := values[0] - previousCPUStats[cpu][0]
			for name, index := range cpuFields {
				value := values[index] - previousCPUStats[cpu][index]

				niceMetric := Metric{
					module:    m.Name(),
					name:      fmt.Sprintf("%s.%s", cpu, name),
					value:     value / totalDiff,
					timestamp: now,
				}
				metrics = append(metrics, niceMetric)
			}
		}
	}

	previousCPUStats = cpus

	return metrics, nil
}
