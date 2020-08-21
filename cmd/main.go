package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexrocco/gotemp/internal"
)

func main() {
	cwl := internal.NewCollectWeatherLoop(10 * time.Second)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		// Wait for the signal
		<-sigc
		// stop the loop
		cwl.Stop()
	}()

	// start the loop
	cwl.Start()
}
