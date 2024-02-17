package routes

import "github.com/labstack/echo/v4"

func homeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(200, "Hello, world!")
	}
}
