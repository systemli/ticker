package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/realtime"
	"github.com/systemli/ticker/internal/api/response"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin != "" || r.URL.Query().Has("origin")
	},
}

func (h *handler) HandleWebSocket(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	log.WithField("ticker_id", ticker.ID).Info("New WebSocket connection attempt")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("WebSocket upgrade failed")
		return
	}

	client := &realtime.Client{
		Engine:   h.realtime,
		Conn:     conn,
		Send:     make(chan realtime.Message, 256), // Buffer to prevent blocking
		TickerID: ticker.ID,
	}

	h.realtime.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
