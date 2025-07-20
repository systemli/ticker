package realtime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
)

type EngineTestSuite struct {
	suite.Suite
	engine *Engine
}

func (s *EngineTestSuite) SetupTest() {
	s.engine = New()
}

func (s *EngineTestSuite) TearDownTest() {
	if s.engine != nil {
		// Ensure hub is shut down after each test with a short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		// Try to shutdown, but don't fail the test if it times out
		// (some tests intentionally don't start the hub)
		_ = s.engine.Shutdown(ctx)
	}
}

func (s *EngineTestSuite) TestShutdown() {
	s.Run("shutdown without clients", func() {
		engine := New()
		go engine.Run()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := engine.Shutdown(ctx)
		s.NoError(err)
	})

	s.Run("shutdown with timeout", func() {
		engine := New()
		go engine.Run()

		// Immediately try to shutdown with a very short timeout
		// This might timeout if the hub is busy
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		err := engine.Shutdown(ctx)
		// Should either succeed or timeout - both are valid
		if err != nil {
			s.Equal(context.DeadlineExceeded, err)
		}
	})

	s.Run("multiple shutdowns should not panic", func() {
		engine := New()
		go engine.Run()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// First shutdown
		err1 := engine.Shutdown(ctx)
		s.NoError(err1)

		// Second shutdown should not panic
		s.NotPanics(func() {
			_ = engine.Shutdown(ctx)
		})
	})

	s.Run("shutdown unstarted hub should not hang", func() {
		engine := New()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// This should return immediately without hanging
		start := time.Now()
		err := engine.Shutdown(ctx)
		elapsed := time.Since(start)

		s.NoError(err)
		s.Less(elapsed, 50*time.Millisecond, "Shutdown should return immediately for non-running hub")
	})

	s.Run("multiple shutdowns timing", func() {
		engine := New()
		go engine.Run()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// First shutdown
		err1 := engine.Shutdown(ctx)
		s.NoError(err1)

		// Second shutdown
		err2 := engine.Shutdown(ctx)
		s.NoError(err2)
	})
}

func (s *EngineTestSuite) TestBroadcast() {
	s.Run("broadcast to empty hub", func() {
		engine := New()
		go engine.Run()

		// Should not panic when broadcasting to empty hub
		s.NotPanics(func() {
			engine.Broadcast(Message{
				Type:     "test",
				TickerID: 1,
				Data:     map[string]any{"test": "data"},
			})
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = engine.Shutdown(ctx)
	})

	s.Run("broadcast with full channel should drop message", func() {
		engine := New()
		go engine.Run()

		// Fill the broadcast channel by not processing messages
		// This tests the non-blocking broadcast behavior
		for i := 0; i < 100; i++ {
			engine.Broadcast(Message{
				Type:     "spam",
				TickerID: 1,
				Data:     map[string]any{"count": i},
			})
		}

		// Should still not panic or block
		s.NotPanics(func() {
			engine.Broadcast(Message{
				Type:     "final",
				TickerID: 1,
				Data:     map[string]any{"last": true},
			})
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = engine.Shutdown(ctx)
	})
}

func (s *EngineTestSuite) TestClientManagement() {
	s.Run("register and unregister clients", func() {
		engine := New()
		go engine.Run()

		// Create mock websocket connections
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true },
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			s.NoError(err)
			defer conn.Close()

			client := &Client{
				Engine:   engine,
				Conn:     conn,
				Send:     make(chan Message, 256),
				TickerID: 1,
			}

			engine.Register(client)
			go client.WritePump()
			go client.ReadPump()

			// Wait for the request context to be cancelled or connection to close
			<-r.Context().Done()
		}))
		defer server.Close()

		// Connect to test server
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		s.NoError(err)
		defer conn.Close()

		// Give time for registration
		time.Sleep(200 * time.Millisecond)

		// Verify client count (we can't directly access the clients map, but we can test behavior)
		engine.Broadcast(Message{
			Type:     "test_registration",
			TickerID: 1,
			Data:     map[string]any{"message": "client registered"},
		})

		// Should be able to read the broadcast message
		var receivedMsg Message
		conn.SetReadDeadline(time.Now().Add(time.Second))
		err = conn.ReadJSON(&receivedMsg)
		s.NoError(err)
		s.Equal("test_registration", receivedMsg.Type)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = engine.Shutdown(ctx)
	})
}

