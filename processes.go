package main

import (
	"io/ioutil"
	"strconv"
)

const ProcessessModuleName = "processes"

type ProcessesInputModule struct{}

func (m *ProcessesInputModule) Name() string {
	return ProcessessModuleName
}

func (m *ProcessesInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *ProcessesInputModule) TearDown() error {
	return nil
}

func (m *ProcessesInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 48)

	states := map[string]int{
		"running":  0,
		"sleeping": 0,
		"blocked":  0,
		"zombies":  0,
		"stopped":  0,
		"paging":   0,
	}

	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pid, err := strconv.ParseInt(file.Name(), 10, 0)
		if err != nil {
			continue
		}

		process, err := GetProcessStats(pid)
		if err != nil {
			continue
		}

		states[process.State]++
	}

	for state, total := range states {
		metrics = append(metrics, NewMetric(state, float64(total)))
	}

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}
