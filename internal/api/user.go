package api

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//GetUsersHandler returns all Users
func GetUsersHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorInsufficientPermissions, "insufficient permissions"))
		return
	}

	var users []User

	//TODO: Discuss need of Pagination
	err := DB.All(&users, storm.Reverse())
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("users", NewUsersResponse(users)))
}

//GetUserHandler returns a User for the given id
func GetUserHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorInsufficientPermissions, "insufficient permissions"))
		return
	}

	var user User
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", userID, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("user", NewUserResponse(user)))
}

//PostUserHandler creates and returns a new Ticker
func PostUserHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorInsufficientPermissions, "insufficient permissions"))
		return
	}

	var body struct {
		Email        string `json:"email,omitempty" binding:"required" validate:"email"`
		Password     string `json:"password,omitempty" binding:"required" validate:"min=10"`
		IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	//TODO: Validation

	user, err := NewUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	user.IsSuperAdmin = body.IsSuperAdmin

	err = DB.Save(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("user", NewUserResponse(*user)))
}

//PutUserHandler updates a user
func PutUserHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorInsufficientPermissions, "insufficient permissions"))
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	var user User
	err = DB.One("ID", userID, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
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
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
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
	//TODO: Check permissions
	user.IsSuperAdmin = body.IsSuperAdmin

	if body.Tickers != nil {
		//TODO: Merge existing Tickers
		user.Tickers = body.Tickers
	}

	err = DB.Save(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("user", NewUserResponse(user)))
}

//DeleteUserHandler deletes a existing User
func DeleteUserHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorInsufficientPermissions, "insufficient permissions"))
		return
	}

	var user User
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", userID, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	err = DB.DeleteStruct(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
}
