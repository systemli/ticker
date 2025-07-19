package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
	"gorm.io/gorm"
)

var (
	configPath string
	cfg        config.Config
	db         *gorm.DB
	store      *storage.SqlStorage

	log = logrus.New()

	rootCmd = &cobra.Command{
		Use:   "ticker",
		Short: "Service to distribute short messages",
		Long:  "Service to distribute short messages in support of events, demonstrations, or other time-sensitive events.",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to config.yml")
}

func initConfig() {
	cfg = config.LoadConfig(configPath)

	// Initialize the global logger with configuration
	log = logger.Initialize(cfg.LogLevel, cfg.LogFormat)

	//TODO: Improve startup routine
	if cfg.Telegram.Enabled() {
		user, err := bridge.BotUser(cfg.Telegram.Token)
		if err != nil {
			log.WithError(err).Error("Unable to retrieve the user information for the Telegram Bot")
		} else {
			cfg.Telegram.User = user
		}
	}

	var err error
	db, err = storage.OpenGormDB(cfg.Database.Type, cfg.Database.DSN, log)
	if err != nil {
		log.WithError(err).Fatal("could not connect to database")
	}
	store = storage.NewSqlStorage(db, cfg.Upload.Path)
	if err := storage.MigrateDB(db); err != nil {
		log.WithError(err).Fatal("could not migrate database")
	}
}

func Execute() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
