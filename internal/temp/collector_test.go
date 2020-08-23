package temp_test

import (
	"errors"
	"testing"

	"github.com/alexrocco/gotemp/internal/logger"
	"github.com/alexrocco/gotemp/internal/temp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSensor struct {
	mock.Mock
}

func (ms *mockSensor) Read() (temp.Data, error) {
	args := ms.Called()

	return args.Get(0).(temp.Data), args.Error(1)
}

func TestSensorCollector_Collect(t *testing.T) {
	t.Run("It should return data when the sensor is ok", func(t *testing.T) {
		mockSensor := mockSensor{}

		data := temp.Data{
			Humidity:    50.1,
			Temperature: 25.5,
		}

		mockSensor.On("Read").Return(data, nil)

		collector := temp.SensorCollector{
			Log:    logger.NewLogger("test"),
			Sensor: &mockSensor,
		}

		result, err := collector.Collect()

		assert.Nil(t, err)
		assert.Equal(t, data, result)
	})
	t.Run("It should fail when sensor is not ok", func(t *testing.T) {
		mockSensor := mockSensor{}

		//nolint:goerr113
		mockSensor.On("Read").Return(temp.Data{}, errors.New("some error"))

		collector := temp.SensorCollector{
			Log:    logger.NewLogger("test"),
			Sensor: &mockSensor,
		}

		_, err := collector.Collect()

		assert.NotNil(t, err)
	})
}
