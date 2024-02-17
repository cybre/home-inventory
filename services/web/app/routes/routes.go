package routes

import (
	"net/http"
	"net/url"

	"github.com/cybre/home-inventory/internal/authenticator"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	AuthSessionCookieName     = "auth-session"
	AuthSessionProfileKey     = "profile"
	AuthSessionAccessTokenKey = "access_token"
)

func Initialize(e *echo.Echo, auth *authenticator.Authenticator) {
	e.GET("/login", loginHandler(auth))
	e.GET("/callback", callbackHandler(auth))
	e.GET("/logout", logoutHandler())
	e.GET("/postlogout", postLogoutHandler())
	e.GET("/", homeHandler())
	e.GET("/protected", func(c echo.Context) error {
		return c.String(http.StatusOK, "protected")
	}, isAuthenticated)
}

func isAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(AuthSessionCookieName, c)
		if err != nil {
			return err
		}

		if sess.Values[AuthSessionProfileKey] == nil {
			return c.Redirect(http.StatusSeeOther, "/login?redirectTo="+url.QueryEscape(c.Request().URL.Path))
		}

		return next(c)
	}
}
