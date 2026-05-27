package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/signal"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetSetting(c *gin.Context) {
	if !helper.IsAdmin(c) {
		c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
		return
	}

	if c.Param("name") == storage.SettingInactiveName {
		setting := storage.GetSettings(h.stores.Settings, storage.InactiveSetting)
		data := map[string]interface{}{"setting": response.InactiveSettingsResponse(setting)}
		c.JSON(http.StatusOK, response.SuccessResponse(data))
		return
	}

	if c.Param("name") == storage.SettingTelegramName {
		setting := storage.GetSettings(h.stores.Settings, storage.TelegramSetting)
		data := map[string]interface{}{"setting": response.TelegramSettingsResponse(setting)}
		c.JSON(http.StatusOK, response.SuccessResponse(data))
		return
	}

	if c.Param("name") == storage.SettingSignalGroupName {
		setting := storage.GetSettings(h.stores.Settings, storage.SignalGroupSetting)
		data := map[string]interface{}{"setting": response.SignalGroupSettingsResponse(setting)}
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

	err = storage.SaveSettings(h.stores.Settings, storage.InactiveSetting, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	setting := storage.GetSettings(h.stores.Settings, storage.InactiveSetting)
	data := map[string]interface{}{"setting": response.InactiveSettingsResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PutTelegramSettings(c *gin.Context) {
	value := storage.TelegramSettings{}
	err := c.Bind(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	// Validate the token by calling the Telegram API if a token is provided
	if value.Token != "" {
		botUser, err := bridge.BotUser(value.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeBadCredentials, response.TelegramError))
			return
		}
		value.BotUsername = botUser.UserName
	}

	err = storage.SaveSettings(h.stores.Settings, storage.TelegramSetting, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	setting := storage.GetSettings(h.stores.Settings, storage.TelegramSetting)
	data := map[string]interface{}{"setting": response.TelegramSettingsResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PutSignalGroupSettings(c *gin.Context) {
	value := storage.SignalGroupSettings{}
	err := c.Bind(&value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	// Validate the connection by calling listGroups if ApiUrl and Account are provided
	if value.ApiUrl != "" && value.Account != "" {
		groupClient := signal.NewGroupClientFromSettings(value)
		_, err := groupClient.ListGroups()
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeBadCredentials, response.SignalGroupError))
			return
		}
	}

	err = storage.SaveSettings(h.stores.Settings, storage.SignalGroupSetting, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	setting := storage.GetSettings(h.stores.Settings, storage.SignalGroupSetting)
	data := map[string]interface{}{"setting": response.SignalGroupSettingsResponse(setting)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
