package http

import (
	"net/http"

	eh "github.com/cybre/home-inventory/internal/echohandler"
	"github.com/cybre/home-inventory/inventory/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func buildHouseholdRoutes(e *echo.Echo, householdService HouseholdService, validate *validator.Validate) {
	e.POST("/household", eh.NewValidateHandler(createHouseholdHandler(householdService), validate))
	e.POST("/household/:householdId/rooms", eh.NewValidateHandler(addHouseholdRoomHandler(householdService), validate))
	e.POST("/household/:householdId/rooms/:roomID/items", eh.NewValidateHandler(addItemHandler(householdService), validate))
	e.POST("/household/:householdId/rooms/:roomID/items/:itemID", eh.NewValidateHandler(updateItemHandler(householdService), validate))
}

func createHouseholdHandler(householdService HouseholdService) eh.Handler[shared.CreateHouseholdCommandData] {
	return func(c echo.Context, data shared.CreateHouseholdCommandData) error {
		if err := householdService.CreateHousehold(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}

func addHouseholdRoomHandler(householdService HouseholdService) eh.Handler[shared.AddRoomCommandData] {
	return func(c echo.Context, data shared.AddRoomCommandData) error {
		if err := householdService.AddRoom(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}

func addItemHandler(householdService HouseholdService) eh.Handler[shared.AddItemCommandData] {
	return func(c echo.Context, data shared.AddItemCommandData) error {
		if err := householdService.AddItem(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}

func updateItemHandler(householdService HouseholdService) eh.Handler[shared.UpdateItemCommandData] {
	return func(c echo.Context, data shared.UpdateItemCommandData) error {
		if err := householdService.UpdateItem(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