func (s *EngineTestSuite) TestGracefulShutdownWithClients() {
	s.Run("shutdown sends close message to clients", func() {
		engine := New()
		go engine.Run()

		// Create a test server that handles WebSocket connections
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true },
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			s.NoError(err)
			defer conn.Close()

			client := &Client{
				Engine:   engine,
				Conn:     conn,
				Send:     make(chan Message, 256),
				TickerID: 1,
			}

			engine.Register(client)
			go client.WritePump()
			go client.ReadPump()

			// Wait for the request context to be cancelled or connection to close
			<-r.Context().Done()
		}))
		defer server.Close()

		// Connect to test server
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		s.NoError(err)
		defer conn.Close()

		// Give time for connection to establish
		time.Sleep(100 * time.Millisecond)

		// Start shutdown in a goroutine
		go func() {
			time.Sleep(50 * time.Millisecond) // Small delay to ensure test is ready
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = engine.Shutdown(ctx)
		}()

		// Should receive shutdown message
		var shutdownMsg Message
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		err = conn.ReadJSON(&shutdownMsg)
		s.NoError(err)
		s.Equal("server_shutdown", shutdownMsg.Type)
		s.Equal(1, shutdownMsg.TickerID)

		// Connection should close after shutdown message
		_, _, err = conn.ReadMessage()
		s.Error(err) // Should get close error
	})
}

func (s *EngineTestSuite) TestConcurrentOperations() {
	s.Run("concurrent broadcasts and shutdowns", func() {
		engine := New()
		go engine.Run()

		// Start multiple goroutines broadcasting
		done := make(chan struct{})

		for i := 0; i < 5; i++ {
			go func(id int) {
				defer func() { done <- struct{}{} }()
				for j := 0; j < 10; j++ {
					engine.Broadcast(Message{
						Type:     "concurrent_test",
						TickerID: id,
						Data:     map[string]any{"worker": id, "message": j},
					})
					time.Sleep(10 * time.Millisecond)
				}
			}(i)
		}

		// Let broadcasts run for a moment
		time.Sleep(50 * time.Millisecond)

		// Shutdown should handle concurrent operations gracefully
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := engine.Shutdown(ctx)
		s.NoError(err)

		// Wait for all goroutines to complete (they should complete quickly after shutdown)
		timeout := time.After(2 * time.Second)
		completed := 0
		for completed < 5 {
			select {
			case <-done:
				completed++
			case <-timeout:
				s.Fail("Goroutines did not complete in time after shutdown")
				return
			}
		}
	})
}

