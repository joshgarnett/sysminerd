// +build linux

package main

import (
	"io/ioutil"
	"strings"
)

func (m *CPUInputModule) GetMetrics() ([]Metric, error) {

	b, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	content := string(b)

	return m.ParseProcStat(content)
}
