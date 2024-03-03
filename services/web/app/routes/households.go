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

		request := client.CreateHouseholdRequest{
			UserID:           user.ID,
			HouseholdID:      uuid.NewString(),
			Name:             c.FormValue("name"),
			Location:         c.FormValue("location"),
			Description:      c.FormValue("description"),
			IsFromOnboarding: strings.HasSuffix(referer, "/onboarding/create-household"),
		}

		if err := householdCreator.CreateHousehold(c.Request().Context(), request); err != nil {
			return fmt.Errorf("failed to create household: %w", err)
		}

		if htmx.ShouldReturnPartial(c) {
			toast.Success(c, "Household has been created successfully")
			htmx.ReplaceUrl(c, "/")
			return c.Render(http.StatusOK, "household_card", map[string]interface{}{
				"Household": shared.UserHousehold{
					UserID:      request.UserID,
					HouseholdID: request.HouseholdID,
					Name:        request.Name,
					Location:    request.Location,
					Description: request.Description,
					Rooms:       []shared.UserHouseholdRoom{},
					Timestamp:   0,
					Order:       0,
				},
			})
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

func createHouseholdViewHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if htmx.ShouldReturnPartial(c) {
			return c.Render(http.StatusOK, "household_create", nil)
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":             "Create Household",
			"CreatingHousehold": true,
		})
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

		if htmx.ShouldReturnPartial(c) {
			htmx.ReplaceUrl(c, "/")
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

		if htmx.ShouldReturnPartial(c) {
			toast.Success(c, "Household has been updated successfully")
			htmx.ReplaceUrl(c, "/")
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

		if htmx.ShouldReturnPartial(c) {
			return c.Render(http.StatusOK, "household_edit", household)
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":            "Edit Household",
			"EditingHousehold": household.HouseholdID,
		})
	}
}

func deleteHouseholdViewHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		householdID := c.Param("householdId")

		data := map[string]interface{}{"HouseholdID": householdID}

		if htmx.ShouldReturnPartial(c) {
			return c.Render(http.StatusOK, "household_confirm_delete", data)
		}

		if c.QueryParam("invalidKey") == "true" {
			data["InvalidDeleteKey"] = true
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":             "Delete Household",
			"EditingHousehold":  householdID,
			"DeletingHousehold": data,
		})
	}
}

type HouseholdDeleter interface {
	DeleteHousehold(ctx context.Context, userID, householdID string) error
}

func deleteHouseholdHandler(householdDeleter HouseholdDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		deleteKey := c.FormValue("delete")
		if deleteKey != "delete" {
			if htmx.IsHTMXRequest(c) {
				return toast.Error("Please enter the text exactly as shown to confirm")
			}

			return c.Redirect(http.StatusFound, "/households/"+c.Param("householdId")+"/delete?invalidKey=true")
		}

		householdID := c.Param("householdId")

		if err := householdDeleter.DeleteHousehold(c.Request().Context(), user.ID, householdID); err != nil {
			return fmt.Errorf("failed to delete household: %w", err)
		}

		if htmx.ShouldReturnPartial(c) {
			toast.Success(c, "Household has been deleted successfully")
			htmx.ReplaceUrl(c, "/")
			return c.NoContent(http.StatusOK)
		}

		return c.Redirect(http.StatusFound, "/")
	}
}
