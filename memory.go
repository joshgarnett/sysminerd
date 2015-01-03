package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const MemoryModuleName = "memory"

type MemoryInputModule struct{}

func (m *MemoryInputModule) Name() string {
	return MemoryModuleName
}

func (m *MemoryInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *MemoryInputModule) TearDown() error {
	return nil
}

func (m *MemoryInputModule) GetMetrics() ([]Metric, error) {
	metrics := make([]Metric, 0, 50)
	now := time.Now()

	meminfo, err := ParseMeminfo("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	used := meminfo["MemTotal"] - (meminfo["MemFree"] + meminfo["Buffers"] + meminfo["Cached"] + meminfo["Slab"])

	metrics = append(metrics, Metric{module: m.Name(), name: "Used", value: used, timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Free", value: meminfo["MemFree"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Buffered", value: meminfo["Buffers"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Cached", value: meminfo["Cached"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "SlabReclaimable", value: meminfo["SReclaimable"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "SlabUnreclaimable", value: meminfo["SUnreclaim"], timestamp: now})

	return metrics, nil
}

func ParseMeminfo(path string) (map[string]float64, error) {
	meminfo := make(map[string]float64)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		name := strings.Split(fields[0], ":")[0]
		value, err := strconv.ParseFloat(fields[1], 64)
		if err == nil {
			if len(fields) > 2 && fields[2] == "kB" {
				value = value * 1024
			}

			meminfo[name] = value
		}
	}

	return meminfo, nil
}
