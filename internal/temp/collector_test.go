package temp

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSensor struct {
	mock.Mock
}

func (ms *mockSensor) Read() (Data, error) {
	args := ms.Called()

	return args.Get(0).(Data), args.Error(1)
}

func TestSensorCollector_Collect(t *testing.T) {
	t.Run("It should return data when the sensor is ok", func(t *testing.T) {
		mockSensor := mockSensor{}

		data := Data{
			Humidity:    50.1,
			Temperature: 25.5,
		}

		mockSensor.On("Read").Return(data, nil)

		collector := SensorCollector{
			logger: *log.New(os.Stdout, "TEST", log.Lmicroseconds),
			sensor: &mockSensor,
		}

		result, err := collector.Collect()

		assert.Nil(t, err)
		assert.Equal(t, data, result)
	})
	t.Run("It should fail when sensor is not ok", func(t *testing.T) {
		mockSensor := mockSensor{}

		mockSensor.On("Read").Return(Data{}, errors.New("some error"))

		collector := SensorCollector{
			logger: *log.New(os.Stdout, "TEST", log.Lmicroseconds),
			sensor: &mockSensor,
		}

		_, err := collector.Collect()

		assert.NotNil(t, err)
	})
}
