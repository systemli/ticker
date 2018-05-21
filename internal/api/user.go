package api

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//GetUsers returns all Users
func GetUsers(c *gin.Context) {
	var users []User

	//TODO: Discuss need of Pagination
	err := DB.All(&users, storm.Reverse())
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("users", NewUsersResponse(users)))
}

//GetUser returns a User for the given id
func GetUser(c *gin.Context) {
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

//PostUser creates and returns a new Ticker
func PostUser(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required" validate:"email"`
		Password string `json:"password" binding:"required" validate:"min=10"`
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

	err = DB.Save(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("user", NewUserResponse(*user)))
}

//PutUser updates a user
func PutUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", userID, &User{})
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	var body struct {
		Email        string `json:"email" validate:"email"`
		Password     string `json:"password" validate:"min=10"`
		Role         string `json:"role"`
		IsSuperAdmin bool   `json:"is_super_admin"`
		Tickers      []int  `json:"tickers"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	var user User
	user.ID = userID

	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		user.UpdatePassword(body.Password)
	}
	if body.Role != "" {
		user.Role = body.Role
	}
	if body.IsSuperAdmin {
		//TODO: Check permissions
		user.IsSuperAdmin = body.IsSuperAdmin
	}

	if len(body.Tickers) > 0 {
		//TODO: Merge existing Tickers
		user.Tickers = body.Tickers
	}

	err = DB.Update(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("user", NewUserResponse(user)))
}

//DeleteUser deletes a existing User
func DeleteUser(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
}
