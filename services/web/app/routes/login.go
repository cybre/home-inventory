package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/labstack/echo/v4"
)

func loginHandler(authenticator *authenticator.Authenticator) echo.HandlerFunc {
	return func(c echo.Context) error {
		state, err := generateRandomState()
		if err != nil {
			return fmt.Errorf("failed to generate random state: %w", err)
		}

		if err := helpers.SessionSet(c, "state", state, "redirectTo", c.QueryParam("redirectTo")); err != nil {
			return fmt.Errorf("failed to save state to session: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, authenticator.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
