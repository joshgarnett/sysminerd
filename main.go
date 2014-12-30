package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Metric struct {
	name      string
	value     float64
	timestamp time.Time
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
		log.Printf("getInputMetrics took %f seconds", tickTime)
	}
}
