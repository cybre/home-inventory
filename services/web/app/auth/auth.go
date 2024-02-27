package auth

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	AuthSessionCookieName     = "auth-session"
	AuthSessionProfileKey     = "profile"
	AuthSessionAccessTokenKey = "access_token"
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

func init() {
	gob.Register(User{})
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(AuthSessionCookieName, c)
		if err != nil {
			return fmt.Errorf("failed to get auth session: %w", err)
		}

		if sess.Values[AuthSessionProfileKey] == nil {
			return c.Redirect(http.StatusSeeOther, "/login?redirectTo="+url.QueryEscape(c.Request().URL.Path))
		}

		return next(c)
	}
}
