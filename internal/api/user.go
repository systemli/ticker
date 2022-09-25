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
		IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
		Tickers      []int  `json:"tickers,omitempty"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	//TODO: Validation
	user, err := storage.NewUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	user.IsSuperAdmin = body.IsSuperAdmin
	user.Tickers = body.Tickers

	err = h.storage.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"user": user}
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
		Role         string `json:"role,omitempty"`
		IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
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
	if body.Role != "" {
		user.Role = body.Role
	}

	// You only can set/unset other users SuperAdmin property
	if me.ID != user.ID {
		user.IsSuperAdmin = body.IsSuperAdmin
	}

	if body.Tickers != nil {
		user.Tickers = body.Tickers
	}

	err = h.storage.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	data := map[string]interface{}{"user": user}
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
