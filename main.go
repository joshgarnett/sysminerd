package main

import (
	"flag"
	"fmt"
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

var configFile = flag.String("c", "", "config file to use")

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()

	if *configFile == "" {
		usage()
	}

	config := parseConfig(*configFile)

	//get all modules
	modules := getModules(config)

	//start loop
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
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

			tearDownModules(&modules)

			os.Exit(0)
		}
	}()

	// run forever
	select {}
}

func tickModules(config Config, modules *Modules) {
	var start = time.Now()

	allMetrics := []Metric{}

	// send metric requests to the input modules in a non blocking manner
	for i, c := range modules.InputChannels {
		module := modules.InputModules[i]
		select {
		case c <- 1:
		default:
			log.Printf("The %s input module is queuing requests", module.Name())
		}
	}

	// collect all metrics that are available
	collectMetrics := true
	for collectMetrics {
		select {
		case metrics := <-modules.InputResponseChan:
			allMetrics = append(allMetrics, metrics...)
		default:
			collectMetrics = false
		}
	}

	// transform metrics
	for _, e := range modules.TransformModules {
		_, ok := e.(TransformModule)
		if !ok {
			log.Printf("%s is not an TransformModule", e.Name())
		}
	}

	// send metrics
	for _, e := range modules.OutputModules {
		module, ok := e.(OutputModule)
		if !ok {
			log.Printf("%s is not an OutputModule", e.Name())
		} else {
			module.SendMetrics(allMetrics)
		}
	}

	// check to make sure the metrics collection isn't taking too long
	maxTime := float64(config.Interval) * 0.9
	tickTime := time.Since(start).Seconds()
	if tickTime >= maxTime {
		log.Printf("getInputMetrics took %f seconds", tickTime)
	}
}
