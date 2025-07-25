package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/api/realtime"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type WebSocketTestSuite struct {
	w        *httptest.ResponseRecorder
	ctx      *gin.Context
	store    *storage.MockStorage
	cfg      config.Config
	realtime *realtime.Engine
	suite.Suite
}

func (s *WebSocketTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *WebSocketTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")
		s.realtime = realtime.New()

		subtest()
	})
}

// setupWebSocketServer creates a test server with WebSocket handler for the given ticker
func (s *WebSocketTestSuite) setupWebSocketServer(ticker storage.Ticker) (*httptest.Server, *realtime.Engine) {
	// Create a new realtime engine for this test to avoid channel conflicts
	realtimeEngine := realtime.New()

	// Start the realtime engine
	go realtimeEngine.Run()

	// Create a real HTTP server with the WebSocket handler
	router := gin.New()
	h := handler{
		storage:  s.store,
		config:   s.cfg,
		realtime: realtimeEngine, // Use the test-specific engine
	}

	// Set up the route with middleware that sets the ticker
	router.GET("/v1/ws", func(c *gin.Context) {
		c.Set("ticker", ticker)
		h.HandleWebSocket(c)
	})

	return httptest.NewServer(router), realtimeEngine
}

// createWebSocketClient creates a WebSocket client connected to the test server
func (s *WebSocketTestSuite) createWebSocketClient(server *httptest.Server, origin string) (*websocket.Conn, *http.Response, error) {
	wsURL := "ws" + server.URL[4:] + "/v1/ws" // Replace "http" with "ws"

	headers := http.Header{}
	headers.Set("Origin", origin)

	return websocket.DefaultDialer.Dial(wsURL, headers)
}

