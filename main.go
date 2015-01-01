package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Metric struct {
	module    string
	name      string
	value     float64
	timestamp time.Time
}

func main() {
	config := parseConfig("config/sysminerd.yaml")

	//initialize all modules
	modules := Modules{}

	modules.inputModules = append(modules.inputModules, CPUInputModule{})
	modules.outputModules = append(modules.outputModules, GraphiteOutputModule{})

	//start loop
	ticker := time.NewTicker(config.interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				tickModules(config, &modules)
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

func tickModules(config Config, modules *Modules) {
	var max = config.interval.Seconds()
	var start = time.Now()

	allMetrics := []Metric{}

	// get metrics
	for _, e := range modules.inputModules {
		module, ok := e.(InputModule)
		if !ok {
			log.Printf("%s is not an InputModule", e.Name())
		} else {
			metrics, err := module.GetMetrics()
			if err != nil {
				log.Printf("There was a problem getting metrics for %s, %v", e.Name(), err)
			} else {
				allMetrics = append(allMetrics, metrics...)
			}
		}
	}

	// transform metrics
	for _, e := range modules.transformModules {
		_, ok := e.(TransformModule)
		if !ok {
			log.Printf("%s is not an TransformModule", e.Name())
		}
	}

	// send metrics
	for _, e := range modules.outputModules {
		module, ok := e.(OutputModule)
		if !ok {
			log.Printf("%s is not an OutputModule", e.Name())
		} else {
			module.SendMetrics(allMetrics)
		}
	}

	// check to make sure the metrics collection isn't taking too long
	tickTime := time.Since(start).Seconds()
	if tickTime >= (max * .9) {
		log.Printf("getInputMetrics took %f seconds", tickTime)
	}
}
