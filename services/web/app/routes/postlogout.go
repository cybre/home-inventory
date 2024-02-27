package routes

import (
	"net/http"

	"github.com/cybre/home-inventory/services/web/app/auth"
	"github.com/labstack/echo/v4"
)

func postLogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:   auth.AuthSessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
