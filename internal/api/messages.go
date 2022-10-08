package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetMessages(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	//TODO: Pagination
	messages, err := h.storage.FindMessagesByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"messages": response.MessagesResponse(messages, h.config)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) GetMessage(c *gin.Context) {
	message, err := helper.Message(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	data := map[string]interface{}{"message": response.MessageResponse(message, h.config)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PostMessage(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Text           string                    `json:"text" binding:"required"`
		GeoInformation geojson.FeatureCollection `json:"geo_information"`
		Attachments    []int                     `json:"attachments"`
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
	message.Ticker = ticker.ID
	message.GeoInformation = body.GeoInformation
	message.AddAttachments(uploads)

	_ = h.bridges.Send(ticker, &message)

	err = h.storage.SaveMessage(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"message": response.MessageResponse(message, h.config)}))
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

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{}))
}
