package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
	buffer *bytes.Buffer
}

func (s *LoggerTestSuite) SetupTest() {
	// Reset the global logger instance before each test
	mu.Lock()
	instance = nil
	mu.Unlock()

	// Create a buffer to capture log output
	s.buffer = new(bytes.Buffer)
}

func (s *LoggerTestSuite) TestInitialization() {
	s.Run("Initialize with debug level and JSON format", func() {
		logger := Initialize("debug", "json")

		s.NotNil(logger)
		s.Equal(logrus.DebugLevel, logger.Level)
		s.IsType(&logrus.JSONFormatter{}, logger.Formatter)
	})

	s.Run("Initialize with info level and text format", func() {
		logger := Initialize("info", "text")

		s.NotNil(logger)
		s.Equal(logrus.InfoLevel, logger.Level)
		s.IsType(&logrus.TextFormatter{}, logger.Formatter)
	})

	s.Run("Initialize with invalid level defaults to info", func() {
		logger := Initialize("invalid", "text")

		s.NotNil(logger)
		s.Equal(logrus.InfoLevel, logger.Level)
	})

	s.Run("Multiple Initialize calls use the last configuration", func() {
		// First initialization
		logger1 := Initialize("error", "json")
		s.Equal(logrus.ErrorLevel, logger1.Level)
		s.IsType(&logrus.JSONFormatter{}, logger1.Formatter)

		// Second initialization should override
		logger2 := Initialize("debug", "text")
		s.Equal(logrus.DebugLevel, logger2.Level)
		s.IsType(&logrus.TextFormatter{}, logger2.Formatter)

		// Both should be the same instance with the latest configuration
		s.Same(logger1, logger2)
		s.Equal(logrus.DebugLevel, logger1.Level) // logger1 should also have the new config
		s.IsType(&logrus.TextFormatter{}, logger1.Formatter)
	})
}

func (s *LoggerTestSuite) TestGet() {
	s.Run("Get without Initialize returns default logger", func() {
		logger := Get()

		s.NotNil(logger)
		s.Equal(logrus.InfoLevel, logger.Level)
		s.IsType(&logrus.TextFormatter{}, logger.Formatter)
	})

	s.Run("Get after Initialize returns configured logger", func() {
		Initialize("warn", "json")
		logger := Get()

		s.NotNil(logger)
		s.Equal(logrus.WarnLevel, logger.Level)
		s.IsType(&logrus.JSONFormatter{}, logger.Formatter)
	})

	s.Run("Multiple Get calls return same instance", func() {
		logger1 := Get()
		logger2 := Get()

		s.Equal(logger1, logger2)
	})
}

func (s *LoggerTestSuite) TestGetWithField() {
	s.Run("GetWithField adds specified field", func() {
		Initialize("debug", "json")
		entry := GetWithField("test_key", "test_value")

		s.NotNil(entry)
		s.Equal("test_value", entry.Data["test_key"])
	})

	s.Run("GetWithField with multiple fields", func() {
		entry1 := GetWithField("key1", "value1")
		entry2 := entry1.WithField("key2", "value2")

		s.Equal("value1", entry1.Data["key1"])
		s.Equal("value1", entry2.Data["key1"])
		s.Equal("value2", entry2.Data["key2"])
	})
}

func (s *LoggerTestSuite) TestGetWithPackage() {
	s.Run("GetWithPackage adds package field", func() {
		Initialize("debug", "json")
		entry := GetWithPackage("test_package")

		s.NotNil(entry)
		s.Equal("test_package", entry.Data["package"])
	})

	s.Run("GetWithPackage with different package names", func() {
		entry1 := GetWithPackage("package1")
		entry2 := GetWithPackage("package2")

		s.Equal("package1", entry1.Data["package"])
		s.Equal("package2", entry2.Data["package"])
	})
}

func (s *LoggerTestSuite) TestLogging() {
	s.Run("Debug messages respect log level", func() {
		logger := Initialize("debug", "json")
		logger.SetOutput(s.buffer)

		entry := GetWithPackage("test")
		entry.Debug("debug message")

		output := s.buffer.String()
		s.Contains(output, "debug message")
		s.Contains(output, `"level":"debug"`)
		s.Contains(output, `"package":"test"`)
	})

	s.Run("Debug messages filtered when level is info", func() {
		// Reset buffer and create fresh logger instance
		s.buffer.Reset()
		logger := Initialize("info", "json")
		logger.SetOutput(s.buffer)

		entry := GetWithPackage("test")
		entry.Debug("debug message")
		entry.Info("info message")

		output := s.buffer.String()
		s.NotContains(output, "debug message")
		s.Contains(output, "info message")
	})

	s.Run("Error messages always shown", func() {
		logger := Initialize("error", "json")
		logger.SetOutput(s.buffer)

		entry := GetWithPackage("test")
		entry.Error("error message")

		output := s.buffer.String()
		s.Contains(output, "error message")
		s.Contains(output, `"level":"error"`)
	})

	s.Run("Text format produces readable output", func() {
		logger := Initialize("info", "text")
		logger.SetOutput(s.buffer)

		entry := GetWithPackage("test_package")
		entry.Info("test message")

		output := s.buffer.String()
		s.Contains(output, "test message")
		s.Contains(output, "package=test_package")
		s.Contains(output, "level=info")
	})
}

func (s *LoggerTestSuite) TestJSONFormat() {
	s.Run("JSON format produces valid JSON", func() {
		logger := Initialize("info", "json")
		logger.SetOutput(s.buffer)

		entry := GetWithPackage("test_package")
		entry.WithField("custom_field", "custom_value").Info("json test message")

		output := strings.TrimSpace(s.buffer.String())

		// Verify it's valid JSON
		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(output), &jsonData)
		s.NoError(err)

		// Verify expected fields
		s.Equal("json test message", jsonData["msg"])
		s.Equal("info", jsonData["level"])
		s.Equal("test_package", jsonData["package"])
		s.Equal("custom_value", jsonData["custom_field"])
	})
}

func (s *LoggerTestSuite) TestConcurrency() {
	s.Run("Concurrent access is thread-safe", func() {
		const numGoroutines = 100
		done := make(chan bool, numGoroutines)

		// Use a discard writer to avoid verbose output
		Initialize("info", "json")
		Get().SetOutput(s.buffer)

		// Start multiple goroutines that access the logger
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				// Some goroutines initialize, others just get
				if id%2 == 0 {
					Initialize("debug", "json")
				} else {
					Get()
				}

				// All should be able to log without panicking
				entry := GetWithPackage("concurrent_test")
				entry.WithField("goroutine_id", id).Info("concurrent access test")
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Verify logger is still functional
		logger := Get()
		s.NotNil(logger)

		// Verify we got some log output (indicating no race conditions prevented logging)
		s.True(s.buffer.Len() > 0)
	})
}

func (s *LoggerTestSuite) TestEdgeCases() {
	s.Run("Empty package name", func() {
		entry := GetWithPackage("")
		s.NotNil(entry)
		s.Equal("", entry.Data["package"])
	})

	s.Run("Empty field value", func() {
		entry := GetWithField("empty", "")
		s.NotNil(entry)
		s.Equal("", entry.Data["empty"])
	})

	s.Run("Nil-like field value", func() {
		entry := GetWithField("nil_value", "")
		entry = entry.WithField("actual_nil", nil)
		s.NotNil(entry)
		s.Nil(entry.Data["actual_nil"])
	})

	s.Run("Special characters in package name", func() {
		specialPackage := "test/package-with_special.chars"
		entry := GetWithPackage(specialPackage)
		s.Equal(specialPackage, entry.Data["package"])
	})
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
