package postlogout

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:   "auth-session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
