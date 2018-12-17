package model

import (
	"fmt"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
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
}

//NewConfig returns config with default values.
func NewConfig() *config {
	secret, _ := password.Generate(64, 12, 12, false, false)

	return &config{
		Listen:    ":8080",
		LogLevel:  "debug",
		Initiator: "admin@systemli.org",
		Secret:    secret,
		Database:  "ticker.db",
	}
}

//TwitterEnabled returns true if required keys not empty.
func (c *config) TwitterEnabled() bool {
	return c.TwitterConsumerKey != "" && c.TwitterConsumerSecret != ""
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
	viper.SetDefault("twitter_consumer_key", "")
	viper.SetDefault("twitter_consumer_secret", "")

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
		fmt.Printf("%s \nFalling back to env vars.", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	return Config
}
