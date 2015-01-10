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

	metrics = append(metrics, NewMetric("CurrentLoad", float64(memData.dwMemoryLoad)))
	metrics = append(metrics, NewMetric("Total", float64(memData.ullTotalPhys)))
	metrics = append(metrics, NewMetric("Free", float64(memData.ullAvailPhys)))

	used := float64(memData.ullTotalPhys) - float64(memData.ullAvailPhys)

	metrics = append(metrics, NewMetric("Used", used))

	virTotal := float64(memData.ullTotalVirtual)
	virFree := float64(memData.ullAvailVirtual)
	virUsed := virTotal - virFree

	metrics = append(metrics, NewMetric("TotalVirtual", virTotal))
	metrics = append(metrics, NewMetric("FreeVirtual", virFree))
	metrics = append(metrics, NewMetric("UsedVirtual", virUsed))

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
