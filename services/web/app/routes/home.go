package routes

import (
	"net/http"

	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/labstack/echo/v4"
)

func homeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !helpers.IsAuthenticated(c) {
			return c.Render(http.StatusOK, "login", map[string]interface{}{"Title": "Home Inventory"})
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{"Title": "Welcome to Home Inventory"})
	}
}
