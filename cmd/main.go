package main

import (
	"github.com/alexrocco/gotemp/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexrocco/gotemp/internal"
)

func main() {
	log := logger.NewLogger("main")

	log.Info("### Starting gotemp app ###")
	defer log.Info("### Stopping gotemp app ###")

	cwl := internal.NewCollectWeatherLoop(10 * time.Second)

	sigNot := make(chan os.Signal, 1)
	// catch commons signais when stopping a process
	signal.Notify(sigNot, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

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
