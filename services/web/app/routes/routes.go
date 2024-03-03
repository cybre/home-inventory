package routes

import (
	"net/http"

	"github.com/cybre/home-inventory/internal/authenticator"
	inventoryclient "github.com/cybre/home-inventory/services/inventory/client"
	"github.com/cybre/home-inventory/services/web/app/auth"
	"github.com/labstack/echo/v4"
)

func Initialize(e *echo.Echo, authenticator *authenticator.Authenticator, inventoryClient *inventoryclient.InventoryClient) {
	e.GET("/", homeHandler(), mustHaveHousehold(inventoryClient))
	e.GET("/households", homeHandler(), mustHaveHousehold(inventoryClient))
	e.GET("/login", loginHandler(authenticator))
	e.GET("/callback", callbackHandler(authenticator, inventoryClient))
	e.GET("/logout", logoutHandler(), auth.IsAuthenticated)
	e.GET("/postlogout", postLogoutHandler())

	e.GET("/onboarding", func(c echo.Context) error {
		return c.Render(http.StatusOK, "onboarding", map[string]interface{}{"Title": "Onboarding"})
	}, auth.IsAuthenticated, mustNotHaveHousehold(inventoryClient))

	e.GET("/onboarding/create-household", func(c echo.Context) error {
		return c.Render(http.StatusOK, "onboarding_create_household", map[string]interface{}{"Title": "Onboarding"})
	}, auth.IsAuthenticated, mustNotHaveHousehold(inventoryClient))

	e.GET("/households/create", createHouseholdViewHandler(), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.POST("/households/create", createHouseholdHandler(inventoryClient), auth.IsAuthenticated)
	e.GET("/households/:householdId", getHouseholdHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId/edit", editHouseholdViewHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.POST("/households/:householdId/edit", editHouseholdHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId/delete", deleteHouseholdViewHandler(), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.POST("/households/:householdId/delete", deleteHouseholdHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))

	e.GET("/households/:householdId/rooms/:roomId", getRoomHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId/rooms/create", createRoomViewHandler(), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.POST("/households/:householdId/rooms/create", createRoomHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId/rooms/:roomId/edit", editRoomViewHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.POST("/households/:householdId/rooms/:roomId/edit", editRoomHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId/rooms/:roomId/delete", deleteRoomViewHandler(), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.POST("/households/:householdId/rooms/:roomId/delete", deleteRoomHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
}
