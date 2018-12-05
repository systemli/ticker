package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

//GetSettingHandler returns a Setting
func GetSettingHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	if c.Param("name") == SettingInactiveName {
		getInactiveSettings(c)
		return
	}

	if c.Param("name") == SettingRefreshInterval {
		getRefreshInterval(c)
		return
	}

	setting, err := FindSetting(c.Param("name"))
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, ErrorSettingNotFound))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("setting", NewSettingResponse(setting)))
}

//PutInactiveSettingsHandler updates the inactive_settings
func PutInactiveSettingsHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	value := &InactiveSettings{}
	err := c.Bind(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var setting Setting
	err = DB.One("Name", SettingInactiveName, &setting)
	if err != nil {
		setting.Name = SettingInactiveName
	}

	setting.Value = value
	err = DB.Save(&setting)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("setting", NewSettingResponse(&setting)))
}

//PutRefreshIntervalHandler updates refresh_interval
func PutRefreshIntervalHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	var payload struct {
		RefreshInterval int `json:"refresh_interval"`
	}
	err := c.Bind(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var setting Setting
	err = DB.One("Name", SettingRefreshInterval, &setting)
	if err != nil {
		setting.Name = SettingRefreshInterval
	}

	setting.Value = payload.RefreshInterval
	err = DB.Save(&setting)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("setting", NewSettingResponse(&setting)))
}

func getInactiveSettings(c *gin.Context) {
	setting := GetInactiveSettings()
	c.JSON(http.StatusOK, NewJSONSuccessResponse("setting", NewSettingResponse(setting)))
}

func getRefreshInterval(c *gin.Context) {
	setting := GetRefreshInterval()
	c.JSON(http.StatusOK, NewJSONSuccessResponse("setting", NewSettingResponse(setting)))
}
