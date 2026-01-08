package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
)

// WebAppManifest represents a Progressive Web App manifest
// https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Manifest
type WebAppManifest struct {
	Name        string         `json:"name"`
	ShortName   string         `json:"short_name"`
	Description string         `json:"description"`
	StartURL    string         `json:"start_url"`
	Display     string         `json:"display"`
	Orientation string         `json:"orientation"`
	Scope       string         `json:"scope"`
	Icons       []ManifestIcon `json:"icons"`
}

// ManifestIcon represents an icon entry in the manifest
type ManifestIcon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose,omitempty"`
}

func (h *handler) HandleManifest(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	manifest := WebAppManifest{
		Name:        ticker.Title,
		ShortName:   ticker.Title,
		Description: ticker.Description,
		StartURL:    "/",
		Display:     "fullscreen",
		Orientation: "portrait-primary",
		Scope:       "/",
		Icons:       []ManifestIcon{},
	}

	c.Header("Content-Type", "application/manifest+json")
	c.JSON(http.StatusOK, manifest)
}
