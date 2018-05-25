package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/dghubble/go-twitter/twitter"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
	"git.codecoop.org/systemli/ticker/internal/bridge"
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

	c.JSON(http.StatusOK, NewJSONSuccessResponse("tickers", NewTickersReponse(tickers)))
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

	if ticker.Twitter.Connected() {
		user, err := bridge.Twitter.User(ticker)
		if err == nil {
			ticker.Twitter.User = *user
		}
	}

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

	DB.Select(q.Eq("Ticker", tickerID)).Delete(new(Message))
	DB.Select(q.Eq("ID", tickerID)).Delete(new(Ticker))

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
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
		Domain      string `json:"domain" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Active      bool   `json:"active"`
		Information struct {
			Author   string `json:"author"`
			URL      string `json:"url"`
			Email    string `json:"email"`
			Twitter  string `json:"twitter"`
			Facebook string `json:"facebook"`
		} `json:"information"`
	}

	err := c.Bind(&body)
	if err != nil {
		return err
	}

	t.Domain = body.Domain
	t.Title = body.Title
	t.Description = body.Description
	t.Active = body.Active
	t.Information.Author = body.Information.Author
	t.Information.URL = body.Information.URL
	t.Information.Email = body.Information.Email
	t.Information.Twitter = body.Information.Twitter
	t.Information.Facebook = body.Information.Facebook

	return nil
}
