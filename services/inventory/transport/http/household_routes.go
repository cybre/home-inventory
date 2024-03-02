package http

import (
	"net/http"

	eh "github.com/cybre/home-inventory/internal/handler"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func buildHouseholdRoutes(e *echo.Echo, householdService HouseholdService, validate *validator.Validate) {
	e.GET(shared.UserHouseholdsRoute, getUserHouseholdsHandler(householdService))
	e.POST(shared.UserHouseholdsRoute, eh.NewValidateHandler(createHouseholdHandler(householdService), validate))
	e.GET(shared.UserHouseholdRoute, getUserHouseholdHandler(householdService))
	e.PUT(shared.UserHouseholdRoute, eh.NewValidateHandler(updateHouseholdHandler(householdService), validate))

	e.POST(shared.UserHouseholdRoomsRoute, eh.NewValidateHandler(addRoomHandler(householdService), validate))
	e.GET(shared.UserHouseholdRoomRoute, getUserHouseholdRoomHandler(householdService))
	e.PUT(shared.UserHouseholdRoomRoute, eh.NewValidateHandler(updateRoomHandler(householdService), validate))
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

func getUserHouseholdsHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")

		households, err := householdService.GetUserHouseholds(c.Request().Context(), userId)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, households)
	}
}

func getUserHouseholdHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		householdId := c.Param("householdId")

		household, err := householdService.GetUserHousehold(c.Request().Context(), userId, householdId)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, household)
	}
}

func updateHouseholdHandler(householdService HouseholdService) eh.Handler[shared.UpdateHouseholdCommandData] {
	return func(c echo.Context, data shared.UpdateHouseholdCommandData) error {
		if err := householdService.UpdateHousehold(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func getUserHouseholdRoomHandler(householdService HouseholdService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		householdId := c.Param("householdId")
		roomId := c.Param("roomId")

		household, err := householdService.GetUserHouseholdRoom(c.Request().Context(), userId, householdId, roomId)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, household)
	}
}

func addRoomHandler(householdService HouseholdService) eh.Handler[shared.AddRoomCommandData] {
	return func(c echo.Context, data shared.AddRoomCommandData) error {
		if err := householdService.AddRoom(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func updateRoomHandler(householdService HouseholdService) eh.Handler[shared.UpdateRoomCommandData] {
	return func(c echo.Context, data shared.UpdateRoomCommandData) error {
		if err := householdService.UpdateRoom(c.Request().Context(), data); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
