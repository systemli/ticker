package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetSetting(c *gin.Context) {
	if !helper.IsAdmin(c) {
		c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
		return
	}

	if c.Param("name") == storage.SettingInactiveName {
		setting := h.storage.GetInactiveSettings()
		data := map[string]interface{}{"setting": response.InactiveSettingsResponse(setting)}
		c.JSON(http.StatusOK, response.SuccessResponse(data))
		return
	}

	c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.SettingNotFound))
}

func (h *handler) PutInactiveSettings(c *gin.Context) {
	value := storage.InactiveSettings{}
	err := c.Bind(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.storage.SaveInactiveSettings(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	setting := h.storage.GetInactiveSettings()
	data := map[string]interface{}{"setting": response.InactiveSettingsResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
