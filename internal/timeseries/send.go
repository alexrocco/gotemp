package timeseries

// Sender send time-series values
type Sender interface {
	Send(measurement string, values map[string]string, tags map[string]string) error
}
