package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Modules struct {
	InputModules      []Module
	TransformModules  []Module
	OutputModules     []Module
	InputResponseChan chan *ModuleMetrics
	InputChannels     []chan int
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
	GetMetrics() (*ModuleMetrics, error)
}

type OutputModule interface {
	SendMetrics([]*ModuleMetrics) error
}

func getModules(config Config) Modules {
	files, err := ioutil.ReadDir(config.ConfigPath)
	if err != nil {
		log.Fatalf("Problem loading modules: %v", err)
	}

	modules := Modules{}

	modules.InputResponseChan = make(chan *ModuleMetrics, len(files)*2)

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
				requestChan, err := InitInputModule(module, modules.InputResponseChan)
				if err != nil {
					log.Fatalf("Failed to initialize input module: %v", err)
				}
				modules.InputChannels = append(modules.InputChannels, requestChan)
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

// InitInputModule creates a channel for making requests on and aggregates the response on the
// responseChan that is passed in.  It returns the request channel.
func InitInputModule(module Module, responseChan chan *ModuleMetrics) (chan int, error) {
	inputModule, ok := module.(InputModule)
	if !ok {
		return nil, errors.New("Not an input module")
	}
	requestChan := make(chan int)

	// Create a goroutine to listen for requests on the request channel
	go func(module InputModule, moduleName string, requestChan chan int, responseChan chan *ModuleMetrics) {
		for _ = range requestChan {
			metrics, err := module.GetMetrics()
			if err != nil {
				log.Printf("Failed to retrieve %s metrics: %v", moduleName, err)
			} else {
				responseChan <- metrics
			}
		}
	}(inputModule, module.Name(), requestChan, responseChan)

	return requestChan, nil
}

func getModule(name string) Module {
	switch name {
	case CpuModuleName:
		return &CPUInputModule{}
	case MemoryModuleName:
		return &MemoryInputModule{}
	case NetworkModuleName:
		return &NetworkInputModule{}
	case GraphiteModuleName:
		return &GraphiteOutputModule{}
	case DiskspaceModuleName:
		return &DiskspaceInputModule{}
	case DiskusageModuleName:
		return &DiskusageInputModule{}
	case LoadModuleName:
		return &LoadInputModule{}
	case ProcessessModuleName:
		return &ProcessesInputModule{}
	case InternalModuleName:
		return &InternalInputModule{}
	case RedisModuleName:
		return &RedisInputModule{}
	default:
		log.Fatalf("Invalid module: %s", name)
	}

	return nil
}

func tearDownModules(modules *Modules) {
	//close all channels
	for _, c := range modules.InputChannels {
		close(c)
	}
	close(modules.InputResponseChan)

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
