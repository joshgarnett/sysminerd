package main

type Module interface {
	init(config map[interface{}]interface{}) error
	name() string
	tearDown() error
}

type InputModule interface {
	getMetrics() ([]Metric, error)
}

type OutputModule interface {
	sendMetrics(metrics []Metric) ([]Metric, error)
}
