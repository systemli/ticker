package logger

import "github.com/sirupsen/logrus"

func NewLogrus(level, format string) *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	return logger
}
