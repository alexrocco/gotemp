package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alexrocco/gotemp/internal"
	"github.com/alexrocco/gotemp/internal/logger"
)

const ticker = 10

func main() {
	log := logger.NewLogger("main")

	log.Info("### Starting gotemp app ###")
	defer log.Info("### Stopping gotemp app ###")

	// Command-lines flags
	tagsF := flag.String("tags", "", "tags used in the metrics")
	influxDBAddrF := flag.String("influxdb", ":8089", "InfluxDB address")

	// Parse the flags before using its values
	flag.Parse()

	tags := make(map[string]string)

	// Create the tags from command-line flags
	if len(*tagsF) > 0 {
		tagsF := strings.Split(*tagsF, ",")

		for _, t := range tagsF {
			kv := strings.Split(t, "=")
			tags[kv[0]] = kv[1]
		}
	}

	log.Infof("Flags - tags: %v, influxdb: %s", tags, *influxDBAddrF)

	cwl := internal.NewCollectWeatherLoop(ticker*time.Second, *influxDBAddrF, tags)

	sigNot := make(chan os.Signal, 1)
	// catch commons signais when stopping a process
	signal.Notify(sigNot, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Wait for the signal
		sig := <-sigNot

		log.Warnf("Received signal: %v", sig)

		// stop the loop
		cwl.Stop()
	}()

	// start the loop
	cwl.Start()
}
