package main

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

// Config stores all the config options for sysminerd
type Config struct {
	interval time.Duration
	hostname string
}

func parseConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	yamlConfig := make(map[interface{}]interface{})

	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	interval, ok := yamlConfig["interval"].(int)
	if !ok || interval < 1 {
		log.Fatalf("Invalid interval specified: %v", yamlConfig["interval"])
	}

	hostname, _ := yamlConfig["hostname"].(string)

	config := Config{}
	config.interval = time.Duration(interval) * time.Second
	config.hostname = hostname
	return config
}
