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
			return err
		}

		if htmx.ShouldReturnPartial(c) {
			htmx.ReplaceUrl(c, "/")
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
			return err
		}

		if htmx.ShouldReturnPartial(c) {
			toast.Success(c, "Room has added successfully")
			htmx.ReplaceUrl(c, "/")
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

		if htmx.ShouldReturnPartial(c) {
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
			return err
		}

		if err := roomUpdater.UpdateRoom(c.Request().Context(), client.UpdateRoomRequest{
			HouseholdID: householdID,
			UserID:      user.ID,
			RoomID:      roomID,
			Name:        c.FormValue("name"),
		}); err != nil {
			return err
		}

		room.Name = c.FormValue("name")

		if htmx.ShouldReturnPartial(c) {
			toast.Success(c, "Room has been updated successfully")
			htmx.ReplaceUrl(c, "/")
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
			return err
		}

		if htmx.ShouldReturnPartial(c) {
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

		deleteKey := c.FormValue("delete")
		if deleteKey != "delete" {
			if htmx.IsHTMXRequest(c) {
				return toast.Error("Please enter the text exactly as shown to confirm")
			}

			return c.Redirect(http.StatusFound, "/households/"+householdID+"/rooms/"+roomID+"/delete?invalidKey=true")
		}

		if err := roomDeleter.DeleteRoom(c.Request().Context(), user.ID, householdID, roomID); err != nil {
			return err
		}

		if htmx.ShouldReturnPartial(c) {
			toast.Success(c, "Room has been deleted successfully")
			htmx.ReplaceUrl(c, "/")
			return c.NoContent(http.StatusOK)
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

func deleteRoomViewHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		householdID := c.Param("householdId")
		roomID := c.Param("roomId")

		data := map[string]interface{}{"HouseholdID": householdID, "RoomID": roomID}

		if htmx.ShouldReturnPartial(c) {
			return c.Render(http.StatusOK, "room_confirm_delete", data)
		}

		if c.QueryParam("invalidKey") == "true" {
			data["InvalidDeleteKey"] = true
		}

		return c.Render(http.StatusOK, "home", map[string]interface{}{
			"Title":        "Delete Room",
			"EditingRoom":  roomID,
			"DeletingRoom": data,
		})
	}
}
