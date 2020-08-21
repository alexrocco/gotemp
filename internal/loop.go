package internal

import (
	"fmt"
	"time"

	"github.com/alexrocco/gotemp/internal/temp"
	"github.com/alexrocco/gotemp/internal/timeseries"
)

// Loop interface describe how to start/stop looping
type Loop interface {
	Start()
	Stop()
}

// CollectWeatherLoop holds the ticker for the weather collection data
type CollectWeatherLoop struct {
	ticker *time.Ticker
	done   chan bool
}

// NewCollectWeatherLoop creates a new loop
func NewCollectWeatherLoop(d time.Duration) *CollectWeatherLoop {
	return &CollectWeatherLoop{
		ticker: time.NewTicker(d),
		done:   make(chan bool),
	}
}

// Start starts CollectWeatherloop
func (cwl *CollectWeatherLoop) Start() {
	for {
		select {
		case <-cwl.ticker.C:
			// collect the weather
			sensorCollector := temp.NewSensorCollector()
			data, err := sensorCollector.Collect()
			if err != nil {
				fmt.Println(err)
				continue
			}

			influxDBSender := timeseries.NewInfluxDBSender(":8089")

			values := map[string]string{
				"humidity":    fmt.Sprintf("%.2f", data.Humidity),
				"temperature": fmt.Sprintf("%.2f", data.Temperature),
			}

			tags := map[string]string{
				"room": "baby",
			}

			// push to InfluxDB
			err = influxDBSender.Send("weather", values, tags)
			if err != nil {
				fmt.Println(err)
			}
		case <-cwl.done:
			return
		}
	}
}

// Stop stops the loop
func (cwl *CollectWeatherLoop) Stop() {
	fmt.Println("Stopping the loop")
	cwl.done <- true
}
