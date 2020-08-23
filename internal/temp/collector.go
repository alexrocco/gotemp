package temp

import (
	"github.com/alexrocco/gotemp/internal/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Collector collects data temperature.
type Collector interface {
	// Collect the data
	Collect() (Data, error)
}

// SensorCollector implements Collector and collect data from the sensor.
type SensorCollector struct {
	Log    *logrus.Entry
	Sensor Sensor
}

// NewSensorCollector creates a collector for temperature.
func NewSensorCollector() *SensorCollector {
	collector := SensorCollector{
		Sensor: NewDHT22Sensor(),
		Log:    logger.NewLogger("sensor_collector"),
	}

	return &collector
}

// Collect collects the data using the DHT22 sensor plugged in the Raspberry PI.
func (sc *SensorCollector) Collect() (Data, error) {
	// read the temperature from the sensor
	data, err := sc.Sensor.Read()
	if err != nil {
		return Data{}, errors.Wrap(err, "error collecting data sensor")
	}

	sc.Log.Infof("Temperature: %v, Humidity: %v", data.Temperature, data.Humidity)

	return data, nil
}
