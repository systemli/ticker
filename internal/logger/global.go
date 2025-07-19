package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	instance *logrus.Logger
	mu       sync.RWMutex
)

// Initialize sets up the global logger instance with the specified level and format.
// This should be called once during application startup.
func Initialize(level, format string) *logrus.Logger {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = NewLogrus(level, format)
	} else {
		// Reconfigure existing instance
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			lvl = logrus.InfoLevel
		}
		instance.SetLevel(lvl)

		switch format {
		case "json":
			instance.SetFormatter(&logrus.JSONFormatter{})
		default:
			instance.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		}
	}

	return instance
}

// Get returns the global logger instance.
// If Initialize was not called before, it returns a default logger.
func Get() *logrus.Logger {
	mu.RLock()
	if instance != nil {
		defer mu.RUnlock()
		return instance
	}
	mu.RUnlock()

	// If not initialized, initialize with defaults
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = NewLogrus("info", "text")
	}
	return instance
}

// GetWithField returns a logger entry with the specified field.
// This is commonly used to add package information to log entries.
func GetWithField(key, value string) *logrus.Entry {
	return Get().WithField(key, value)
}

// GetWithPackage returns a logger entry with package field set.
// This is a convenience function for the common case of logging with package information.
func GetWithPackage(packageName string) *logrus.Entry {
	return GetWithField("package", packageName)
}
