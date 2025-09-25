package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	envs map[string]string
	suite.Suite
}

func (s *ConfigTestSuite) SetupTest() {
	log.Logger.SetOutput(io.Discard)

	s.envs = map[string]string{
		"TICKER_LISTEN":         ":7070",
		"TICKER_LOG_LEVEL":      "trace",
		"TICKER_LOG_FORMAT":     "text",
		"TICKER_SECRET":         "secret",
		"TICKER_DATABASE_TYPE":  "mysql",
		"TICKER_DATABASE_DSN":   "user:password@tcp(localhost:3306)/ticker?charset=utf8mb4&parseTime=True&loc=Local",
		"TICKER_METRICS_LISTEN": ":9191",
		"TICKER_UPLOAD_PATH":    "/data/uploads",
		"TICKER_UPLOAD_URL":     "https://example.com",
	}
}

func (s *ConfigTestSuite) TestConfig() {
	s.Run("LoadConfig", func() {
		s.Run("when path is empty", func() {
			s.Run("loads config with default values", func() {
				c := LoadConfig("")
				s.Equal(":8080", c.Listen)
				s.Equal("debug", c.LogLevel)
				s.Equal("json", c.LogFormat)
				s.NotEmpty(c.Secret)
				s.Equal("sqlite", c.Database.Type)
				s.Equal("ticker.db", c.Database.DSN)
				s.Equal(":8181", c.MetricsListen)
				s.Equal("uploads", c.Upload.Path)
				s.Equal("http://localhost:8080", c.Upload.URL)
				s.Empty(c.SignalGroup.ApiUrl)
				s.Empty(c.SignalGroup.Account)
				s.False(c.SignalGroup.Enabled())
			})

			s.Run("loads config from env", func() {
				for key, value := range s.envs {
					err := os.Setenv(key, value)
					s.NoError(err)
				}

				c := LoadConfig("")
				s.Equal(s.envs["TICKER_LISTEN"], c.Listen)
				s.Equal(s.envs["TICKER_LOG_LEVEL"], c.LogLevel)
				s.Equal(s.envs["TICKER_LOG_FORMAT"], c.LogFormat)
				s.Equal(s.envs["TICKER_SECRET"], c.Secret)
				s.Equal(s.envs["TICKER_DATABASE_TYPE"], c.Database.Type)
				s.Equal(s.envs["TICKER_DATABASE_DSN"], c.Database.DSN)
				s.Equal(s.envs["TICKER_METRICS_LISTEN"], c.MetricsListen)
				s.Equal(s.envs["TICKER_UPLOAD_PATH"], c.Upload.Path)
				s.Equal(s.envs["TICKER_UPLOAD_URL"], c.Upload.URL)
				s.Equal(s.envs["TICKER_SIGNAL_GROUP_API_URL"], c.SignalGroup.ApiUrl)
				s.Equal(s.envs["TICKER_SIGNAL_GROUP_ACCOUNT"], c.SignalGroup.Account)

				for key := range s.envs {
					os.Unsetenv(key)
				}
			})
		})

		s.Run("when path is not empty", func() {
			s.Run("when path is absolute", func() {
				s.Run("loads config from file", func() {
					path, err := filepath.Abs("../../testdata/config_valid.yml")
					s.NoError(err)
					c := LoadConfig(path)
					s.Equal("127.0.0.1:8888", c.Listen)
				})
			})

			s.Run("when path is relative", func() {
				s.Run("loads config from file", func() {
					c := LoadConfig("../../testdata/config_valid.yml")
					s.Equal("127.0.0.1:8888", c.Listen)
				})
			})

			s.Run("when file does not exist", func() {
				s.Run("loads config with default values", func() {
					c := LoadConfig("config_notfound.yml")
					s.Equal(":8080", c.Listen)
				})
			})

			s.Run("when file is invalid", func() {
				s.Run("loads config with default values", func() {
					c := LoadConfig("../../testdata/config_invalid.txt")
					s.Equal(":8080", c.Listen)
				})
			})
		})
	})
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
