package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the ticker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("starting ticker api on %s", cfg.Listen)
			if GitCommit != "" && GitVersion != "" {
				log.Infof("build info: %s (commit: %s)", GitVersion, GitCommit)
			}

			go func() {
				http.Handle("/metrics", promhttp.Handler())
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(`<html>
					<head><title>Ticker Metrics Exporter</title></head>
					<body>
					<h1>Ticker Metrics Exporter</h1>
					<p><a href="/metrics">Metrics</a></p>
					</body>
					</html>`))
				})
				log.Fatal(http.ListenAndServe(cfg.MetricsListen, nil))
			}()

			var dialector gorm.Dialector
			switch cfg.Database.Type {
			case "sqlite":
				dialector = sqlite.Open(cfg.Database.DSN)
			case "mysql":
				dialector = mysql.Open(cfg.Database.DSN)
			case "postgres":
				dialector = postgres.Open(cfg.Database.DSN)
			default:
				log.Fatalf("unknown database type %s", cfg.Database.Type)
			}

			db, err := gorm.Open(dialector, &gorm.Config{
				Logger: logger.NewGormLogger(log),
			})
			if err != nil {
				log.WithError(err).Fatal("could not connect to database")
			}
			store := storage.NewSqlStorage(db, cfg.UploadPath)
			err = db.AutoMigrate(
				&storage.Attachment{},
				&storage.Message{},
				&storage.Setting{},
				&storage.Ticker{},
				&storage.TickerInformation{},
				&storage.TickerMastodon{},
				&storage.TickerTelegram{},
				&storage.Upload{},
				&storage.User{},
			)
			if err != nil {
				log.WithError(err).Fatal("could not migrate database")
			}

			router := api.API(cfg, store, log)
			server := &http.Server{
				Addr:    cfg.Listen,
				Handler: router,
			}

			firstRun(store, cfg)

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}()

			// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit

			log.Infoln("Shutdown Ticker")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func firstRun(store storage.Storage, config config.Config) {
	count, err := store.CountUser()
	if err != nil {
		log.Fatal("error using database")
	}

	if count == 0 {
		pw, err := password.Generate(24, 3, 3, false, false)
		if err != nil {
			log.Fatal(err)
		}

		user, err := storage.NewUser(config.Initiator, pw)
		user.IsSuperAdmin = true
		if err != nil {
			log.Fatal("could not create first user")
		}

		err = store.SaveUser(&user)
		if err != nil {
			log.Fatal("could not persist first user")
		}

		log.WithField("email", user.Email).WithField("password", pw).Info("admin user created (change password now!)")
	}
}
