package api

import (
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
)

func (h *handler) PostAuthTwitter(c *gin.Context) {
	oauthConfig := oauthConfig(c, h.config)
	token, secret, err := oauthConfig.AccessToken(c.Query("oauth_token"), "", c.Query("oauth_verifier"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.Unauthorized))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"access_token": token, "access_secret": secret})
}

func (h *handler) PostTwitterRequestToken(c *gin.Context) {
	oauthConfig := oauthConfig(c, h.config)
	token, secret, err := oauthConfig.RequestToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.Unauthorized))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"oauth_token": token, "oauth_token_secret": secret})
}

func oauthConfig(c *gin.Context, config config.Config) *oauth1.Config {
	return &oauth1.Config{
		ConsumerKey:    config.TwitterConsumerKey,
		ConsumerSecret: config.TwitterConsumerSecret,
		CallbackURL:    c.Query("callback"),
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://api.twitter.com/oauth/request_token",
			AuthorizeURL:    "https://api.twitter.com/oauth/authenticate",
			AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
		},
	}
}
