package web

import (
	"net/http"
	"strings"

	"github.com/indece-official/go-gousu/gousuchi/v2"
)

func (c *Controller) checkAuth(r *http.Request) gousuchi.IResponse {
	if *disableAuth {
		return nil
	}

	authHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return gousuchi.Unauthorized(r, "Missing authorization bearer token")
	}

	authToken := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
	if authToken == "" {
		return gousuchi.Unauthorized(r, "Empty authorization bearer token")
	}

	err := c.sessionService.VerifySessionToken(authToken)
	if err != nil {
		return gousuchi.Unauthorized(r, "Invalid authorization bearer token: %s", err)
	}

	return nil
}
