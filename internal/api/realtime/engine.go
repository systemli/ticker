package realtime

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/systemli/ticker/internal/logger"
)

const (
	// Time allowed writing a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed reading the next pong message from the peer.
	pongWait = 60 * time.Second

	// Time allowed for the client to close the connection gracefully.
	clientCloseWait = 100 * time.Millisecond

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var log = logger.GetWithPackage("realtime")

// Engine is the core component of the real-time messaging system. It manages WebSocket
// connections, tracks active clients, and facilitates message broadcasting. The Engine
// maintains a mapping of ticker IDs to their connected clients, allowing targeted
// communication for specific tickers.
//
// Key responsibilities:
//   - Registering and unregistering clients: The Engine tracks clients as they connect
//     and disconnect, ensuring proper resource management.
//   - Broadcasting messages: Messages are sent to all relevant clients based on their
//     associated ticker IDs.
//   - Managing WebSocket connections: The Engine handles the lifecycle of WebSocket
//     connections, including sending and receiving messages.
//
// Relationship with other components:
//   - Client: Represents an individual WebSocket connection. Each client is associated
//     with a specific ticker ID and communicates with the Engine for message exchange.
//   - Message: Represents the data structure for messages sent between the Engine and
//     clients, including metadata like type and ticker ID.
type Engine struct {
	clients      map[int]map[*Client]bool // clients maps ticker IDs to their connected clients
	broadcast    chan Message
	register     chan *Client
	unregister   chan *Client
	shutdown     chan struct{}
	done         chan struct{}
	running      bool
	shuttingDown bool
	mu           sync.RWMutex
}

// Client represents a single connection to the Engine. Each Client is associated with a specific TickerID
// and is responsible for sending and receiving messages over a WebSocket connection.
//
// Lifecycle:
// - A Client is created when a new WebSocket connection is established.
// - It is registered with the Engine, which manages all active clients.
// - The Client listens for messages to send via the `Send` channel and writes them to the WebSocket.
// - When the connection is closed, the Client is unregistered and cleaned up.
//
// Fields:
//   - Engine: A reference to the Engine managing this Client.
//   - Conn: The WebSocket connection associated with this Client.
//   - Send: A channel for outgoing messages. It is buffered to prevent blocking the Engine
//     when sending messages. The buffer size should be chosen based on expected message
//     volume and latency requirements.
//   - TickerID: The ID of the ticker this Client is subscribed to.
//   - Origin: The origin of the WebSocket connection.
//   - closed: A flag indicating whether the Client has been closed.
//   - mu: A mutex to protect concurrent access to the Client's fields.
//   - unregisterOnce: Ensures unregistration happens only once.
type Client struct {
	Engine         *Engine
	Conn           *websocket.Conn
	Send           chan Message
	TickerID       int
	Origin         string
	closed         bool
	mu             sync.Mutex
	unregisterOnce sync.Once
}

// Message represents a message sent to clients in the realtime engine.
//
// Fields:
// - Type: A string indicating the type of message. Expected values include:
//   - "message_deleted": Indicates that a message was deleted.
//   - "message_created": Indicates that a new message was created.
//     Additional types may be added as needed.
//   - TickerID: An integer representing the ID of the ticker associated with the message.
//   - Data: A flexible field of type `any` that contains additional data related to the message.
//   - Origin: A string representing the origin of the message, used for logging and metrics.
//     The structure of this data depends on the `Type` field. For example:
//   - For "message_created", `Data` might include the content of the new message.
//   - For "message_deleted", `Data` might include the ID of the deleted message.
type Message struct {
	Type     string `json:"type"`
	TickerID int    `json:"tickerId"`
	Data     any    `json:"data"`
	Origin   string `json:"-"` // The origin of the message, used for logging and metrics
}

// New creates a new realtime messaging engine.
func New() *Engine {
	return &Engine{
		clients:    make(map[int]map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		shutdown:   make(chan struct{}),
		done:       make(chan struct{}),
	}
}

// Run starts the realtime engine's main loop.
func (e *Engine) Run() {
	defer func() {
		e.mu.Lock()
		defer e.mu.Unlock()

		e.running = false

		// Ensure channels are closed properly to prevent deadlocks
		select {
		case <-e.done:
			// Channel was already closed
		default:
			close(e.done)
		}
	}()

	e.mu.Lock()
	e.running = true
	e.mu.Unlock()

	log.Info("websocket engine started")

	for {
		select {
		case client := <-e.register:
			e.registerClient(client)

		case client := <-e.unregister:
			e.unregisterClient(client)

		case message := <-e.broadcast:
			e.broadcastMessage(message)

		case <-e.shutdown:
			log.Info("websocket engine shutting down")
			e.mu.Lock()
			e.shuttingDown = true
			e.mu.Unlock()
			e.closeAllConnections()
			return
		}
	}
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown(ctx context.Context) error {
	log.Info("initiating websocket engine shutdown")

	e.mu.Lock()
	isRunning := e.running
	isShuttingDown := e.shuttingDown

	if !isRunning {
		e.mu.Unlock()
		log.Info("websocket engine was not running")
		return nil
	}

	if isShuttingDown {
		e.mu.Unlock()
		log.Info("websocket engine already shutting down")
		// Wait for the engine to finish shutting down or context timeout
		select {
		case <-e.done:
			log.Info("websocket engine shutdown completed")
			return nil
		case <-ctx.Done():
			log.Warn("websocket engine shutdown timed out")
			return ctx.Err()
		}
	}

	// Mark as shutting down and signal the engine to shut down
	e.shuttingDown = true
	e.mu.Unlock()

	// Signal the engine to shut down safely
	defer func() {
		if r := recover(); r != nil {
			// Channel was already closed, that's fine
			log.Debug("shutdown channel was already closed")
		}
	}()
	close(e.shutdown)

	// Wait for the engine to finish shutting down or context timeout
	select {
	case <-e.done:
		log.Info("websocket engine shutdown completed")
		return nil
	case <-ctx.Done():
		log.Warn("websocket engine shutdown timed out")
		// Force close in a separate goroutine to avoid blocking on locks
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.WithField("panic", r).Error("panic during force close")
				}
			}()
			e.forceCloseAllConnections()
		}()
		return ctx.Err()
	}
}

