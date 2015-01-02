package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

type Modules struct {
	InputModules     []Module
	TransformModules []Module
	OutputModules    []Module
}

type Module interface {
	Init(config *Config, moduleConfig *ModuleConfig) error
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

func getModules(config Config) Modules {
	files, err := ioutil.ReadDir(config.ConfigPath)
	if err != nil {
		log.Fatalf("Problem loading modules: %v", err)
	}

	modules := Modules{}

	for _, file := range files {
		if file.IsDir() {
			log.Printf("Found subdirectory %s in config path %s", file.Name(), config.ConfigPath)
			continue
		}

		extension := file.Name()[len(file.Name())-4:]
		if extension != "yaml" {
			log.Printf("Found non yaml extension %s in config path %s", extension, config.ConfigPath)
			continue
		}

		fullPath := filepath.Join(config.ConfigPath, file.Name())

		moduleConfig := parseModuleConfig(fullPath)

		if moduleConfig.Enabled {
			module := getModule(moduleConfig.Name)

			log.Printf("Module %s enabled: %v", moduleConfig.Name, moduleConfig)

			switch v := module.(type) {
			default:
				log.Fatalf("unexpected type %T", v)
			case InputModule:
				modules.InputModules = append(modules.InputModules, module)
			case TransformModule:
				modules.TransformModules = append(modules.TransformModules, module)
			case OutputModule:
				modules.OutputModules = append(modules.OutputModules, module)
			}

			module.Init(&config, &moduleConfig)
		}
	}

	return modules
}

func getModule(name string) Module {
	switch name {
	case "cpu":
		return &CPUInputModule{}
	case "graphite":
		return &GraphiteOutputModule{}
	default:
		log.Fatalf("Invalid module: %s", name)
	}

	return nil
}

func tearDownModules(modules *Modules) {
	// input modules
	for _, e := range modules.InputModules {
		e.TearDown()
	}

	// transform modules
	for _, e := range modules.TransformModules {
		e.TearDown()
	}

	// output modules
	for _, e := range modules.OutputModules {
		e.TearDown()
	}
}
