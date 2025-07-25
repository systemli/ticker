package realtime

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// WebSocket Prometheus metrics
var (
	// connectedClients tracks the current number of connected WebSocket clients per ticker
	connectedClients = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "websocket_connected_clients",
			Help: "Current number of connected WebSocket clients per ticker",
		},
		[]string{"ticker_id"},
	)

	// totalConnections tracks the total number of WebSocket connections established
	totalConnections = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_connections_total",
			Help: "Total number of WebSocket connections established",
		},
		[]string{"ticker_id"},
	)

	// disconnections tracks the total number of WebSocket disconnections
	disconnections = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_disconnections_total",
			Help: "Total number of WebSocket disconnections",
		},
		[]string{"ticker_id", "reason"},
	)

	// messagesSent tracks the total number of messages sent to WebSocket clients
	messagesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_total",
			Help: "Total number of messages sent to WebSocket clients",
		},
		[]string{"ticker_id", "message_type"},
	)

	// messagesDropped tracks messages that couldn't be delivered to clients
	messagesDropped = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_dropped_total",
			Help: "Total number of messages dropped due to client unavailability",
		},
		[]string{"ticker_id", "message_type"},
	)

	// broadcastDuration tracks how long it takes to broadcast messages to all clients
	broadcastDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "websocket_broadcast_duration_seconds",
			Help:    "Time spent broadcasting messages to WebSocket clients",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"ticker_id", "message_type"},
	)

	// totalClientsGauge tracks the total number of connected clients across all tickers
	totalClientsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_total_connected_clients",
			Help: "Total number of connected WebSocket clients across all tickers",
		},
	)
)

// recordClientConnected increments connection metrics when a client connects
func recordClientConnected(tickerID int) {
	tickerIDStr := strconv.Itoa(tickerID)
	connectedClients.WithLabelValues(tickerIDStr).Inc()
	totalConnections.WithLabelValues(tickerIDStr).Inc()
	totalClientsGauge.Inc()
}

// recordClientDisconnected decrements connection metrics when a client disconnects
func recordClientDisconnected(tickerID int, reason string) {
	tickerIDStr := strconv.Itoa(tickerID)
	connectedClients.WithLabelValues(tickerIDStr).Dec()
	disconnections.WithLabelValues(tickerIDStr, reason).Inc()
	totalClientsGauge.Dec()
}

// recordMessageSent increments the messages sent counter
func recordMessageSent(tickerID int, messageType string) {
	tickerIDStr := strconv.Itoa(tickerID)
	messagesSent.WithLabelValues(tickerIDStr, messageType).Inc()
}

// recordMessageDropped increments the messages dropped counter
func recordMessageDropped(tickerID int, messageType string) {
	tickerIDStr := strconv.Itoa(tickerID)
	messagesDropped.WithLabelValues(tickerIDStr, messageType).Inc()
}

// recordBroadcastDuration records the time taken to broadcast a message
func recordBroadcastDuration(tickerID int, messageType string, duration time.Duration) {
	tickerIDStr := strconv.Itoa(tickerID)
	broadcastDuration.WithLabelValues(tickerIDStr, messageType).Observe(duration.Seconds())
}
