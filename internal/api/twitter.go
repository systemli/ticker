package api

import (
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
)

//PostAuthTwitterHandler returns access token and secret for twitter access.
func PostAuthTwitterHandler(c *gin.Context) {
	token, secret, err := config(c).AccessToken(c.Query("oauth_token"), "", c.Query("oauth_verifier"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"access_token": token, "access_secret": secret})
	return
}

//PostTwitterRequestTokenHandler returns request tokens for twitter login process.
func PostTwitterRequestTokenHandler(c *gin.Context) {
	token, secret, err := config(c).RequestToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"oauth_token": token, "oauth_token_secret": secret})
	return
}

func config(c *gin.Context) *oauth1.Config {
	return &oauth1.Config{
		ConsumerKey:    Config.TwitterConsumerKey,
		ConsumerSecret: Config.TwitterConsumerSecret,
		CallbackURL:    c.Query("callback"),
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://api.twitter.com/oauth/request_token",
			AuthorizeURL:    "https://api.twitter.com/oauth/authenticate",
			AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
		},
	}
}
