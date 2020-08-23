package timeseries

import (
	"fmt"
	"github.com/alexrocco/gotemp/internal/logger"
	"net"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// InfluxDBSender holds configuration for InfluxDB
type InfluxDBSender struct {
	address string
	log     *logrus.Entry
}

// NewInfluxDBSender creates a InfluxDBSender
func NewInfluxDBSender(address string) *InfluxDBSender {
	return &InfluxDBSender{
		address: address,
		log:     logger.NewLogger("influxdb_sender"),
	}
}

// Send sends time-series data to InfluxDB
func (is *InfluxDBSender) Send(measurement string, values map[string]string, tags map[string]string) error {
	if len(measurement) == 0 || len(values) == 0 {
		return errors.New("measurement and values should not be empty")
	}

	// Format tags
	tagsf := ""
	if len(tags) > 0 {
		for _, entry := range sortMap(tags) {
			tagsf = fmt.Sprintf("%s,%s=%s", tagsf, entry.key, entry.value)
		}
	}

	// Format values
	valuesF := ""
	for _, entry := range sortMap(values) {
		valuesF = fmt.Sprintf("%s,%s=%s", valuesF, entry.key, entry.value)
	}
	// remove comma
	valuesF = valuesF[1:]

	// create the point for the time-series database
	point := fmt.Sprintf("%s%s %s %d", measurement, tagsf, valuesF, time.Now().UnixNano())

	// open UDP connection
	conn, err := net.Dial("udp", is.address)
	if err != nil {
		return errors.Wrap(err, "cannot open UDP connection")
	}
	defer conn.Close()

	is.log.Infof("Sending point to InfluxDB: %s", point)
	_, err = conn.Write([]byte(point))
	if err != nil {
		return errors.Wrap(err, "could not write on InfluxUDP port")
	}

	return nil
}

type entry struct {
	key   string
	value string
}

// sortMap sort the map and convert to a slice of entris since map is an unordered collection
func sortMap(values map[string]string) []entry {
	keys := make([]string, 0)
	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	result := make([]entry, 0)
	for _, k := range keys {
		result = append(result, entry{
			key:   k,
			value: values[k],
		})
	}

	return result
}
