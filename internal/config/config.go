package config

import (
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type Config struct {
	Listen           string `mapstructure:"listen"`
	LogLevel         string `mapstructure:"log_level"`
	LogFormat        string `mapstructure:"log_format"`
	Initiator        string `mapstructure:"initiator"`
	Secret           string `mapstructure:"secret"`
	Database         string `mapstructure:"database"`
	TelegramBotToken string `mapstructure:"telegram_bot_token"`
	TelegramBotUser  tgbotapi.User
	MetricsListen    string `mapstructure:"metrics_listen"`
	UploadPath       string `mapstructure:"upload_path"`
	UploadURL        string `mapstructure:"upload_url"`
	FileBackend      afero.Fs
}

// NewConfig returns config with default values.
func NewConfig() Config {
	secret, _ := password.Generate(64, 12, 12, false, false)

	return Config{
		Listen:        ":8080",
		LogLevel:      "debug",
		LogFormat:     "json",
		Initiator:     "admin@systemli.org",
		Secret:        secret,
		Database:      "ticker.db",
		MetricsListen: ":8181",
		UploadPath:    "uploads",
		UploadURL:     "http://localhost:8080",
		FileBackend:   afero.NewOsFs(),
	}
}

// TelegramEnabled returns true if the required token is not empty.
func (c *Config) TelegramEnabled() bool {
	return c.TelegramBotToken != ""
}

// LoadConfig loads config from file.
func LoadConfig(path string) Config {
	c := NewConfig()
	viper.SetEnvPrefix("ticker")
	viper.AutomaticEnv()

	viper.SetDefault("listen", c.Listen)
	viper.SetDefault("log_level", c.LogLevel)
	viper.SetDefault("log_format", c.LogFormat)
	viper.SetDefault("initiator", c.Initiator)
	viper.SetDefault("secret", c.Secret)
	viper.SetDefault("database", c.Database)
	viper.SetDefault("metrics_listen", c.MetricsListen)
	viper.SetDefault("telegram_bot_token", "")
	viper.SetDefault("upload_path", c.UploadPath)
	viper.SetDefault("upload_url", c.UploadURL)

	//TODO: Make configurable
	fs := afero.NewOsFs()
	c.FileBackend = fs

	viper.SetFs(fs)

	if path != "" {
		dir, file := filepath.Split(path)
		// use current directory as default
		if dir == "" {
			dir = "."
		}
		// remove file name extensions
		file = strings.TrimSuffix(file, filepath.Ext(file))

		viper.SetConfigName(file)
		viper.AddConfigPath(dir)

		err := viper.ReadInConfig()
		if err != nil {
			log.WithError(err).Error("Falling back to ENV vars.")
		}
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		log.WithError(err).Panic("Unable to decode config into struct")
	}

	return c
}
