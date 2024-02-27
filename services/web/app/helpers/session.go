package helpers

import (
	"fmt"

	"github.com/cybre/home-inventory/services/web/app/auth"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(c echo.Context) bool {
	sess, err := session.Get(auth.AuthSessionCookieName, c)
	if err != nil {
		return false
	}

	if sess.Values[auth.AuthSessionProfileKey] != nil {
		return true
	}

	return false
}

func GetUser(c echo.Context) (auth.User, bool) {
	return SessionGet[auth.User](c, auth.AuthSessionProfileKey)
}

func SessionSet(c echo.Context, args ...interface{}) error {
	sess, err := session.Get(auth.AuthSessionCookieName, c)
	if err != nil {
		return fmt.Errorf("failed to get auth session: %w", err)
	}

	for len(args) >= 2 {
		key := args[0]
		value := args[1]
		args = args[2:]
		sess.Values[key] = value
	}

	return sess.Save(c.Request(), c.Response())
}

func SessionGet[T any](c echo.Context, key string) (T, bool) {
	sess, err := session.Get(auth.AuthSessionCookieName, c)
	if err != nil {
		return *new(T), false
	}

	value, ok := sess.Values[key].(T)

	return value, ok
}
