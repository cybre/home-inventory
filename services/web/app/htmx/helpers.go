package htmx

import "github.com/labstack/echo/v4"

func IsHTMXRequest(e echo.Context) bool {
	return e.Request().Header.Get("HX-Request") == "true"
}
