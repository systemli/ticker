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

	var setting storage.Setting
	if c.Param("name") == storage.SettingInactiveName {
		setting = h.storage.GetInactiveSetting()
	}

	if c.Param("name") == storage.SettingRefreshInterval {
		setting = h.storage.GetRefreshIntervalSetting()
	}

	data := map[string]interface{}{"setting": response.SettingResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PutInactiveSettings(c *gin.Context) {
	value := storage.InactiveSettings{}
	err := c.Bind(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.storage.SaveInactiveSetting(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	setting := h.storage.GetInactiveSetting()
	data := map[string]interface{}{"setting": response.SettingResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PutRefreshInterval(c *gin.Context) {
	value := storage.RefreshIntervalSettings{}
	err := c.Bind(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.storage.SaveRefreshInterval(value.RefreshInterval)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	setting := h.storage.GetRefreshIntervalSetting()
	data := map[string]interface{}{"setting": response.SettingResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
