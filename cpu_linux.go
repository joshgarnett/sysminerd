// +build linux

package main

import (
	"io/ioutil"
)

func (m *CPUInputModule) GetMetrics() (*ModuleMetrics, error) {

	b, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	content := string(b)

	return m.ParseProcStat(content)
}
