package logger

import "github.com/sirupsen/logrus"

// NewLogger creates an common logging struct.
func NewLogger(context string) *logrus.Entry {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	entry := logrus.NewEntry(log)
	entry = entry.WithField("context", context)

	return entry
}
