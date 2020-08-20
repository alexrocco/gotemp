package main

import (
	"fmt"
	"log"

	"github.com/alexrocco/gotemp/internal/temp"
)

func main() {
	collector := temp.NewSensorCollector()

	data, err := collector.Collect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temp %v, Hum %v", data.Temperature, data.Humidity)
}