func (s *WebSocketTestSuite) TestHandleWebSocket() {
	s.Run("when ticker is missing", func() {
		// Don't set ticker in context, should return 404
		h := s.handler()
		h.HandleWebSocket(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.Contains(s.w.Body.String(), "ticker not found")
	})

	s.Run("when ticker is present but upgrade fails", func() {
		ticker := storage.Ticker{ID: 1, Websites: []storage.TickerWebsite{{ID: 1, Origin: "https://example.org"}}}
		s.ctx.Set("ticker", ticker)

		req := httptest.NewRequest(http.MethodGet, "/v1/ws", nil)
		req.Header.Set("Origin", "https://example.org")

		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/ws", nil)

		h := s.handler()
		h.HandleWebSocket(s.ctx)

		s.NotEqual(http.StatusOK, s.w.Code)
	})

	s.Run("single client receives message", func() {
		ticker := storage.Ticker{ID: 1, Websites: []storage.TickerWebsite{{ID: 1, Origin: "https://example.org"}}}

		server, realtimeEngine := s.setupWebSocketServer(ticker)
		defer server.Close()
		defer realtimeEngine.Shutdown(s.ctx)

		// Connect to the WebSocket server
		client, response, err := s.createWebSocketClient(server, "https://example.org")
		s.NoError(err, "Should successfully connect to WebSocket")
		s.Equal(http.StatusSwitchingProtocols, response.StatusCode)
		defer client.Close()

		// Test sending a message through the realtime engine
		testMessage := realtime.Message{
			Type:     "message_created",
			TickerID: ticker.ID,
			Data:     map[string]any{"id": 42, "text": "Test message"},
		}

		// Give the connection time to register
		time.Sleep(100 * time.Millisecond)

		// Broadcast the message
		realtimeEngine.Broadcast(testMessage)

		// Read the message from the WebSocket
		var receivedMessage realtime.Message
		err = client.ReadJSON(&receivedMessage)
		s.NoError(err, "Should receive message from WebSocket")

		// Verify the message content
		s.Equal("message_created", receivedMessage.Type)
		s.Equal(ticker.ID, receivedMessage.TickerID)

		data, ok := receivedMessage.Data.(map[string]interface{})
		s.True(ok, "Data should be a map")
		s.Equal(float64(42), data["id"]) // JSON numbers are float64
		s.Equal("Test message", data["text"])
	})

	s.Run("multiple clients receive broadcasts", func() {
		ticker := storage.Ticker{ID: 2, Websites: []storage.TickerWebsite{{ID: 1, Origin: "https://test.org"}}}

		server, realtimeEngine := s.setupWebSocketServer(ticker)
		defer server.Close()
		defer realtimeEngine.Shutdown(s.ctx)

		// Connect two clients
		client1, _, err := s.createWebSocketClient(server, "https://test.org")
		s.NoError(err)
		defer client1.Close()

		client2, _, err := s.createWebSocketClient(server, "https://test.org")
		s.NoError(err)
		defer client2.Close()

		// Give connections time to register
		time.Sleep(200 * time.Millisecond)

		// Broadcast a message
		testMessage := realtime.Message{
			Type:     "message_deleted",
			TickerID: ticker.ID,
			Data:     map[string]any{"id": 123},
		}

		realtimeEngine.Broadcast(testMessage)

		// Both clients should receive the message
		var msg1, msg2 realtime.Message

		err = client1.ReadJSON(&msg1)
		s.NoError(err, "Client 1 should receive message")
		s.Equal("message_deleted", msg1.Type)
		s.Equal(ticker.ID, msg1.TickerID)

		err = client2.ReadJSON(&msg2)
		s.NoError(err, "Client 2 should receive message")
		s.Equal("message_deleted", msg2.Type)
		s.Equal(ticker.ID, msg2.TickerID)
	})

	s.Run("clients only receive messages for their ticker", func() {
		ticker1 := storage.Ticker{ID: 10, Websites: []storage.TickerWebsite{{ID: 1, Origin: "https://ticker1.org"}}}
		ticker2 := storage.Ticker{ID: 20, Websites: []storage.TickerWebsite{{ID: 1, Origin: "https://ticker2.org"}}}

		// Set up two separate servers for different tickers
		server1, realtimeEngine1 := s.setupWebSocketServer(ticker1)
		defer server1.Close()
		defer realtimeEngine1.Shutdown(s.ctx)

		server2, realtimeEngine2 := s.setupWebSocketServer(ticker2)
		defer server2.Close()
		defer realtimeEngine2.Shutdown(s.ctx)

		// Connect clients to different ticker servers
		client1, _, err := s.createWebSocketClient(server1, "https://ticker1.org")
		s.NoError(err)
		defer client1.Close()

		client2, _, err := s.createWebSocketClient(server2, "https://ticker2.org")
		s.NoError(err)
		defer client2.Close()

		// Give connections time to register
		time.Sleep(200 * time.Millisecond)

		// Send message only to ticker1
		testMessage := realtime.Message{
			Type:     "message_created",
			TickerID: ticker1.ID,
			Data:     map[string]any{"content": "Message for ticker 1"},
		}

		realtimeEngine1.Broadcast(testMessage)

		// Client1 should receive the message
		var msg1 realtime.Message
		err = client1.ReadJSON(&msg1)
		s.NoError(err, "Client 1 should receive message for ticker 1")
		s.Equal(ticker1.ID, msg1.TickerID)

		// Client2 should NOT receive the message (timeout expected)
		client2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		var msg2 realtime.Message
		err = client2.ReadJSON(&msg2)
		s.Error(err, "Client 2 should not receive message for ticker 1")
	})
}

func (s *WebSocketTestSuite) TestUpgraderConfiguration() {
	s.Run("CheckOrigin function works correctly", func() {
		testCases := []struct {
			name     string
			setupReq func() *http.Request
			expected bool
		}{
			{
				name: "allows with Origin header",
				setupReq: func() *http.Request {
					req := httptest.NewRequest("GET", "/ws", nil)
					req.Header.Set("Origin", "https://example.com")
					return req
				},
				expected: true,
			},
			{
				name: "allows with origin query parameter",
				setupReq: func() *http.Request {
					req := httptest.NewRequest("GET", "/ws?origin=test", nil)
					return req
				},
				expected: true,
			},
			{
				name: "denies without Origin header or query parameter",
				setupReq: func() *http.Request {
					req := httptest.NewRequest("GET", "/ws", nil)
					return req
				},
				expected: false,
			},
			{
				name: "allows with empty Origin header but has origin query param",
				setupReq: func() *http.Request {
					req := httptest.NewRequest("GET", "/ws?origin=", nil)
					req.Header.Set("Origin", "")
					return req
				},
				expected: true,
			},
		}

		for _, tc := range testCases {
			s.Run(tc.name, func() {
				req := tc.setupReq()
				result := upgrader.CheckOrigin(req)
				s.Equal(tc.expected, result)
			})
		}
	})
}

func (s *WebSocketTestSuite) handler() handler {
	return handler{
		storage:  s.store,
		config:   s.cfg,
		realtime: s.realtime,
	}
}

func TestWebSocketTestSuite(t *testing.T) {
	suite.Run(t, new(WebSocketTestSuite))
}
