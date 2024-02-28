package routes

import (
	"context"
	"net/http"

	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/labstack/echo/v4"
)

const (
	ContextHouseholdsKey = "households"
)

func mustHaveHousehold(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		households, ok := helpers.ContextGet[[]shared.UserHousehold](c, ContextHouseholdsKey)
		if !ok {
			return next(c)
		}

		if len(households) == 0 {
			return c.Redirect(http.StatusTemporaryRedirect, "/onboarding")
		}

		return next(c)
	}
}

func mustNotHaveHousehold(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		households, ok := helpers.ContextGet[[]shared.UserHousehold](c, ContextHouseholdsKey)
		if !ok {
			return next(c)
		}

		if len(households) > 0 {
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		return next(c)
	}
}

type HouseholdsGetter interface {
	GetUserHouseholds(ctx context.Context, userID string) ([]shared.UserHousehold, error)
}

func LoadHouseholdsIntoContext(householdsGetter HouseholdsGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := helpers.GetUser(c)
			if !ok {
				c.Set(ContextHouseholdsKey, []shared.UserHousehold{})
				return next(c)
			}

			households, err := householdsGetter.GetUserHouseholds(c.Request().Context(), user.ID)
			if err != nil {
				c.Set(ContextHouseholdsKey, []shared.UserHousehold{})
				return next(c)
			}

			c.Set(ContextHouseholdsKey, households)
			return next(c)
		}
	}
}
