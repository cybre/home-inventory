package routes

import (
	"net/http"

	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/labstack/echo/v4"
)

const (
	SessionHasHouseholdKey = "has_household"
)

func mustHaveHousehold(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		hasHousehold, ok := helpers.SessionGet[bool](c, SessionHasHouseholdKey)
		if !ok {
			return next(c)
		}

		if !hasHousehold {
			return c.Redirect(http.StatusTemporaryRedirect, "/onboarding")
		}

		return next(c)
	}
}

func mustNotHaveHousehold(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		hasHousehold, ok := helpers.SessionGet[bool](c, SessionHasHouseholdKey)
		if !ok {
			return next(c)
		}

		if hasHousehold {
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		return next(c)
	}
}
