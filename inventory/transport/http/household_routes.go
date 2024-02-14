package http

import (
	"net/http"

	"github.com/cybre/home-inventory/inventory/shared"
	"github.com/labstack/echo/v4"
)

func buildHouseholdRoutes(e *echo.Echo, householdService HouseholdService) {
	e.POST("/household", createHouseholdHandler(householdService))
	e.POST("/household/:householdId/rooms", addHouseholdRoomHandler(householdService))
	e.POST("/household/:householdId/rooms/:roomID/items", addItemHandler(householdService))
	e.POST("/household/:householdId/rooms/:roomID/items/:itemID", updateItemHandler(householdService))
}

func createHouseholdHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data shared.CreateHouseholdCommandData
		if err := c.Bind(&data); err != nil {
			return err
		}

		if err := householdService.CreateHousehold(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}

func addHouseholdRoomHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data shared.AddRoomCommandData
		if err := c.Bind(&data); err != nil {
			return err
		}

		if err := householdService.AddRoom(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}

func addItemHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data shared.AddItemCommandData
		if err := c.Bind(&data); err != nil {
			return err
		}

		if err := householdService.AddItem(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusCreated)
	}
}

func updateItemHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data shared.UpdateItemCommandData
		if err := c.Bind(&data); err != nil {
			return err
		}

		if err := householdService.UpdateItem(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
