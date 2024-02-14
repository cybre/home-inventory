package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/inventory/shared"
	"github.com/labstack/echo/v4"
)

type HouseholdService interface {
	CreateHousehold(context.Context, shared.CreateHouseholdCommandData) error
	AddRoom(context.Context, shared.AddRoomCommandData) error
	AddItem(context.Context, shared.AddItemCommandData) error
	UpdateItem(context.Context, shared.UpdateItemCommandData) error
}

func NewHTTPTransport(ctx context.Context, householdService HouseholdService) error {
	e := echo.New()
	buildHouseholdRoutes(e, householdService)

	go func() {
		if err := e.Start(":8080"); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			panic(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
