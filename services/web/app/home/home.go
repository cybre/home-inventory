package home

import "github.com/labstack/echo/v4"

func Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(200, "Hello, world!")
	}
}
