package main

import (
	"time"

	"github.com/alexrocco/gotemp/internal"
)

func main() {
	done := make(chan bool)
	cwl := internal.NewLoop(10*time.Second, &done)

	cwl.Start()
}
