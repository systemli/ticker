package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetUsers(c *gin.Context) {
	//TODO: Discuss need of Pagination
	users, err := h.storage.FindUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	data := map[string]interface{}{"users": response.UsersResponse(users)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) GetUser(c *gin.Context) {
	me, _ := helper.Me(c)
	user, err := helper.User(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	if !helper.IsAdmin(c) && user.ID != me.ID {
		c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
		return
	}

	data := map[string]interface{}{"user": response.UserResponse(user)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PostUser(c *gin.Context) {
	var body struct {
		Email        string `json:"email,omitempty" binding:"required" validate:"email"`
		Password     string `json:"password,omitempty" binding:"required" validate:"min=10"`
		IsSuperAdmin bool   `json:"isSuperAdmin,omitempty"`
		Tickers      []int  `json:"tickers,omitempty"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	user, err := storage.NewUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	tickers, err := h.storage.FindTickersByIDs(body.Tickers)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	user.IsSuperAdmin = body.IsSuperAdmin
	user.Tickers = tickers

	err = h.storage.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"user": response.UserResponse(user)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) PutUser(c *gin.Context) {
	me, _ := helper.Me(c)
	user, err := helper.User(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Email        string `json:"email,omitempty" validate:"email"`
		Password     string `json:"password,omitempty" validate:"min=10"`
		IsSuperAdmin bool   `json:"isSuperAdmin,omitempty"`
		Tickers      []int  `json:"tickers,omitempty"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		user.UpdatePassword(body.Password)
	}

	// You only can set/unset other users SuperAdmin property
	if me.ID != user.ID {
		user.IsSuperAdmin = body.IsSuperAdmin
	}

	if body.Tickers != nil {
		tickers, err := h.storage.FindTickersByIDs(body.Tickers)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
			return
		}

		user.Tickers = tickers
	}

	err = h.storage.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"user": response.UserResponse(user)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}

func (h *handler) DeleteUser(c *gin.Context) {
	me, _ := helper.Me(c)
	user, err := helper.User(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	if me.ID == user.ID {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, "self deletion is forbidden"))
		return
	}

	err = h.storage.DeleteUser(user)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(nil))
}

func (h *handler) PutMe(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeDefault, response.Unauthorized))
		return
	}

	var body struct {
		Password    string `json:"password" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required" validate:"min=10"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	if !me.Authenticate(body.Password) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.PasswordError))
		return
	}

	me.UpdatePassword(body.NewPassword)

	err = h.storage.SaveUser(&me)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"user": response.UserResponse(me)}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
