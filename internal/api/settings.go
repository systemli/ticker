package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
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

func getInactiveSettings(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	setting := GetInactiveSettings()
	c.JSON(http.StatusOK, NewJSONSuccessResponse("setting", NewSettingResponse(setting)))
}
