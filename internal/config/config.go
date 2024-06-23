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
	Listen        string      `yaml:"listen"`
	LogLevel      string      `yaml:"log_level"`
	LogFormat     string      `yaml:"log_format"`
	Secret        string      `yaml:"secret"`
	Database      Database    `yaml:"database"`
	Telegram      Telegram    `yaml:"telegram"`
	SignalGroup   SignalGroup `yaml:"signal_group"`
	MetricsListen string      `yaml:"metrics_listen"`
	Upload        Upload      `yaml:"upload"`
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

type SignalGroup struct {
	ApiUrl  string `yaml:"api_url"`
	Avatar  string `yaml:"avatar"`
	Account string `yaml:"account"`
}

type Upload struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func defaultConfig() Config {
	secret, _ := password.Generate(64, 12, 12, false, true)

	return Config{
		Listen:        ":8080",
		LogLevel:      "debug",
		LogFormat:     "json",
		Secret:        secret,
		Database:      Database{Type: "sqlite", DSN: "ticker.db"},
		MetricsListen: ":8181",
		Upload: Upload{
			Path: "uploads",
			URL:  "http://localhost:8080",
		},
		FileBackend: afero.NewOsFs(),
	}
}

// Enabled returns true if the required token is not empty.
func (t *Telegram) Enabled() bool {
	return t.Token != ""
}

// Enabled returns true if requried API URL and account are set.
func (t *SignalGroup) Enabled() bool {
	return t.ApiUrl != "" && t.Account != ""
}

// LoadConfig loads config from file.
func LoadConfig(path string) Config {
	c := defaultConfig()
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
		c.Upload.Path = os.Getenv("TICKER_UPLOAD_PATH")
	}
	if os.Getenv("TICKER_UPLOAD_URL") != "" {
		c.Upload.URL = os.Getenv("TICKER_UPLOAD_URL")
	}
	if os.Getenv("TICKER_TELEGRAM_TOKEN") != "" {
		c.Telegram.Token = os.Getenv("TICKER_TELEGRAM_TOKEN")
	}
	if os.Getenv("TICKER_SIGNAL_GROUP_API_URL") != "" {
		c.SignalGroup.ApiUrl = os.Getenv("TICKER_SIGNAL_GROUP_API_URL")
	}
	if os.Getenv("TICKER_SIGNAL_GROUP_ACCOUNT") != "" {
		c.SignalGroup.ApiUrl = os.Getenv("TICKER_SIGNAL_GROUP_ACCOUNT")
	}
	if os.Getenv("TICKER_SIGNAL_GROUP_AVATAR") != "" {
		c.SignalGroup.ApiUrl = os.Getenv("TICKER_SIGNAL_GROUP_AVATAR")
	}

	return c
}
