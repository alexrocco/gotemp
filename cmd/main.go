package main

import (
	"os"
	"os/signal"
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

	cwl := internal.NewCollectWeatherLoop(ticker * time.Second)

	sigNot := make(chan os.Signal, 1)
	// catch commons signais when stopping a process
	signal.Notify(sigNot, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Wait for the signal
		sig := <-sigNot

		log.Infof("Received signal: %v", sig)

		// stop the loop
		cwl.Stop()
	}()

	// start the loop
	cwl.Start()
}
