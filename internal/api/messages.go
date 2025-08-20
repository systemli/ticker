package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/pagination"
	"github.com/systemli/ticker/internal/api/realtime"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetMessages(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	pagination := pagination.NewPagination(c)
	messages, err := h.storage.FindMessagesByTickerAndPagination(ticker, *pagination, storage.WithAttachments())
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]any{"messages": response.MessagesResponse(messages, h.config)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) GetMessage(c *gin.Context) {
	message, err := helper.Message(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	data := map[string]any{"message": response.MessageResponse(message, h.config)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PostMessage(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Text        string `json:"text" binding:"required"`
		Attachments []int  `json:"attachments"`
	}
	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	var uploads []storage.Upload
	if len(body.Attachments) > 0 {
		uploads, err = h.storage.FindUploadsByIDs(body.Attachments)
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.UploadsNotFound))
			return
		}
	}

	message := storage.NewMessage()
	message.Text = body.Text
	message.TickerID = ticker.ID
	message.AddAttachments(uploads)

	_ = h.bridges.Send(ticker, &message)

	err = h.storage.SaveMessage(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	serializedMessage := response.MessageResponse(message, h.config)
	h.realtime.Broadcast(realtime.Message{
		Type:     "message_created",
		TickerID: ticker.ID,
		Origin:   helper.GetOriginHost(c),
		Data: map[string]any{
			"message": serializedMessage,
		},
	})

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]any{"message": serializedMessage}))
}

func (h *handler) DeleteMessage(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	message, err := helper.Message(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	_ = h.bridges.Delete(ticker, &message)

	err = h.storage.DeleteMessage(message)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	h.ClearMessagesCache(&ticker)

	h.realtime.Broadcast(realtime.Message{
		Type:     "message_deleted",
		TickerID: ticker.ID,
		Origin:   helper.GetOriginHost(c),
		Data: map[string]any{
			"messageId": message.ID,
		},
	})

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]any{}))
}

// ClearMessagesCache clears the cache for the timeline endpoint of a ticker
func (h *handler) ClearMessagesCache(ticker *storage.Ticker) {
	h.cache.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), fmt.Sprintf("response:%s:/v1/timeline", ticker.Domain)) {
			h.cache.Delete(key)
		}

		return true
	})
}
