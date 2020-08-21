package internal

import (
	"fmt"
	"time"

	"github.com/alexrocco/gotemp/internal/temp"
	"github.com/alexrocco/gotemp/internal/timeseries"
)

// Loop interface describe how to keep looping
type Loop interface {
	Start()
}

// CollectWeatherloop holds the ticker for the weather collection data
type CollectWeatherloop struct {
	ticker *time.Ticker
	done   *chan bool
}

// NewLoop creates a new loop
func NewLoop(d time.Duration, done *chan bool) *CollectWeatherloop {
	return &CollectWeatherloop{
		ticker: time.NewTicker(d),
		done:   done,
	}
}

// Start starts CollectWeatherloop
func (cwl *CollectWeatherloop) Start() {
	for {
		select {
		case <-*cwl.done:
			return
		case <-cwl.ticker.C:
			// collect the weather
			sensorCollector := temp.NewSensorCollector()
			data, err := sensorCollector.Collect()
			if err != nil {
				fmt.Println(err)
				continue
			}

			influxDBSender := timeseries.NewInfluxDBSender(":8086")

			values := map[string]string{
				"humidity":    fmt.Sprintf("%f", data.Humidity),
				"temperature": fmt.Sprintf("%f", data.Temperature),
			}

			tags := map[string]string{
				"room": "baby",
			}

			// push to InfluxDB
			err = influxDBSender.Send("weather", values, tags)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
