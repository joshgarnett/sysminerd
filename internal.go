package main

import (
	"runtime"
)

const InternalModuleName = "internal"

type InternalInputModule struct {
	lastStats *runtime.MemStats
}

func (m *InternalInputModule) Name() string {
	return InternalModuleName
}

func (m *InternalInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *InternalInputModule) TearDown() error {
	return nil
}

func (m *InternalInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0)

	memstats := runtime.MemStats{}

	runtime.ReadMemStats(&memstats)

	if m.lastStats != nil {
		previous := m.lastStats

		lookups := memstats.Lookups - previous.Lookups
		mallocs := memstats.Mallocs - previous.Mallocs
		frees := memstats.Frees - previous.Frees

		metrics = append(metrics, NewMetric("general.Alloc", float64(memstats.Alloc)))
		metrics = append(metrics, NewMetric("general.Sys", float64(memstats.Sys)))
		metrics = append(metrics, NewMetric("general.Lookups", float64(lookups)))
		metrics = append(metrics, NewMetric("general.Mallocs", float64(mallocs)))
		metrics = append(metrics, NewMetric("general.Frees", float64(frees)))

		metrics = append(metrics, NewMetric("heap.Alloc", float64(memstats.HeapAlloc)))
		metrics = append(metrics, NewMetric("heap.Sys", float64(memstats.HeapSys)))
		metrics = append(metrics, NewMetric("heap.Idle", float64(memstats.HeapIdle)))
		metrics = append(metrics, NewMetric("heap.Inuse", float64(memstats.HeapInuse)))
		metrics = append(metrics, NewMetric("heap.Released", float64(memstats.HeapReleased)))
		metrics = append(metrics, NewMetric("heap.Objects", float64(memstats.HeapObjects)))

		numGC := int(memstats.NumGC - previous.NumGC)
		var averageGC float64
		if numGC > 0 {
			var total uint64
			for i := 0; i < numGC; i++ {
				total += memstats.PauseNs[(memstats.NumGC+255-uint32(i))%256]
			}
			averageGC = float64(total) / float64(numGC) / 1000000
		}

		metrics = append(metrics, NewMetric("gc.PauseTotalNs", float64(memstats.PauseTotalNs)))
		metrics = append(metrics, NewMetric("gc.AverageGCMs", averageGC))
		metrics = append(metrics, NewMetric("gc.NumGC", float64(numGC)))
	}

	m.lastStats = &memstats

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
