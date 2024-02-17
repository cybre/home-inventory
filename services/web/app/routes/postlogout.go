package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func postLogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:   AuthSessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
