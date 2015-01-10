package main

import (
	"time"
)

type ModuleMetrics struct {
	Module  string
	Metrics []Metric
}

// Metric stores all data related to a system metric
type Metric struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

// NewMetric creates a new Metric structure and automatically sets the timestamp
func NewMetric(name string, value float64) Metric {
	m := Metric{}
	m.Name = name
	m.Value = value
	m.Timestamp = time.Now()
	return m
}
