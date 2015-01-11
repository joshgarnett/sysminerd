// +build linux

package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func (m *MemoryInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 50)

	meminfo, err := ParseMeminfo("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, NewMetric("total", meminfo["MemTotal"]))
	metrics = append(metrics, NewMetric("free", meminfo["MemFree"]))

	used := meminfo["MemTotal"] - (meminfo["MemFree"] + meminfo["Buffers"] + meminfo["Cached"])

	metrics = append(metrics, NewMetric("used", used))
	metrics = append(metrics, NewMetric("cached", meminfo["Cached"]))
	metrics = append(metrics, NewMetric("buffer", meminfo["Buffers"]))

	bcTotal := meminfo["Cached"] + meminfo["Buffers"]
	bcFree := meminfo["MemFree"] + bcTotal

	metrics = append(metrics, NewMetric("buffer_cache_total", bcTotal))
	metrics = append(metrics, NewMetric("buffer_cache_free", bcFree))

	metrics = append(metrics, NewMetric("swap_total", meminfo["SwapTotal"]))
	metrics = append(metrics, NewMetric("swap_free", meminfo["SwapFree"]))

	metrics = append(metrics, NewMetric("high_total", meminfo["HighTotal"]))
	metrics = append(metrics, NewMetric("high_free", meminfo["HighFree"]))

	metrics = append(metrics, NewMetric("low_total", meminfo["LowTotal"]))
	metrics = append(metrics, NewMetric("low_free", meminfo["LowFree"]))

	metrics = append(metrics, NewMetric("slab_reclaimable", meminfo["SReclaimable"]))
	metrics = append(metrics, NewMetric("slab_unreclaimable", meminfo["SUnreclaim"]))

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
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
