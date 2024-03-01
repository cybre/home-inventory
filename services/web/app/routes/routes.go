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

	e.POST("/households", createHouseholdHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId", getHouseholdHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.GET("/households/:householdId/edit", editHouseholdViewHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
	e.PUT("/households/:householdId/edit", editHouseholdHandler(inventoryClient), auth.IsAuthenticated, mustHaveHousehold(inventoryClient))
}
