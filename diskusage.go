package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type DiskStats struct {
	Reads                  float64
	ReadsMerged            float64
	ReadsSectors           float64
	ReadsMilliseconds      float64
	Writes                 float64
	WritesMerged           float64
	WritesSectors          float64
	WritesMilliseconds     float64
	IOInProgress           float64
	IOMilliseconds         float64
	IOMillisecondsWeighted float64
}

const DiskusageModuleName = "diskusage"

type DiskusageInputModule struct {
	previousDiskStats map[string]DiskStats
	previousTime      time.Time
}

func (m *DiskusageInputModule) Name() string {
	return DiskusageModuleName
}

func (m *DiskusageInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *DiskusageInputModule) TearDown() error {
	return nil
}

func (m *DiskusageInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 50)
	now := time.Now()
	timeDiff := time.Since(m.previousTime).Seconds()

	allStats, err := GetDiskStats("/proc/diskstats")
	if err != nil {
		return nil, err
	}

	if m.previousDiskStats != nil {
		for device, stats := range allStats {
			previous, ok := m.previousDiskStats[device]

			if ok {
				readsPerSecond := (stats.Reads - previous.Reads) / timeDiff
				writesPerSecond := (stats.Writes - previous.Writes) / timeDiff
				readsMergedPerSecond := (stats.ReadsMerged - previous.ReadsMerged) / timeDiff
				writesMergedPerSecond := (stats.WritesMerged - previous.WritesMerged) / timeDiff
				readBytesPerSecond := ((stats.ReadsSectors - previous.ReadsSectors) * 512) / timeDiff
				writeBytesPerSecond := ((stats.WritesSectors - previous.WritesSectors) * 512) / timeDiff

				metrics = append(metrics, NewMetric(fmt.Sprintf("%s.reads", device), readsPerSecond))
				metrics = append(metrics, NewMetric(fmt.Sprintf("%s.writes", device), writesPerSecond))
				metrics = append(metrics, NewMetric(fmt.Sprintf("%s.reads_merged", device), readsMergedPerSecond))
				metrics = append(metrics, NewMetric(fmt.Sprintf("%s.writes_merged", device), writesMergedPerSecond))
				metrics = append(metrics, NewMetric(fmt.Sprintf("%s.read_bytes", device), readBytesPerSecond))
				metrics = append(metrics, NewMetric(fmt.Sprintf("%s.write_bytes", device), writeBytesPerSecond))
			}
		}
	}

	m.previousDiskStats = allStats
	m.previousTime = now

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}

func GetDiskStats(path string) (map[string]DiskStats, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")

	stats := make(map[string]DiskStats)

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) < 14 || len(fields) == 0 {
			continue
		}

		device := fields[2]

		if strings.HasPrefix(device, "ram") || strings.HasPrefix(device, "loop") {
			continue
		}

		stats[device] = DiskStats{
			Reads:                  fieldValue(fields[3]),
			ReadsMerged:            fieldValue(fields[4]),
			ReadsSectors:           fieldValue(fields[5]),
			ReadsMilliseconds:      fieldValue(fields[6]),
			Writes:                 fieldValue(fields[7]),
			WritesMerged:           fieldValue(fields[8]),
			WritesSectors:          fieldValue(fields[9]),
			WritesMilliseconds:     fieldValue(fields[10]),
			IOInProgress:           fieldValue(fields[11]),
			IOMilliseconds:         fieldValue(fields[12]),
			IOMillisecondsWeighted: fieldValue(fields[13]),
		}
	}

	return stats, nil
}

func fieldValue(field string) float64 {
	value, err := strconv.ParseFloat(field, 64)
	if err != nil {
		value = 0
	}
	return value
}
