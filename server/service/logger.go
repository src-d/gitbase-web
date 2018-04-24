package service

import (
	"github.com/sirupsen/logrus"
)

// NewLogger returns a logrus Logger
func NewLogger(env string) *logrus.Logger {
	logger := logrus.New()

	if env == "dev" {
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.Formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
			},
		}
		logger.SetLevel(logrus.WarnLevel)
	}

	return logger
}
