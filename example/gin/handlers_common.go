//go:build gin

package main

import (
	"net/http"

	ginlib "github.com/gin-gonic/gin"
)

// handleHealthz is a small, taggable endpoint that doesn't use groups.
func handleHealthz(c *ginlib.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
