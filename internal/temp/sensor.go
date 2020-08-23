package temp

import (
	"github.com/d2r2/go-dht"
	"github.com/pkg/errors"
)

// Data holds sensor information.
type Data struct {
	Temperature float32
	Humidity    float32
}

// Sensor interface reads data from the sensor plugged. It was created to wrap go-dht package for better testing.
type Sensor interface {
	Read() (Data, error)
}

// DHT22Sensor implements Sensor interface for DHT22 sensor.
type DHT22Sensor struct {
}

// NewDHT22Sensor creates a sensor for DHT22.
func NewDHT22Sensor() *DHT22Sensor {
	return &DHT22Sensor{}
}

func (s *DHT22Sensor) Read() (Data, error) {
	// read the temperature from the sensor
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 2, false, 10)
	if err != nil {
		return Data{}, errors.Wrap(err, "error reading data from DHT22 on pin 2")
	}

	return Data{
		Temperature: temperature,
		Humidity:    humidity,
	}, nil
}
