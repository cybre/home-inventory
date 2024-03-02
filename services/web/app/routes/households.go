package routes

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cybre/home-inventory/services/inventory/client"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/cybre/home-inventory/services/web/app/htmx"
	"github.com/cybre/home-inventory/services/web/app/toast"
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

		referer := c.Request().Header.Get("Referer")

		if err := householdCreator.CreateHousehold(c.Request().Context(), client.CreateHouseholdRequest{
			UserID:           user.ID,
			HouseholdID:      uuid.NewString(),
			Name:             c.FormValue("name"),
			Location:         c.FormValue("location"),
			Description:      c.FormValue("description"),
			IsFromOnboarding: strings.HasSuffix(referer, "/onboarding/create-household"),
		}); err != nil {
			return fmt.Errorf("failed to create household: %w", err)
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

type HouseholdGetter interface {
	GetUserHousehold(ctx context.Context, userID, householdID string) (shared.UserHousehold, error)
}

func getHouseholdHandler(householdGetter HouseholdGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")
		household, err := householdGetter.GetUserHousehold(c.Request().Context(), user.ID, householdID)
		if err != nil {
			return fmt.Errorf("failed to get household: %w", err)
		}

		if htmx.IsHTMXRequest(c) {
			return c.Render(http.StatusOK, "household_card", map[string]interface{}{
				"Household": household,
			})
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

type HouseholdUpdater interface {
	HouseholdGetter
	UpdateHousehold(ctx context.Context, household client.UpdateHouseholdRequest) error
}

func editHouseholdHandler(householdUpdater HouseholdUpdater) echo.HandlerFunc {
	return func(c echo.Context) error {
		// return toast.Error("Failed to update household")

		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")

		household, err := householdUpdater.GetUserHousehold(c.Request().Context(), user.ID, householdID)
		if err != nil {
			return fmt.Errorf("failed to get household: %w", err)
		}

		if err := householdUpdater.UpdateHousehold(c.Request().Context(), client.UpdateHouseholdRequest{
			UserID:      user.ID,
			HouseholdID: householdID,
			Name:        c.FormValue("name"),
			Location:    c.FormValue("location"),
			Description: c.FormValue("description"),
		}); err != nil {
			return toast.Error("Failed to update household")
		}

		household.Name = c.FormValue("name")
		household.Location = c.FormValue("location")
		household.Description = c.FormValue("description")

		toast.Success(c, "Household has been updated successfully")

		if htmx.IsHTMXRequest(c) {
			return c.Render(http.StatusOK, "household_card", map[string]interface{}{
				"Household": household,
			})
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

func editHouseholdViewHandler(householdGetter HouseholdGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")
		household, err := householdGetter.GetUserHousehold(c.Request().Context(), user.ID, householdID)
		if err != nil {
			return fmt.Errorf("failed to get household: %w", err)
		}

		if htmx.IsHTMXRequest(c) {
			return c.Render(http.StatusOK, "household_edit", household)
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":   "Edit Household",
			"Editing": household,
		})
	}
}
