package routes

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/cybre/home-inventory/services/web/app/htmx"
	"github.com/labstack/echo/v4"
)

const (
	ContextHouseholdsKey = "households"
)

type HouseholdsGetter interface {
	GetUserHouseholds(ctx context.Context, userID string) ([]shared.UserHousehold, error)
}

func mustHaveHousehold(householdsGetter HouseholdsGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			loadHouseholdsIntoContext(c, householdsGetter)

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
}

func mustNotHaveHousehold(householdsGetter HouseholdsGetter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			loadHouseholdsIntoContext(c, householdsGetter)

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
}

func loadHouseholdsIntoContext(c echo.Context, householdsGetter HouseholdsGetter) {
	if htmx.IsHTMXRequest(c) {
		return
	}

	user, ok := helpers.GetUser(c)
	if !ok {
		c.Set(ContextHouseholdsKey, nil)
		return
	}

	households, err := householdsGetter.GetUserHouseholds(c.Request().Context(), user.ID)
	if err != nil {
		logging.FromContext(c.Request().Context()).Error("failed to load households into context", slog.Any("error", err))
		c.Set(ContextHouseholdsKey, nil)
		return
	}

	c.Set(ContextHouseholdsKey, households)
}
