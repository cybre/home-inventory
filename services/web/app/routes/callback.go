package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/cybre/home-inventory/services/web/app/auth"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type UserHouseholdGetter interface {
	GetUserHouseholds(ctx context.Context, userID string) ([]shared.UserHousehold, error)
}

func callbackHandler(authenticator *authenticator.Authenticator, userHouseholdGetter UserHouseholdGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := session.Get(auth.AuthSessionCookieName, c)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}

		state := session.Values["state"].(string)
		if state != c.QueryParam("state") {
			return fmt.Errorf("state mismatch")
		}

		code := c.QueryParam("code")
		token, err := authenticator.Exchange(c.Request().Context(), code)
		if err != nil {
			return fmt.Errorf("failed to exchange code for token: %w", err)
		}

		idToken, err := authenticator.VerifyIDToken(c.Request().Context(), token)
		if err != nil {
			return fmt.Errorf("failed to verify ID token: %w", err)
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			return fmt.Errorf("failed to parse ID token claims: %w", err)
		}

		session.Values[auth.AuthSessionAccessTokenKey] = token.AccessToken
		session.Values[auth.AuthSessionProfileKey] = auth.NewUserFromProfile(profile)

		redirectTo, ok := session.Values["redirectTo"].(string)
		if !ok || redirectTo == "" {
			redirectTo = "/"
		}
		delete(session.Values, "redirectTo")

		if err := session.Save(c.Request(), c.Response()); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
	}
}
