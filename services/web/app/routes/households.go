package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cybre/home-inventory/services/inventory/client"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HouseholdCreator interface {
	CreateHousehold(ctx context.Context, household client.CreateHouseholdRequest) error
}

func createHouseholdHandler(householdCreator HouseholdCreator) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}

		if err := householdCreator.CreateHousehold(c.Request().Context(), client.CreateHouseholdRequest{
			UserID:      user.ID,
			HouseholdID: uuid.NewString(),
			Name:        c.FormValue("name"),
			Location:    c.FormValue("location"),
			Description: c.FormValue("description"),
		}); err != nil {
			return fmt.Errorf("failed to create household: %w", err)
		}

		if err := helpers.SessionSet(c, SessionHasHouseholdKey, true); err != nil {
			return fmt.Errorf("failed to set session: %w", err)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
