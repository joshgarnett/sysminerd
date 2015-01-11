package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const CpuModuleName = "cpu"

// CPU Fields for /proc/stat
var cpuFields = map[string]int{
	"user":       1,
	"nice":       2,
	"system":     3,
	"idle":       4,
	"io_wait":    5,
	"irq":        6,
	"soft_irq":   7,
	"steal":      8,
	"guest":      9,
	"guest_nice": 10,
}

type CPUInputModule struct {
	previousCPUStats map[string][]float64
}

func (m *CPUInputModule) Name() string {
	return CpuModuleName
}

func (m *CPUInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *CPUInputModule) TearDown() error {
	return nil
}

func (m *CPUInputModule) ParseProcStat(content string) (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 10)

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

	if m.previousCPUStats != nil {
		for cpu, values := range cpus {
			totalDiff := values[0] - m.previousCPUStats[cpu][0]
			for name, index := range cpuFields {
				value := values[index] - m.previousCPUStats[cpu][index]

				metric := NewMetric(fmt.Sprintf("%s.%s", cpu, name), (value/totalDiff)*100)

				metrics = append(metrics, metric)
			}
		}
	}

	m.previousCPUStats = cpus

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
