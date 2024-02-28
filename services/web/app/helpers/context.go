package helpers

import "github.com/labstack/echo/v4"

func ContextGet[T any](c echo.Context, key string) (T, bool) {
	value, ok := c.Get(key).(T)

	return value, ok
}
