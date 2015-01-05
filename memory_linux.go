// +build linux

package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func (m *MemoryInputModule) GetMetrics() ([]Metric, error) {
	metrics := make([]Metric, 0, 50)
	now := time.Now()

	meminfo, err := ParseMeminfo("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, Metric{module: m.Name(), name: "Total", value: meminfo["MemTotal"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Free", value: meminfo["MemFree"], timestamp: now})

	used := meminfo["MemTotal"] - (meminfo["MemFree"])

	metrics = append(metrics, Metric{module: m.Name(), name: "Used", value: used, timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Cached", value: meminfo["Cached"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Buffer", value: meminfo["Buffers"], timestamp: now})

	bcTotal := meminfo["Cached"] + meminfo["Buffers"]
	bcUsed := used - bcTotal
	bcFree := meminfo["MemFree"] + bcTotal

	metrics = append(metrics, Metric{module: m.Name(), name: "BufferCacheTotal", value: bcTotal, timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "BufferCacheUser", value: bcUsed, timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "BufferCacheFree", value: bcFree, timestamp: now})

	metrics = append(metrics, Metric{module: m.Name(), name: "SwapTotal", value: meminfo["SwapTotal"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "SwapFree", value: meminfo["SwapFree"], timestamp: now})

	metrics = append(metrics, Metric{module: m.Name(), name: "HighTotal", value: meminfo["HighTotal"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "HighFree", value: meminfo["HighFree"], timestamp: now})

	metrics = append(metrics, Metric{module: m.Name(), name: "LowTotal", value: meminfo["LowTotal"], timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "LowFree", value: meminfo["LowFree"], timestamp: now})

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
