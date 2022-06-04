package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/systemli/ticker/internal/bridge"

	"github.com/sethvargo/go-password/password"

	log "github.com/sirupsen/logrus"

	. "github.com/systemli/ticker/internal/api"
	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

var (
	GitCommit  string
	GitVersion string
)

func main() {
	router := API()
	server := &http.Server{
		Addr:    Config.Listen,
		Handler: router,
	}

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

func init() {
	var cp = flag.String("config", "config.yml", "path to config.yml")
	flag.Parse()

	Config = LoadConfig(*cp)
	//TODO: Improve startup routine
	if Config.TelegramEnabled() {
		user, err := bridge.BotUser(Config.TelegramBotToken)
		if err != nil {
			log.WithError(err).Error("Unable to retrieve the user information for the Telegram Bot")
		} else {
			Config.TelegramBotUser = user
		}
	}
	DB = OpenDB(Config.Database)

	firstRun()

	log.Println("Starting Ticker API")
	log.Printf("Listen on %s", Config.Listen)

	buildInfo()

	lvl, err := log.ParseLevel(Config.LogLevel)
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
		log.Fatal(http.ListenAndServe(Config.MetricsListen, nil))
	}()
}

func buildInfo() {
	if GitCommit != "" && GitVersion != "" {
		log.Println("Build Information")
		log.Printf("Version: %s", GitVersion)
		log.Printf("Commit: %s", GitCommit)
	}
}

func firstRun() {
	count, err := DB.Count(&User{})
	if err != nil {
		log.Fatal("error using database")
	}

	if count == 0 {
		pw, err := password.Generate(24, 3, 3, false, false)
		if err != nil {
			log.Fatal(err)
		}

		user, err := NewAdminUser(Config.Initiator, pw)
		if err != nil {
			log.Fatal("could not create first user")
		}

		err = DB.Save(user)
		if err != nil {
			log.Fatal("could not persist first user")
		}

		log.WithField("email", user.Email).WithField("password", pw).Info("admin user created (change password now!)")
	}
}