// safeCloseClient safely closes a client's sent channel only once
func (e *Engine) safeCloseClient(client *Client) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if !client.closed {
		close(client.Send)
		client.closed = true
	}
}

// closeAllConnections gracefully closes all client connections
func (e *Engine) closeAllConnections() {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Count clients for logging
	clientCount := 0
	for _, clients := range e.clients {
		clientCount += len(clients)
	}
	log.WithField("client_count", clientCount).Info("closing all websocket connections")

	// Send close messages to all clients and close their sending channels
	for tickerID, clients := range e.clients {
		for client := range clients {
			// Send a close message to a client
			select {
			case client.Send <- Message{
				Type:     "server_shutdown",
				TickerID: tickerID,
				Data:     map[string]any{"message": "server is shutting down"},
			}:
				recordMessageSent(client.Origin, "server_shutdown")
			default:
				// Channel might be full, skip
				recordMessageDropped(client.Origin, "server_shutdown")
			}

			// Close the channel to signal WritePump to send a close message
			e.safeCloseClient(client)

			// Record disconnection metric
			recordClientDisconnected(client.Origin, "server_shutdown")
		}
	}

	// Give clients a moment to receive the close message (release lock temporarily)
	e.mu.Unlock()
	time.Sleep(clientCloseWait)
	e.mu.Lock()

	// Force close any remaining connections (already have lock)
	e.forceCloseAllConnectionsUnsafe()
}

// forceCloseAllConnections forcefully closes all connections
func (e *Engine) forceCloseAllConnections() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.forceCloseAllConnectionsUnsafe()
}

// forceCloseAllConnectionsUnsafe forcefully closes all connections (assumes lock is held)
func (e *Engine) forceCloseAllConnectionsUnsafe() {
	for _, clients := range e.clients {
		for client := range clients {
			_ = client.Conn.Close()
		}
	}

	// Clear all clients
	e.clients = make(map[int]map[*Client]bool)
}

// Broadcast sends a message to all relevant clients
func (e *Engine) Broadcast(message Message) {
	log.WithFields(logrus.Fields{"message_type": message.Type, "ticker_id": message.TickerID}).Debug("broadcasting message")

	select {
	case e.broadcast <- message:
	default:
		log.WithFields(logrus.Fields{"message_type": message.Type, "ticker_id": message.TickerID}).Warn("broadcast channel full, dropping message")
	}
}

