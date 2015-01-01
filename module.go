package main

type Modules struct {
	inputModules     []Module
	transformModules []Module
	outputModules    []Module
}

type Module interface {
	Init(config map[interface{}]interface{}) error
	Name() string
	TearDown() error
}

type TransformModule interface {
	TransformMetrics(metrics []Metric) ([]Metric, error)
}

type InputModule interface {
	GetMetrics() ([]Metric, error)
}

type OutputModule interface {
	SendMetrics(metrics []Metric) ([]Metric, error)
}
