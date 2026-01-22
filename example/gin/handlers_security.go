//go:build gin && security && !typed

package main

import (
	"net/http"
	"strings"

	"github.com/aizacoders/openapigo/openapi"
	ginlib "github.com/gin-gonic/gin"
)

func requireBearer(c *ginlib.Context) bool {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.Status(http.StatusUnauthorized)
		return false
	}
	return true
}

func requireAPIKey(c *ginlib.Context) bool {
	if c.GetHeader("X-API-Key") == "" {
		c.Status(http.StatusUnauthorized)
		return false
	}
	return true
}

func handleSecureHealthz(c *ginlib.Context) {
	if !requireBearer(c) {
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func handleSecureListUsers(c *ginlib.Context) {
	if !requireBearer(c) {
		return
	}
	c.JSON(http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
}

func handleSecureCreateUser(c *ginlib.Context) {
	if !requireAPIKey(c) {
		return
	}
	c.Status(http.StatusCreated)
}

func handleSecureDemoErrors(c *ginlib.Context) {
	if !requireBearer(c) {
		c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{Error: "unauthorized"})
		return
	}
	switch c.Query("code") {
	case "400":
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Error: "bad request"})
		return
	case "500":
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Error: "internal error"})
		return
	case "503":
		c.JSON(http.StatusServiceUnavailable, openapi.ErrorResponse{Error: "service unavailable"})
		return
	default:
		c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}