func (s *EngineTestSuite) TestWebSocketIntegration() {
	s.Run("full websocket lifecycle with graceful shutdown", func() {
		engine := New()
		go engine.Run()

		// Create a test WebSocket server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true },
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			s.NoError(err)
			defer conn.Close()

			client := &Client{
				Engine:   engine,
				Conn:     conn,
				Send:     make(chan Message, 256),
				TickerID: 1,
			}

			engine.Register(client)
			go client.WritePump()
			go client.ReadPump()

			// Wait for the request context to be cancelled or connection to close
			<-r.Context().Done()
		}))
		defer server.Close()

		// Connect to the test server
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		s.NoError(err)
		defer conn.Close()

		// Give some time for connection to establish
		time.Sleep(100 * time.Millisecond)

		// Test broadcasting a regular message
		engine.Broadcast(Message{
			Type:     "test_message",
			TickerID: 1,
			Data:     map[string]any{"content": "Hello WebSocket!"},
		})

		// Read the message
		var receivedMsg Message
		conn.SetReadDeadline(time.Now().Add(time.Second))
		err = conn.ReadJSON(&receivedMsg)
		s.NoError(err)
		s.Equal("test_message", receivedMsg.Type)
		s.Equal(1, receivedMsg.TickerID)
		s.Equal("Hello WebSocket!", receivedMsg.Data.(map[string]any)["content"])

		// Test graceful shutdown
		shutdownDone := make(chan error, 1)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			shutdownDone <- engine.Shutdown(ctx)
		}()

		// Should receive shutdown message
		var shutdownMsg Message
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		err = conn.ReadJSON(&shutdownMsg)
		s.NoError(err)
		s.Equal("server_shutdown", shutdownMsg.Type)

		// Wait for shutdown to complete
		select {
		case err := <-shutdownDone:
			s.NoError(err)
		case <-time.After(3 * time.Second):
			s.Fail("Shutdown took too long")
		}

		// Connection should be closed after shutdown
		_, _, err = conn.ReadMessage()
		s.Error(err) // Should get close error
	})

	s.Run("multiple clients different tickers", func() {
		engine := New()
		go engine.Run()

		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true },
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			s.NoError(err)
			defer conn.Close()

			// Get ticker ID from query param
			tickerID := 1
			if r.URL.Query().Get("ticker") == "2" {
				tickerID = 2
			}

			client := &Client{
				Engine:   engine,
				Conn:     conn,
				Send:     make(chan Message, 256),
				TickerID: tickerID,
			}

			engine.Register(client)
			go client.WritePump()
			go client.ReadPump()

			// Wait for the request context to be cancelled or connection to close
			<-r.Context().Done()
		}))
		defer server.Close()

		// Connect two clients with different ticker IDs
		wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "?ticker=1"
		wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + "?ticker=2"

		conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
		s.NoError(err)
		defer conn1.Close()

		conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
		s.NoError(err)
		defer conn2.Close()

		time.Sleep(100 * time.Millisecond) // Let connections establish

		// Broadcast to ticker 1 only
		engine.Broadcast(Message{
			Type:     "ticker1_message",
			TickerID: 1,
			Data:     map[string]any{"content": "Only for ticker 1"},
		})

		// Client 1 should receive the message
		var msg1 Message
		conn1.SetReadDeadline(time.Now().Add(time.Second))
		err = conn1.ReadJSON(&msg1)
		s.NoError(err)
		s.Equal("ticker1_message", msg1.Type)
		s.Equal(1, msg1.TickerID)

		// Client 2 should not receive the message (should timeout)
		conn2.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		var msg2 Message
		err = conn2.ReadJSON(&msg2)
		s.Error(err) // Should timeout

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = engine.Shutdown(ctx)
	})
}

func (s *EngineTestSuite) TestIntegrationWebSocketConnection() {
	s.Run("websocket connection and shutdown", func() {
		engine := New()
		go engine.Run()

		// Create a test WebSocket server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool { return true },
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			s.NoError(err)
			defer conn.Close()

			client := &Client{
				Engine:   engine,
				Conn:     conn,
				Send:     make(chan Message, 256),
				TickerID: 1,
			}

			engine.Register(client)
			go client.WritePump()
			go client.ReadPump()

			// Wait for the request context to be cancelled or connection to close
			<-r.Context().Done()
		}))
		defer server.Close()

		// Connect to the test server
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		s.NoError(err)
		defer conn.Close()

		// Give some time for connection to establish
		time.Sleep(100 * time.Millisecond)

		// Test broadcasting
		engine.Broadcast(Message{
			Type:     "test_message",
			TickerID: 1,
			Data:     map[string]any{"content": "Hello WebSocket!"},
		})

		// Read the message
		var receivedMsg Message
		conn.SetReadDeadline(time.Now().Add(time.Second))
		err = conn.ReadJSON(&receivedMsg)
		s.NoError(err)

		s.Equal("test_message", receivedMsg.Type)
		s.Equal(1, receivedMsg.TickerID)
		s.Equal("Hello WebSocket!", receivedMsg.Data.(map[string]any)["content"])

		// Test graceful shutdown
		shutdownDone := make(chan error, 1)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			shutdownDone <- engine.Shutdown(ctx)
		}()

		// Should receive shutdown message
		var shutdownMsg Message
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		err = conn.ReadJSON(&shutdownMsg)
		s.NoError(err)
		s.Equal("server_shutdown", shutdownMsg.Type)

		// Wait for shutdown to complete
		select {
		case err := <-shutdownDone:
			s.NoError(err)
		case <-time.After(3 * time.Second):
			s.Fail("Shutdown took too long")
		}

		// Connection should be closed after shutdown
		_, _, err = conn.ReadMessage()
		s.Error(err)
	})
}

func TestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}
