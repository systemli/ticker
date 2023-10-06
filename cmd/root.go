package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/config"
)

var (
	GitCommit  string
	GitVersion string

	configPath string
	cfg        config.Config

	rootCmd = &cobra.Command{
		Use:   "ticker",
		Short: "Service to distribute short messages",
		Long:  "Service to distribute short messages in support of events, demonstrations, or other time-sensitive events.",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yml", "path to config.yml")
}

func initConfig() {
	cfg = config.LoadConfig(configPath)
	//TODO: Improve startup routine
	if cfg.TelegramEnabled() {
		user, err := bridge.BotUser(cfg.TelegramBotToken)
		if err != nil {
			log.WithError(err).Error("Unable to retrieve the user information for the Telegram Bot")
		} else {
			cfg.TelegramBotUser = user
		}
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	if cfg.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	}
	log.SetLevel(lvl)
}

func Execute() {
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