// Register queues a client for registration
func (e *Engine) Register(client *Client) {
	select {
	case e.register <- client:
	default:
		log.WithField("origin", client.Origin).Warn("register channel full, cannot register client")
	}
}

// registerClient handles the actual client registration (called from Run loop)
func (e *Engine) registerClient(client *Client) {
	log.WithField("origin", client.Origin).Debug("registering client")

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.clients[client.TickerID] == nil {
		e.clients[client.TickerID] = make(map[*Client]bool)
	}
	e.clients[client.TickerID][client] = true

	// Record metrics for new connection
	recordClientConnected(client.Origin)
}

// Unregister queues a client for unregistration
func (e *Engine) Unregister(client *Client) {
	select {
	case e.unregister <- client:
	default:
		log.WithField("origin", client.Origin).Warn("unregister channel full, cannot unregister client")
	}
}

// unregisterClient removes a client and closes its sent channel
func (e *Engine) unregisterClient(client *Client) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if clients, exists := e.clients[client.TickerID]; exists {
		if _, exists := clients[client]; exists {
			delete(clients, client)

			// Close the client's sent channel safely
			e.safeCloseClient(client)

			// Record metrics for disconnection
			recordClientDisconnected(client.Origin, "normal")

			// If no more clients exist for this ticker, delete the map
			if len(clients) == 0 {
				delete(e.clients, client.TickerID)
			}
		}
	}
}

// broadcastMessage sends a message to all clients of a specific ticker
func (e *Engine) broadcastMessage(message Message) {
	start := time.Now()

	e.mu.Lock()
	defer e.mu.Unlock()

	// Send a message only to clients of the specific ticker
	if clients, exists := e.clients[message.TickerID]; exists {
		var deadClients []*Client
		sentCount := 0
		droppedCount := 0

		for client := range clients {
			select {
			case client.Send <- message:
				sentCount++
			default:
				// Client cannot receive, mark for removal
				deadClients = append(deadClients, client)
				droppedCount++
			}
		}

		// Remove dead clients
		for _, deadClient := range deadClients {
			delete(clients, deadClient)
			e.safeCloseClient(deadClient)
			recordClientDisconnected(deadClient.Origin, "channel_full")
		}

		// Clean up empty client maps
		if len(clients) == 0 {
			delete(e.clients, message.TickerID)
		}

		// Record metrics
		for i := 0; i < sentCount; i++ {
			recordMessageSent(message.Origin, message.Type)
		}
		for i := 0; i < droppedCount; i++ {
			recordMessageDropped(message.Origin, message.Type)
		}
		recordBroadcastDuration(message.Origin, message.Type, time.Since(start))
	}
}

// unregisterSafely ensures that the client is unregistered only once
func (c *Client) unregisterSafely() {
	c.unregisterOnce.Do(func() {
		c.Engine.Unregister(c)
	})
}

// WritePump pumps messages from the client to the engine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.unregisterSafely()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The engine closed the channel - send close message and exit
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.WithError(err).WithField("origin", c.Origin).Warn("error writing to websocket")
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithError(err).WithField("origin", c.Origin).Warn("error sending ping")
				return
			}
		}
	}
}

// ReadPump handles the read side of the WebSocket connection.
// It's optimized for a broadcast-only system - we don't process incoming messages
// but need to handle connection health monitoring and proper cleanup.
func (c *Client) ReadPump() {
	defer func() {
		c.unregisterSafely()
		_ = c.Conn.Close()
	}()

	// Set read limit to prevent large message attacks
	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// Handle pong messages for connection health monitoring
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Optimized read loop: we only care about connection state, not message content
	for {
		// Use NextReader instead of ReadJSON to avoid the overhead of JSON unmarshaling.
		messageType, reader, err := c.Conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).WithField("origin", c.Origin).Info("websocket connection closed unexpectedly")
			}
			return
		}

		// Handle control messages properly
		if messageType == websocket.CloseMessage {
			return
		}

		// For text/binary messages: discard efficiently using io.Discard.
		// This system is broadcast-only, so incoming messages are not processed.
		// Discarding them ensures efficient resource usage without unnecessary allocations.
		if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
			_, _ = io.Copy(io.Discard, reader)
		}
	}
}
