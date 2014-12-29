package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	interval time.Duration
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

	config := Config{}
	config.interval = time.Duration(interval) * time.Second
	return config
}

func main() {
	config := parseConfig("config/sysminerd.yaml")

	//initialize all modules

	//start loop
	ticker := time.NewTicker(config.interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				tickModules(config)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	//catch sigkill, stop loop, cleanup
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range sigChan {
			log.Printf("captured %v, stopping metrics collection and exiting", sig)

			//run any cleanup steps
			close(quit)

			os.Exit(1)
		}
	}()

	// run forever
	select {}
}

func tickModules(config Config) {
	var max = config.interval.Seconds()
	var start = time.Now()

	// get metrics

	// transform metrics

	// send metrics

	// check to make sure the metrics collection isn't taking too long
	tickTime := time.Since(start).Seconds()
	if tickTime >= (max * .9) {
		log.Println("getInputMetrics took %f seconds", tickTime)
	}
}
