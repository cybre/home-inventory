package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func loginHandler(auth *authenticator.Authenticator) echo.HandlerFunc {
	return func(c echo.Context) error {
		state, err := generateRandomState()
		if err != nil {
			return fmt.Errorf("failed to generate random state: %w", err)
		}

		sess, err := session.Get(AuthSessionCookieName, c)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}

		sess.Values["state"] = state

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
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