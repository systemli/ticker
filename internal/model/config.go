package model

import (
	"path/filepath"
	"strings"

	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var Config *config

type config struct {
	Listen                string `mapstructure:"listen"`
	LogLevel              string `mapstructure:"log_level"`
	Initiator             string `mapstructure:"initiator"`
	Secret                string `mapstructure:"secret"`
	Database              string `mapstructure:"database"`
	TwitterConsumerKey    string `mapstructure:"twitter_consumer_key"`
	TwitterConsumerSecret string `mapstructure:"twitter_consumer_secret"`
	TelegramBotToken      string `mapstructure:"telegram_bot_token"`
	MetricsListen         string `mapstructure:"metrics_listen"`
	UploadPath            string `mapstructure:"upload_path"`
	UploadURL             string `mapstructure:"upload_url"`
	FileBackend           afero.Fs
}

//NewConfig returns config with default values.
func NewConfig() *config {
	secret, _ := password.Generate(64, 12, 12, false, false)

	return &config{
		Listen:        ":8080",
		LogLevel:      "debug",
		Initiator:     "admin@systemli.org",
		Secret:        secret,
		Database:      "ticker.db",
		MetricsListen: ":8181",
		UploadPath:    "uploads",
		UploadURL:     "http://localhost:8080",
		FileBackend:   afero.NewOsFs(),
	}
}

//TwitterEnabled returns true if required keys not empty.
func (c *config) TwitterEnabled() bool {
	return c.TwitterConsumerKey != "" && c.TwitterConsumerSecret != ""
}

//TelegramBotEnabled returns true if the required token is not empty.
func (c *config) TelegramBotEnabled() bool {
	return c.TelegramBotToken != ""
}

//LoadConfig loads config from file.
func LoadConfig(path string) *config {
	c := NewConfig()
	viper.SetEnvPrefix("ticker")
	viper.AutomaticEnv()

	viper.SetDefault("listen", c.Listen)
	viper.SetDefault("log_level", c.LogLevel)
	viper.SetDefault("initiator", c.Initiator)
	viper.SetDefault("secret", c.Secret)
	viper.SetDefault("database", c.Database)
	viper.SetDefault("metrics_listen", c.MetricsListen)
	viper.SetDefault("twitter_consumer_key", "")
	viper.SetDefault("twitter_consumer_secret", "")
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

	Config = c
	return Config
}
