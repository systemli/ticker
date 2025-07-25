package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/systemli/ticker/internal/api"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the ticker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("starting ticker (version: %s, commit: %s) on %s", version, commit, cfg.Listen)

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

			apiServer := api.API(cfg, store)
			server := &http.Server{
				Addr:    cfg.Listen,
				Handler: apiServer.Router,
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

			log.Infoln("Shutdown Ticker")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Shutdown realtime engine first
			if err := apiServer.Realtime.Shutdown(ctx); err != nil {
				log.WithError(err).Warn("Realtime engine shutdown failed")
			}

			// Shutdown HTTP server
			if err := server.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		},
	}
)
