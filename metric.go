package main

import (
	"time"
)

// Metric stores all data related to a system metric
type Metric struct {
	module    string
	name      string
	value     float64
	timestamp time.Time
}
