package realtime

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/suite"
)

type MetricsTestSuite struct {
	// Store original metrics to restore after tests
	originalConnectedClients  *prometheus.GaugeVec
	originalTotalConnections  *prometheus.CounterVec
	originalDisconnections    *prometheus.CounterVec
	originalMessagesSent      *prometheus.CounterVec
	originalMessagesDropped   *prometheus.CounterVec
	originalBroadcastDuration *prometheus.HistogramVec
	originalTotalClientsGauge prometheus.Gauge
	suite.Suite
}

func (s *MetricsTestSuite) SetupSuite() {
	// Save current metrics
	s.originalConnectedClients = connectedClients
	s.originalTotalConnections = totalConnections
	s.originalDisconnections = disconnections
	s.originalMessagesSent = messagesSent
	s.originalMessagesDropped = messagesDropped
	s.originalBroadcastDuration = broadcastDuration
	s.originalTotalClientsGauge = totalClientsGauge

	// Create new metrics for testing to avoid conflicts
	connectedClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "test_websocket_connected_clients",
			Help: "Test metric for connected clients",
		},
		[]string{"origin"},
	)

	totalConnections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_websocket_connections_total",
			Help: "Test metric for total connections",
		},
		[]string{"origin"},
	)

	disconnections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_websocket_disconnections_total",
			Help: "Test metric for disconnections",
		},
		[]string{"origin", "reason"},
	)

	messagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_websocket_messages_sent_total",
			Help: "Test metric for messages sent",
		},
		[]string{"origin", "message_type"},
	)

	messagesDropped = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_websocket_messages_dropped_total",
			Help: "Test metric for messages dropped",
		},
		[]string{"origin", "message_type"},
	)

	broadcastDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "test_websocket_broadcast_duration_seconds",
			Help:    "Test metric for broadcast duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"origin", "message_type"},
	)

	totalClientsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_websocket_total_connected_clients",
			Help: "Test metric for total clients",
		},
	)
}

func (s *MetricsTestSuite) TearDownSuite() {
	// Restore original metrics after tests
	connectedClients = s.originalConnectedClients
	totalConnections = s.originalTotalConnections
	disconnections = s.originalDisconnections
	messagesSent = s.originalMessagesSent
	messagesDropped = s.originalMessagesDropped
	broadcastDuration = s.originalBroadcastDuration
	totalClientsGauge = s.originalTotalClientsGauge
}

func (s *MetricsTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		// Reset metrics before each subtest
		totalClientsGauge.Set(0)
		subtest()
	})
}

func (s *MetricsTestSuite) TestRecordClientConnected() {
	s.Run("client connection increments metrics", func() {
		origin := "TestRecordClientConnected"

		recordClientConnected(origin)

		// Check connected clients gauge
		s.Equal(float64(1), testutil.ToFloat64(connectedClients.WithLabelValues(origin)))

		// Check total connections counter
		s.Equal(float64(1), testutil.ToFloat64(totalConnections.WithLabelValues(origin)))

		// Check total clients gauge
		s.Equal(float64(1), testutil.ToFloat64(totalClientsGauge))

		// Connect another client to same ticker
		recordClientConnected(origin)
		s.Equal(float64(2), testutil.ToFloat64(connectedClients.WithLabelValues(origin)))
		s.Equal(float64(2), testutil.ToFloat64(totalConnections.WithLabelValues(origin)))
		s.Equal(float64(2), testutil.ToFloat64(totalClientsGauge))
	})
}

func (s *MetricsTestSuite) TestRecordClientDisconnected() {
	s.Run("client disconnection decrements metrics", func() {
		origin := "TestRecordClientDisconnected"
		reason := "normal"

		// First connect a client
		recordClientConnected(origin)
		s.Equal(float64(1), testutil.ToFloat64(connectedClients.WithLabelValues(origin)))

		// Then disconnect
		recordClientDisconnected(origin, reason)

		// Check connected clients gauge decremented
		s.Equal(float64(0), testutil.ToFloat64(connectedClients.WithLabelValues(origin)))

		// Check disconnections counter incremented
		s.Equal(float64(1), testutil.ToFloat64(disconnections.WithLabelValues(origin, "normal")))

		// Check total clients gauge decremented
		s.Equal(float64(0), testutil.ToFloat64(totalClientsGauge))
	})
}

func (s *MetricsTestSuite) TestRecordMessageSent() {
	s.Run("message sent increments counter", func() {
		origin := "TestRecordMessageSent"
		messageType := "message_created"

		recordMessageSent(origin, messageType)
		recordMessageSent(origin, messageType)

		s.Equal(float64(2), testutil.ToFloat64(messagesSent.WithLabelValues(origin, "message_created")))
	})
}

func (s *MetricsTestSuite) TestRecordMessageDropped() {
	s.Run("message dropped increments counter", func() {
		origin := "TestRecordMessageDropped"
		messageType := "message_deleted"

		recordMessageDropped(origin, messageType)

		s.Equal(float64(1), testutil.ToFloat64(messagesDropped.WithLabelValues(origin, "message_deleted")))
	})
}

func (s *MetricsTestSuite) TestRecordBroadcastDuration() {
	s.Run("broadcast duration records histogram", func() {
		origin := "TestRecordBroadcastDuration"
		messageType := "message_created"
		duration := 50 * time.Millisecond

		recordBroadcastDuration(origin, messageType, duration)

		// Check that histogram was updated by getting a metric sample
		metric := &dto.Metric{}
		histogram := broadcastDuration.WithLabelValues(origin, "message_created")
		err := histogram.(prometheus.Histogram).Write(metric)
		s.NoError(err)
		s.Equal(uint64(1), metric.GetHistogram().GetSampleCount())
		s.InDelta(0.05, metric.GetHistogram().GetSampleSum(), 0.001) // 50ms = 0.05s
	})
}

func (s *MetricsTestSuite) TestMultipleTickers() {
	s.Run("metrics work with multiple tickers", func() {
		// Test metrics with multiple tickers
		recordClientConnected("1")
		recordClientConnected("2")
		recordClientConnected("1") // Second client for ticker 1

		s.Equal(float64(2), testutil.ToFloat64(connectedClients.WithLabelValues("1")))
		s.Equal(float64(1), testutil.ToFloat64(connectedClients.WithLabelValues("2")))
		s.Equal(float64(3), testutil.ToFloat64(totalClientsGauge))

		// Disconnect one client from ticker 1
		recordClientDisconnected("1", "normal")

		s.Equal(float64(1), testutil.ToFloat64(connectedClients.WithLabelValues("1")))
		s.Equal(float64(1), testutil.ToFloat64(connectedClients.WithLabelValues("2")))
		s.Equal(float64(2), testutil.ToFloat64(totalClientsGauge))
	})
}

func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}
