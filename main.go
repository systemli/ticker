package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"

	"github.com/sethvargo/go-password/password"

	log "github.com/sirupsen/logrus"
)

var (
	GitCommit  string
	GitVersion string
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yml", "path to config.yml")
	flag.Parse()

	config := config.LoadConfig(configPath)
	//TODO: Improve startup routine
	if config.TelegramEnabled() {
		user, err := bridge.BotUser(config.TelegramBotToken)
		if err != nil {
			log.WithError(err).Error("Unable to retrieve the user information for the Telegram Bot")
		} else {
			config.TelegramBotUser = user
		}
	}

	log.Println("Starting Ticker API")
	log.Printf("Listen on %s", config.Listen)

	buildInfo()

	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)

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
		log.Fatal(http.ListenAndServe(config.MetricsListen, nil))
	}()

	store := storage.NewStorage(config.Database, config.UploadPath)
	router := api.API(config, store)
	server := &http.Server{
		Addr:    config.Listen,
		Handler: router,
	}

	firstRun(store, config)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Ticker")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func buildInfo() {
	if GitCommit != "" && GitVersion != "" {
		log.Println("Build Information")
		log.Printf("Version: %s", GitVersion)
		log.Printf("Commit: %s", GitCommit)
	}
}

func firstRun(store storage.TickerStorage, config config.Config) {
	count, err := store.CountUser()
	if err != nil {
		log.Fatal("error using database")
	}

	if count == 0 {
		pw, err := password.Generate(24, 3, 3, false, false)
		if err != nil {
			log.Fatal(err)
		}

		user, err := storage.NewAdminUser(config.Initiator, pw)
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
