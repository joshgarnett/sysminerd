package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config stores all the config options for sysminerd
type Config struct {
	Interval   float64
	Hostname   string
	ConfigPath string `yaml:"config_path"`
}

type ModuleConfig struct {
	Name     string
	Enabled  bool
	Settings map[string]string
}

func parseConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	yamlConfig := Config{}

	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	log.Printf("Config: %v", yamlConfig)

	return yamlConfig
}

func parseModuleConfig(path string) ModuleConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	moduleConfig := ModuleConfig{}
	err = yaml.Unmarshal(data, &moduleConfig)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	return moduleConfig
}
