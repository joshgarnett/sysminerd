package main

import (
	"golang.org/x/sys/unix" //see https://godoc.org/golang.org/x/sys/unix
)

const LoadModuleName = "loadavg"

// LinuxSysinfoLoadsScale magic number
const LinuxSysinfoLoadsScale = 65536.0

type LoadInputModule struct{}

func (m *LoadInputModule) Name() string {
	return LoadModuleName
}

func (m *LoadInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *LoadInputModule) TearDown() error {
	return nil
}

func (m *LoadInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 48)

	info := unix.Sysinfo_t{}
	err := unix.Sysinfo(&info)
	if err != nil {
		return nil, err
	}

	shortterm := float64(info.Loads[0]) / LinuxSysinfoLoadsScale
	midterm := float64(info.Loads[1]) / LinuxSysinfoLoadsScale
	longterm := float64(info.Loads[2]) / LinuxSysinfoLoadsScale

	metrics = append(metrics, NewMetric("shortterm", shortterm))
	metrics = append(metrics, NewMetric("midterm", midterm))
	metrics = append(metrics, NewMetric("longterm", longterm))

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
