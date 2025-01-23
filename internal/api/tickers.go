package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-mastodon"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/bluesky"
	"github.com/systemli/ticker/internal/signal"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetTickers(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	filter := storage.NewTickerFilter(c.Request)
	tickers, err := h.storage.FindTickersByUser(me, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"tickers": response.TickersResponse(tickers, h.config)}))
}

func (h *handler) GetTicker(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) GetTickerUsers(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	users, _ := h.storage.FindUsersByTicker(ticker)

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"users": response.UsersResponse(users)}))
}

func (h *handler) PostTicker(c *gin.Context) {
	ticker := storage.NewTicker()
	err := updateTicker(&ticker, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTicker(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	err = updateTicker(&ticker, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.bridges.Update(ticker)
	if err != nil {
		log.WithError(err).Error("failed to update ticker in bridges")
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.BridgeError))
		return
	}

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	h.ClearTickerCache(&ticker)

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTickerUsers(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Users []storage.User `json:"users" binding:"required"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	ticker.Users = body.Users

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"users": response.UsersResponse(ticker.Users)}))
}

func (h *handler) PostTickerWebsite(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Origin string `json:"origin" binding:"required,http_url"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.storage.SaveTickerWebsite(&ticker, body.Origin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	h.ClearTickerCache(&ticker)

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) DeleteTickerWebsite(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Origin string `json:"origin" binding:"required"`
	}
	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	err = h.storage.DeleteTickerWebsite(&ticker, body.Origin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	h.ClearTickerCache(&ticker)

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTickerTelegram(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body storage.TickerTelegram
	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeNotFound, response.FormError))
		return
	}

	ticker.Telegram.Active = body.Active
	if body.ChannelName != "" {
		ticker.Telegram.ChannelName = body.ChannelName
	}

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) DeleteTickerTelegram(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	err = h.storage.DeleteTelegram(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTickerMastodon(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body storage.TickerMastodon
	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeNotFound, response.FormError))
		return
	}

	if body.Secret != "" || body.Token != "" || body.AccessToken != "" || body.Server != "" {
		client := mastodon.NewClient(&mastodon.Config{
			Server:       body.Server,
			ClientID:     body.Token,
			ClientSecret: body.Secret,
			AccessToken:  body.AccessToken,
		})

		account, err := client.GetAccountCurrentUser(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeBadCredentials, response.MastodonError))
			return
		}

		ticker.Mastodon.Server = body.Server
		ticker.Mastodon.Secret = body.Secret
		ticker.Mastodon.Token = body.Token
		ticker.Mastodon.AccessToken = body.AccessToken
		ticker.Mastodon.User = storage.MastodonUser{
			Username:    account.Username,
			Avatar:      account.Avatar,
			DisplayName: account.DisplayName,
		}
	}

	ticker.Mastodon.Active = body.Active

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) DeleteTickerMastodon(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	err = h.storage.DeleteMastodon(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTickerBluesky(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body storage.TickerBluesky
	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeNotFound, response.FormError))
		return
	}

	if body.Handle != "" && body.AppKey != "" {
		_, err = bluesky.Authenticate(body.Handle, body.AppKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeBadCredentials, response.BlueskyError))
			return
		}

		ticker.Bluesky.Handle = body.Handle
		ticker.Bluesky.AppKey = body.AppKey
	}
	ticker.Bluesky.Active = body.Active

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) DeleteTickerBluesky(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	err = h.storage.DeleteBluesky(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTickerSignalGroup(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body storage.TickerSignalGroup
	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeNotFound, response.FormError))
		return
	}

	groupClient := signal.NewGroupClient(h.config)
	err = groupClient.CreateOrUpdateGroup(&ticker)
	if err != nil {
		log.WithError(err).Error("failed to create or update group")
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.SignalGroupError))
		return
	}
	ticker.SignalGroup.Active = body.Active

	err = h.storage.SaveTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) DeleteTickerSignalGroup(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	groupClient := signal.NewGroupClient(h.config)

	// Remove all members except the account number
	err = groupClient.RemoveAllMembers(ticker.SignalGroup.GroupID)
	if err != nil {
		log.WithError(err).Error("failed to remove members")
		return
	}

	// Quit the group
	err = groupClient.QuitGroup(ticker.SignalGroup.GroupID)
	if err != nil {
		log.WithError(err).Error("failed to quit group")
		return
	}

	err = h.storage.DeleteSignalGroup(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

func (h *handler) PutTickerSignalGroupAdmin(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	var body struct {
		Number string `json:"number" binding:"required"`
	}

	err = c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeNotFound, response.FormError))
		return
	}

	groupClient := signal.NewGroupClient(h.config)
	err = groupClient.AddAdminMember(ticker.SignalGroup.GroupID, body.Number)
	if err != nil {
		log.WithError(err).Error("failed to add member")
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.SignalGroupError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{}))
}

func (h *handler) DeleteTicker(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	err = h.storage.DeleteTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	h.ClearTickerCache(&ticker)

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{}))
}

func (h *handler) DeleteTickerUser(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.UserIdentifierMissing))
		return
	}

	user, err := h.storage.FindUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	err = h.storage.DeleteTickerUser(&ticker, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"users": response.UsersResponse(ticker.Users)}))
}

func (h *handler) ResetTicker(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	err = h.storage.ResetTicker(&ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.StorageError))
		return
	}

	h.ClearTickerCache(&ticker)

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config)}))
}

// ClearTickerCache clears the cache for the init endpoint of a ticker
func (h *handler) ClearTickerCache(ticker *storage.Ticker) {
	h.cache.Range(func(key, value interface{}) bool {
		for _, website := range ticker.Websites {
			if strings.HasPrefix(key.(string), fmt.Sprintf("response:%s:/v1/init", website.Origin)) {
				h.cache.Delete(key)
			}
		}
		return true
	})
}

func updateTicker(t *storage.Ticker, c *gin.Context) error {
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Active      bool   `json:"active"`
		Information struct {
			Author   string `json:"author"`
			URL      string `json:"url"`
			Email    string `json:"email"`
			Twitter  string `json:"twitter"`
			Facebook string `json:"facebook"`
			Telegram string `json:"telegram"`
			Mastodon string `json:"mastodon"`
			Bluesky  string `json:"bluesky"`
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

	t.Title = body.Title
	t.Description = body.Description
	t.Active = body.Active
	t.Information.Author = body.Information.Author
	t.Information.URL = body.Information.URL
	t.Information.Email = body.Information.Email
	t.Information.Twitter = body.Information.Twitter
	t.Information.Facebook = body.Information.Facebook
	t.Information.Telegram = body.Information.Telegram
	t.Information.Mastodon = body.Information.Mastodon
	t.Information.Bluesky = body.Information.Bluesky
	t.Location.Lat = body.Location.Lat
	t.Location.Lon = body.Location.Lon

	return nil
}
