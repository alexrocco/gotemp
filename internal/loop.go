package internal

import (
	"fmt"
	"time"

	"github.com/alexrocco/gotemp/internal/logger"
	"github.com/alexrocco/gotemp/internal/temp"
	"github.com/alexrocco/gotemp/internal/timeseries"
	"github.com/sirupsen/logrus"
)

// Loop interface describe how to start/stop looping.
type Loop interface {
	Start()
	Stop()
}

// CollectWeatherLoop holds the ticker for the weather collection data.
type CollectWeatherLoop struct {
	log    *logrus.Entry
	ticker *time.Ticker
	done   chan bool
	tags   map[string]string

	sensorCollector *temp.SensorCollector
	influxDBSender  *timeseries.InfluxDBSender
}

// NewCollectWeatherLoop creates a new loop.
func NewCollectWeatherLoop(d time.Duration, influxDBAddr string, tags map[string]string) *CollectWeatherLoop {
	return &CollectWeatherLoop{
		ticker: time.NewTicker(d),
		tags:   tags,

		done: make(chan bool),
		log:  logger.NewLogger("collect_weather_loop"),

		sensorCollector: temp.NewSensorCollector(),
		influxDBSender:  timeseries.NewInfluxDBSender(influxDBAddr),
	}
}

// Start starts the loop.
func (cwl *CollectWeatherLoop) Start() {
	for {
		select {
		case <-cwl.ticker.C:
			// collect the weather
			data, err := cwl.sensorCollector.Collect()
			if err != nil {
				cwl.log.Error(err)

				continue
			}

			values := map[string]string{
				"humidity":    fmt.Sprintf("%.2f", data.Humidity),
				"temperature": fmt.Sprintf("%.2f", data.Temperature),
			}

			tags := map[string]string{
				"room": "baby",
			}

			// push to InfluxDB
			err = cwl.influxDBSender.Send("weather", values, tags)
			if err != nil {
				cwl.log.Error(err)
			}
		case <-cwl.done:
			return
		}
	}
}

// Stop stops the loop.
func (cwl *CollectWeatherLoop) Stop() {
	cwl.log.Info("Stopping the loop")
	cwl.done <- true
}
