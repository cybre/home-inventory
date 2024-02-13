package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/inventory/shared"
	"github.com/go-chi/chi/v5"
)

type HouseholdService interface {
	CreateHousehold(context.Context, shared.CreateHouseholdCommandData) error
	AddRoom(context.Context, shared.AddRoomCommandData) error
	AddItem(context.Context, shared.AddItemCommandData) error
	UpdateItem(context.Context, shared.UpdateItemCommandData) error
}

func NewHTTPTransport(ctx context.Context, householdService HouseholdService) error {
	router := chi.NewRouter()
	buildHouseholdRoutes(router, householdService)

	server := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			panic(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
