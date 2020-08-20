package temp

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

// Collector collects data temperature
type Collector interface {
	// Collect the data
	Collect() (Data, error)
}

// SensorCollector implements Collector and collect data from the sensor
type SensorCollector struct {
	logger log.Logger
	sensor Sensor
}

// NewSensorCollector creates a colector for temperature
func NewSensorCollector() *SensorCollector {
	colector := SensorCollector{
		sensor: NewDHT22Sensor(),
		logger: *log.New(os.Stdout, "SENSOR_COLLECTOR ", log.Ltime),
	}

	return &colector
}

// Collect collects the data using the DHT22 sensor plugged in the Raspberry PI
func (sc *SensorCollector) Collect() (Data, error) {
	// read the temperature from the sensor
	data, err := sc.sensor.Read()
	if err != nil {
		return Data{}, errors.Wrap(err, "error collecting data sensor")
	}

	sc.logger.Printf("Temperature: %v, Humidity: %v", data.Temperature, data.Humidity)

	return data, nil
}
