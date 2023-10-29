package config

import (
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

var log = logrus.WithField("package", "config")

type Config struct {
	Listen        string   `yaml:"listen"`
	LogLevel      string   `yaml:"log_level"`
	LogFormat     string   `yaml:"log_format"`
	Secret        string   `yaml:"secret"`
	Database      Database `yaml:"database"`
	Telegram      Telegram `yaml:"telegram"`
	MetricsListen string   `yaml:"metrics_listen"`
	UploadPath    string   `yaml:"upload_path"`
	UploadURL     string   `yaml:"upload_url"`
	FileBackend   afero.Fs
}

type Database struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

type Telegram struct {
	Token string `yaml:"token"`
	User  tgbotapi.User
}

// NewConfig returns config with default values.
func NewConfig() Config {
	secret, _ := password.Generate(64, 12, 12, false, true)

	return Config{
		Listen:        ":8080",
		LogLevel:      "debug",
		LogFormat:     "json",
		Secret:        secret,
		Database:      Database{Type: "sqlite", DSN: "ticker.db"},
		MetricsListen: ":8181",
		UploadPath:    "uploads",
		UploadURL:     "http://localhost:8080",
		FileBackend:   afero.NewOsFs(),
	}
}

// Enabled returns true if the required token is not empty.
func (t *Telegram) Enabled() bool {
	return t.Token != ""
}

// LoadConfig loads config from file.
func LoadConfig(path string) Config {
	c := NewConfig()
	c.FileBackend = afero.NewOsFs()

	if path != "" {
		bytes, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			log.WithError(err).Error("Unable to load config")
		}
		if err := yaml.Unmarshal(bytes, &c); err != nil {
			log.WithError(err).Error("Unable to load config")
		}
	}

	if os.Getenv("TICKER_LISTEN") != "" {
		c.Listen = os.Getenv("TICKER_LISTEN")
	}
	if os.Getenv("TICKER_LOG_LEVEL") != "" {
		c.LogLevel = os.Getenv("TICKER_LOG_LEVEL")
	}
	if os.Getenv("TICKER_LOG_FORMAT") != "" {
		c.LogFormat = os.Getenv("TICKER_LOG_FORMAT")
	}
	if os.Getenv("TICKER_SECRET") != "" {
		c.Secret = os.Getenv("TICKER_SECRET")
	}
	if os.Getenv("TICKER_DATABASE_TYPE") != "" {
		c.Database.Type = os.Getenv("TICKER_DATABASE_TYPE")
	}
	if os.Getenv("TICKER_DATABASE_DSN") != "" {
		c.Database.DSN = os.Getenv("TICKER_DATABASE_DSN")
	}
	if os.Getenv("TICKER_METRICS_LISTEN") != "" {
		c.MetricsListen = os.Getenv("TICKER_METRICS_LISTEN")
	}
	if os.Getenv("TICKER_UPLOAD_PATH") != "" {
		c.UploadPath = os.Getenv("TICKER_UPLOAD_PATH")
	}
	if os.Getenv("TICKER_UPLOAD_URL") != "" {
		c.UploadURL = os.Getenv("TICKER_UPLOAD_URL")
	}
	if os.Getenv("TICKER_TELEGRAM_TOKEN") != "" {
		c.Telegram.Token = os.Getenv("TICKER_TELEGRAM_TOKEN")
	}

	return c
}
