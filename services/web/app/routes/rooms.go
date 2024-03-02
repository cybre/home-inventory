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

		toast.Success(c, "Room has been updated successfully")

		if htmx.IsHTMXRequest(c) {
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
			"EditingRoom": room,
		})
	}
}