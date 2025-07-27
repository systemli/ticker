package realtime

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// connectedClients tracks the current number of connected WebSocket clients per ticker
	connectedClients = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "websocket_connected_clients",
			Help: "Current number of connected WebSocket clients per ticker",
		},
		[]string{"origin"},
	)

	// totalConnections tracks the total number of WebSocket connections established
	totalConnections = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_connections_total",
			Help: "Total number of WebSocket connections established",
		},
		[]string{"origin"},
	)

	// disconnections tracks the total number of WebSocket disconnections
	disconnections = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_disconnections_total",
			Help: "Total number of WebSocket disconnections",
		},
		[]string{"origin", "reason"},
	)

	// messagesSent tracks the total number of messages sent to WebSocket clients
	messagesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_sent_total",
			Help: "Total number of messages sent to WebSocket clients",
		},
		[]string{"origin", "message_type"},
	)

	// messagesDropped tracks messages that couldn't be delivered to clients
	messagesDropped = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_dropped_total",
			Help: "Total number of messages dropped due to client unavailability",
		},
		[]string{"origin", "message_type"},
	)

	// broadcastDuration tracks how long it takes to broadcast messages to all clients
	broadcastDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "websocket_broadcast_duration_seconds",
			Help:    "Time spent broadcasting messages to WebSocket clients",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"origin", "message_type"},
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
func recordClientConnected(origin string) {
	connectedClients.WithLabelValues(origin).Inc()
	totalConnections.WithLabelValues(origin).Inc()
	totalClientsGauge.Inc()
}

// recordClientDisconnected decrements connection metrics when a client disconnects
func recordClientDisconnected(origin string, reason string) {
	connectedClients.WithLabelValues(origin).Dec()
	disconnections.WithLabelValues(origin, reason).Inc()
	totalClientsGauge.Dec()
}

// recordMessageSent increments the messages sent counter
func recordMessageSent(origin string, messageType string) {
	messagesSent.WithLabelValues(origin, messageType).Inc()
}

// recordMessageDropped increments the messages dropped counter
func recordMessageDropped(origin string, messageType string) {
	messagesDropped.WithLabelValues(origin, messageType).Inc()
}

// recordBroadcastDuration records the time taken to broadcast a message
func recordBroadcastDuration(origin string, messageType string, duration time.Duration) {
	broadcastDuration.WithLabelValues(origin, messageType).Observe(duration.Seconds())
}
