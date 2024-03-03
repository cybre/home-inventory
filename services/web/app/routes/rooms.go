package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cybre/home-inventory/services/inventory/client"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/cybre/home-inventory/services/web/app/helpers"
	"github.com/cybre/home-inventory/services/web/app/htmx"
	"github.com/cybre/home-inventory/services/web/app/toast"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoomGetter interface {
	GetUserHouseholdRoom(ctx context.Context, userID, householdID, roomID string) (shared.UserHouseholdRoom, error)
}

func getRoomHandler(roomGetter RoomGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")
		roomID := c.Param("roomId")
		room, err := roomGetter.GetUserHouseholdRoom(c.Request().Context(), user.ID, householdID, roomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}

		if htmx.IsHTMXRequest(c) {
			return c.Render(http.StatusOK, "room_card", room)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

type RoomCreator interface {
	AddRoom(ctx context.Context, room client.AddRoomRequest) error
}

func createRoomHandler(roomCreator RoomCreator) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")

		request := client.AddRoomRequest{
			HouseholdID: householdID,
			UserID:      user.ID,
			RoomID:      uuid.NewString(),
			Name:        c.FormValue("name"),
		}

		if err := roomCreator.AddRoom(c.Request().Context(), request); err != nil {
			return toast.Error("Failed to add room")
		}

		if htmx.IsHTMXRequest(c) {
			toast.Success(c, "Room has added successfully")
			return c.Render(http.StatusOK, "room_card", shared.UserHouseholdRoom{
				HouseholdID: request.HouseholdID,
				RoomID:      request.RoomID,
				Name:        request.Name,
			})
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

func createRoomViewHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		householdId := c.Param("householdId")

		if htmx.IsHTMXRequest(c) {
			return c.Render(http.StatusOK, "room_add", householdId)
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":      "Add Room",
			"AddingRoom": householdId,
		})
	}
}

type RoomUpdater interface {
	RoomGetter
	UpdateRoom(ctx context.Context, room client.UpdateRoomRequest) error
}

func editRoomHandler(roomUpdater RoomUpdater) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")
		roomID := c.Param("roomId")

		room, err := roomUpdater.GetUserHouseholdRoom(c.Request().Context(), user.ID, householdID, roomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}

		if err := roomUpdater.UpdateRoom(c.Request().Context(), client.UpdateRoomRequest{
			HouseholdID: householdID,
			UserID:      user.ID,
			RoomID:      roomID,
			Name:        c.FormValue("name"),
		}); err != nil {
			return toast.Error("Failed to update room")
		}

		room.Name = c.FormValue("name")

		if htmx.IsHTMXRequest(c) {
			toast.Success(c, "Room has been updated successfully")
			return c.Render(http.StatusOK, "room_card", room)
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

func editRoomViewHandler(roomGetter RoomGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdId := c.Param("householdId")
		roomID := c.Param("roomId")

		room, err := roomGetter.GetUserHouseholdRoom(c.Request().Context(), user.ID, householdId, roomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}

		if htmx.IsHTMXRequest(c) {
			return c.Render(http.StatusOK, "room_edit", room)
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":       "Edit Room",
			"EditingRoom": room.RoomID,
		})
	}
}

type RoomDeleter interface {
	DeleteRoom(ctx context.Context, userId, householdId, roomId string) error
}

func deleteRoomHandler(roomDeleter RoomDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := helpers.GetUser(c)
		if !ok {
			return fmt.Errorf("user not found")
		}

		householdID := c.Param("householdId")
		roomID := c.Param("roomId")

		if err := roomDeleter.DeleteRoom(c.Request().Context(), user.ID, householdID, roomID); err != nil {
			return toast.Error("Failed to delete room")
		}

		if htmx.IsHTMXRequest(c) {
			toast.Success(c, "Room has been deleted successfully")
			return c.NoContent(http.StatusOK)
		}

		return c.Redirect(http.StatusFound, "/")
	}
}
