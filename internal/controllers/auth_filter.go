package controllers

import (
	"net/http"

	"github.com/zeusito/gochat/internal/models"
)

func (c *WebSocketServer) authFilter(w http.ResponseWriter, r *http.Request) (*models.MyClaims, bool) {
	auth := r.Header.Get("Authorization")

	if len(auth) == 0 {
		c.log.Errorf("No authorization header found")
		return nil, false
	}

	// Remove the "Bearer " prefix
	auth = auth[7:]

	claims, err := c.joseSvc.Parse(auth)
	if err != nil {
		c.log.Errorf("Auth failed: %v", err)
		return nil, false
	}

	return claims, true
}
