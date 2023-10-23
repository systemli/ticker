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

	log = logger.NewLogrus(cfg.LogLevel, cfg.LogFormat)

	var err error
	db, err = storage.OpenGormDB(cfg.Database.Type, cfg.Database.DSN, log)
	if err != nil {
		log.WithError(err).Fatal("could not connect to database")
	}
	store = storage.NewSqlStorage(db, cfg.UploadPath)
	if err := storage.MigrateDB(db); err != nil {
		log.WithError(err).Fatal("could not migrate database")
	}
}

func Execute() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
