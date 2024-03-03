package htmx

import "github.com/labstack/echo/v4"

func IsHTMXRequest(e echo.Context) bool {
	return e.Request().Header.Get("HX-Request") == "true"
}

func ReplaceUrl(e echo.Context, url string) {
	e.Response().Header().Add("HX-Replace-Url", url)
}

func IsHTMXHistoryRestoreRequest(e echo.Context) bool {
	return e.Request().Header.Get("HX-History-Restore-Request") == "true"
}

func ShouldReturnPartial(e echo.Context) bool {
	return IsHTMXRequest(e) && !IsHTMXHistoryRestoreRequest(e)
}
