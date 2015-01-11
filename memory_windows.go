// +build windows

package main

import (
	"unsafe"
)

func (m *MemoryInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 50)

	var memData memorystatusex
	memData.dwLength = dword(unsafe.Sizeof(memData))

	r1, _, err := _globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memData)))

	if r1 != 1 {
		return nil, err
	}

	metrics = append(metrics, NewMetric("current_load", float64(memData.dwMemoryLoad)))
	metrics = append(metrics, NewMetric("total", float64(memData.ullTotalPhys)))
	metrics = append(metrics, NewMetric("free", float64(memData.ullAvailPhys)))

	used := float64(memData.ullTotalPhys) - float64(memData.ullAvailPhys)

	metrics = append(metrics, NewMetric("used", used))

	virTotal := float64(memData.ullTotalVirtual)
	virFree := float64(memData.ullAvailVirtual)
	virUsed := virTotal - virFree

	metrics = append(metrics, NewMetric("total_virtual", virTotal))
	metrics = append(metrics, NewMetric("free_virtual", virFree))
	metrics = append(metrics, NewMetric("used_virtual", virUsed))

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
