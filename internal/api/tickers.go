package api

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/systemli/ticker/internal/bridge"
	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

//GetTickersHandler returns all Ticker with paging
func GetTickersHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	var tickers []*Ticker
	if me.IsSuperAdmin {
		err = DB.All(&tickers, storm.Reverse())
	} else {
		allowed := me.Tickers
		err = DB.Select(q.In("ID", allowed)).Reverse().Find(&tickers)
		if err == storm.ErrNotFound {
			err = nil
			tickers = []*Ticker{}
		}
	}
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("tickers", NewTickersResponse(tickers)))
}

//GetTickerHandler returns a Ticker for the given id
func GetTickerHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", NewTickerResponse(&ticker)))
}

//GetTickerUsersHandler returns Users for the given ticker
func GetTickerUsersHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	//TODO: Discuss need of Pagination
	users, _ := FindUsersByTicker(ticker)

	c.JSON(http.StatusOK, NewJSONSuccessResponse("users", NewUsersResponse(users)))
}

//PostTickerHandler creates and returns a new Ticker
func PostTickerHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	ticker := NewTicker()
	err := updateTicker(ticker, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = DB.Save(ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", NewTickerResponse(ticker)))
}

//PutTickerHandler updates and returns a existing Ticker
func PutTickerHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = updateTicker(&ticker, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = DB.Save(&ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", NewTickerResponse(&ticker)))
}

//PutTickerUsersHandler changes the allowed users for a ticker
func PutTickerUsersHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var body struct {
		Users []int `json:"users" binding:"required"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = AddUsersToTicker(ticker, body.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	users, _ := FindUsersByTicker(ticker)

	c.JSON(http.StatusOK, NewJSONSuccessResponse("users", NewUsersResponse(users)))
}

//
func PutTickerTwitterHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var body struct {
		Active     bool   `json:"active,omitempty"`
		Disconnect bool   `json:"disconnect"`
		Token      string `json:"token,omitempty"`
		Secret     string `json:"secret,omitempty"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if body.Disconnect {
		ticker.Twitter.Token = ""
		ticker.Twitter.Secret = ""
		ticker.Twitter.Active = false
		ticker.Twitter.User = twitter.User{}
	} else {
		if body.Token != "" {
			ticker.Twitter.Token = body.Token
		}
		if body.Secret != "" {
			ticker.Twitter.Secret = body.Secret
		}
		ticker.Twitter.Active = body.Active
	}

	tu, err := bridge.TwitterUser(&ticker)
	if err != nil {
		log.WithError(err).Error("cant fetch user information from twitter")
	} else {
		ticker.Twitter.User = *tu
	}

	err = DB.Save(&ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", NewTickerResponse(&ticker)))
}

func PutTickerTelegramHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var tg Telegram
	err = c.Bind(&tg)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	ticker.Telegram = tg

	err = DB.Save(&ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", NewTickerResponse(&ticker)))
}

//DeleteTickerHandler deletes a existing Ticker
func DeleteTickerHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	var ticker Ticker
	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	_ = DB.Select(q.Eq("Ticker", tickerID)).Delete(new(Message))
	_ = DB.Select(q.Eq("ID", tickerID)).Delete(new(Ticker))

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
}

//DeleteTickerUserHandler removes ticker credentials for a user
func DeleteTickerUserHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var user User
	err = DB.One("ID", userID, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	err = RemoveTickerFromUser(ticker, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	users, _ := FindUsersByTicker(ticker)

	c.JSON(http.StatusOK, NewJSONSuccessResponse("users", NewUsersResponse(users)))
}

func ResetTickerHandler(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
		return
	}

	var ticker Ticker
	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	err = DeleteMessages(&ticker)
	if err != nil {
		log.WithError(err).WithField("ticker", ticker.ID).Error("error while deleting messages")
	}
	err = DeleteUploadsByTicker(&ticker)
	if err != nil {
		log.WithError(err).WithField("ticker", ticker.ID).Error("error while deleting remaining uploads")
	}

	ticker.Reset()

	err = DB.Save(&ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", NewTickerResponse(&ticker)))
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func updateTicker(t *Ticker, c *gin.Context) error {
	var body struct {
		Domain      string   `json:"domain" binding:"required"`
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Active      bool     `json:"active"`
		Hashtags    []string `json:"hashtags"`
		Information struct {
			Author   string `json:"author"`
			URL      string `json:"url"`
			Email    string `json:"email"`
			Twitter  string `json:"twitter"`
			Facebook string `json:"facebook"`
			Telegram string `json:"telegram"`
		} `json:"information"`
		Location struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}
	}

	err := c.Bind(&body)
	if err != nil {
		return err
	}

	t.Domain = body.Domain
	t.Title = body.Title
	t.Description = body.Description
	t.Active = body.Active
	t.Hashtags = body.Hashtags
	t.Information.Author = body.Information.Author
	t.Information.URL = body.Information.URL
	t.Information.Email = body.Information.Email
	t.Information.Twitter = body.Information.Twitter
	t.Information.Facebook = body.Information.Facebook
	t.Information.Telegram = body.Information.Telegram
	t.Location.Lat = body.Location.Lat
	t.Location.Lon = body.Location.Lon

	return nil
}
