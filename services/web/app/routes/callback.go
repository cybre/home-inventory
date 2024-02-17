package routes

import (
	"fmt"
	"net/http"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID             string
	FirstName      string
	LastName       string
	ProfilePicture string
}

func NewUserFromProfile(profile map[string]interface{}) *User {
	return &User{
		ID:             profile["sub"].(string),
		FirstName:      profile["given_name"].(string),
		LastName:       profile["family_name"].(string),
		ProfilePicture: profile["picture"].(string),
	}
}

func callbackHandler(auth *authenticator.Authenticator) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := session.Get(AuthSessionCookieName, c)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}

		state := session.Values["state"].(string)
		if state != c.QueryParam("state") {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid state parameter")
		}

		code := c.QueryParam("code")
		token, err := auth.Exchange(c.Request().Context(), code)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to exchange code for token")
		}

		idToken, err := auth.VerifyIDToken(c.Request().Context(), token)
		if err != nil {
			return fmt.Errorf("failed to verify ID token: %w", err)
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			return fmt.Errorf("failed to parse ID token claims: %w", err)
		}

		session.Values[AuthSessionAccessTokenKey] = token.AccessToken
		session.Values[AuthSessionProfileKey] = NewUserFromProfile(profile)
		if err := session.Save(c.Request(), c.Response()); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		redirectTo, ok := redirectMap[state]
		if !ok || redirectTo == "" {
			redirectTo = "/"
		}
		delete(redirectMap, state)

		return c.Redirect(http.StatusTemporaryRedirect, redirectTo)
	}
}
