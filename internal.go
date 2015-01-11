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

		metrics = append(metrics, NewMetric("general.allocated", float64(memstats.Alloc)))
		metrics = append(metrics, NewMetric("general.system", float64(memstats.Sys)))
		metrics = append(metrics, NewMetric("general.lookups", float64(lookups)))
		metrics = append(metrics, NewMetric("general.mallocs", float64(mallocs)))
		metrics = append(metrics, NewMetric("general.frees", float64(frees)))

		metrics = append(metrics, NewMetric("heap.allocated", float64(memstats.HeapAlloc)))
		metrics = append(metrics, NewMetric("heap.system", float64(memstats.HeapSys)))
		metrics = append(metrics, NewMetric("heap.idle", float64(memstats.HeapIdle)))
		metrics = append(metrics, NewMetric("heap.in_use", float64(memstats.HeapInuse)))
		metrics = append(metrics, NewMetric("heap.released", float64(memstats.HeapReleased)))
		metrics = append(metrics, NewMetric("heap.objects", float64(memstats.HeapObjects)))

		numGC := int(memstats.NumGC - previous.NumGC)
		var averageGC float64
		if numGC > 0 {
			var total uint64
			for i := 0; i < numGC; i++ {
				total += memstats.PauseNs[(memstats.NumGC+255-uint32(i))%256]
			}
			averageGC = float64(total) / float64(numGC) / 1000000
		}

		metrics = append(metrics, NewMetric("gc.average_gc", averageGC))
		metrics = append(metrics, NewMetric("gc.num_gc", float64(numGC)))
	}

	m.lastStats = &memstats

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
