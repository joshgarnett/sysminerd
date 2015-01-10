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

	metrics = append(metrics, NewMetric("Total", meminfo["MemTotal"]))
	metrics = append(metrics, NewMetric("Free", meminfo["MemFree"]))

	used := meminfo["MemTotal"] - (meminfo["MemFree"] + meminfo["Buffers"] + meminfo["Cached"])

	metrics = append(metrics, NewMetric("Used", used))
	metrics = append(metrics, NewMetric("Cached", meminfo["Cached"]))
	metrics = append(metrics, NewMetric("Buffer", meminfo["Buffers"]))

	bcTotal := meminfo["Cached"] + meminfo["Buffers"]
	bcFree := meminfo["MemFree"] + bcTotal

	metrics = append(metrics, NewMetric("BufferCacheTotal", bcTotal))
	metrics = append(metrics, NewMetric("BufferCacheFree", bcFree))

	metrics = append(metrics, NewMetric("SwapTotal", meminfo["SwapTotal"]))
	metrics = append(metrics, NewMetric("SwapFree", meminfo["SwapFree"]))

	metrics = append(metrics, NewMetric("HighTotal", meminfo["HighTotal"]))
	metrics = append(metrics, NewMetric("HighFree", meminfo["HighFree"]))

	metrics = append(metrics, NewMetric("LowTotal", meminfo["LowTotal"]))
	metrics = append(metrics, NewMetric("LowFree", meminfo["LowFree"]))

	metrics = append(metrics, NewMetric("SlabReclaimable", meminfo["SReclaimable"]))
	metrics = append(metrics, NewMetric("SlabUnreclaimable", meminfo["SUnreclaim"]))

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
