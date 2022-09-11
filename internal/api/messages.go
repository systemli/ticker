package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/api/util"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetMessages(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
			return
		}
	}

	ticker, err := h.storage.FindTickerByID(tickerID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.TickerNotFound))
		return
	}

	//TODO: Pagination
	messages, err := h.storage.FindMessagesByTicker(ticker, util.Pagination{})
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"messages": []string{}}))
			return
		}

		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"messages": response.MessagesResponse(messages, h.config)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) GetMessage(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
			return
		}
	}

	messageID, err := strconv.Atoi(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.MessageIdentierMissing))
		return
	}
	message, err := h.storage.FindMessage(tickerID, messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.MessageNotFound))
		return
	}

	data := map[string]interface{}{"message": response.MessageResponse(message, h.config)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PostMessage(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
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

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
			return
		}
	}

	ticker, err := h.storage.FindTickerByID(tickerID)
	if err != nil {
		log.WithError(err).Error("failed to find the ticker")
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.TickerNotFound))
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
	message.Ticker = tickerID
	message.GeoInformation = body.GeoInformation
	message.AddAttachments(uploads)

	_ = h.bridges.Send(ticker, message)

	err = h.storage.SaveMessage(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"message": response.MessageResponse(message, h.config)}))
}

func (h *handler) DeleteMessage(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
			return
		}
	}

	ticker, err := h.storage.FindTickerByID(tickerID)
	if err != nil {
		log.WithError(err).Error("failed to find the ticker")
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.TickerNotFound))
		return
	}

	messageID, err := strconv.Atoi(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.MessageIdentierMissing))
		return
	}

	message, err := h.storage.FindMessage(tickerID, messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.MessageNotFound))
		return
	}

	_ = h.bridges.Delete(ticker, message)

	err = h.storage.DeleteMessage(message)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{}))
}
