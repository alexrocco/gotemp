package timeseries_test

import (
	"net"
	"regexp"
	"testing"

	"github.com/alexrocco/gotemp/internal/timeseries"
	"github.com/stretchr/testify/assert"
)

const (
	udpAddr = ":8180"
)

//nolint:funlen
func TestInfluxDBSender_Send(t *testing.T) {
	t.Run("It should send points to influxdb", func(t *testing.T) {
		type testData struct {
			regexMsg    string
			measurement string
			values      map[string]string
			tags        map[string]string
		}
		tests := []testData{
			{
				regexMsg:    "weather,house=123,room=baby humidity=50,temperature=21\\.5 \\d+$",
				measurement: "weather",
				values: map[string]string{
					"temperature": "21.5",
					"humidity":    "50",
				},
				tags: map[string]string{
					"room":  "baby",
					"house": "123",
				},
			},
			{
				regexMsg:    "weather humidity=50,temperature=21 \\d+$",
				measurement: "weather",
				values: map[string]string{
					"temperature": "21",
					"humidity":    "50",
				},
				tags: nil,
			},
			{
				regexMsg:    "test value=10 \\d+$",
				measurement: "test",
				values: map[string]string{
					"value": "10",
				},
				tags: nil,
			},
		}

		// Start UDP server to mock InfluxDB server
		ludpAddr, err := net.ResolveUDPAddr("udp", udpAddr)
		assert.Nil(t, err)

		server, err := net.ListenUDP("udp", ludpAddr)
		assert.Nil(t, err)
		defer server.Close()

		for _, tt := range tests {
			influxDBSender := timeseries.NewInfluxDBSender(udpAddr)

			err = influxDBSender.Send(tt.measurement, tt.values, tt.tags)
			assert.Nil(t, err)

			// Read the message from the UDP server
			buf := make([]byte, 2048)
			n, _, err := server.ReadFromUDP(buf)
			if err != nil {
				t.Fatal("error reading message from UDP")
			}

			// Assert the message from the UDP server and ignore timestamp
			assert.Regexp(t, regexp.MustCompile(tt.regexMsg), string(buf[:n]))
		}
	})
	t.Run("It should fail when measurement is empty", func(t *testing.T) {
		influxDBSender := timeseries.NewInfluxDBSender(udpAddr)
		err := influxDBSender.Send("", nil, nil)
		assert.NotNil(t, err)
	})
	t.Run("It should fail when values is empty", func(t *testing.T) {
		influxDBSender := timeseries.NewInfluxDBSender(udpAddr)
		err := influxDBSender.Send("test", nil, nil)
		assert.NotNil(t, err)
	})
}
