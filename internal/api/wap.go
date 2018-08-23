package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/storage"
	. "git.codecoop.org/systemli/ticker/internal/util"
	. "git.codecoop.org/systemli/ticker/internal/model"
)

//
func GetWAPHandler(c *gin.Context) {
	domain, err := GetDomain(c)
	if err != nil {

		return
	}

	c.Header("Content-Type", "text/vnd.wap.wml")

	ticker, err := FindTicker(domain)
	if err != nil || !ticker.Active {
		settings := GetInactiveSettings().Value.(*InactiveSettings)

		c.HTML(http.StatusOK, "wap.tmpl", gin.H{
			"Author":   settings.Author,
			"Title":    settings.Headline,
			"Messages": []interface{}{},
		})

		return
	}

	pagination := NewPagination(c)
	messages, err := FindByTicker(ticker, pagination)
	if err != nil {

		return
	}

	c.HTML(http.StatusOK, "wap.tmpl", gin.H{
		"Author":   ticker.Information.Author,
		"Title":    ticker.Title,
		"Messages": messages,
	})
}
