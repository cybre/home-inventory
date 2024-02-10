package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/internal/app/inventory/item"
	"github.com/go-chi/chi/v5"
)

type ItemService interface {
	AddItem(context.Context, item.AddItemCommandData) error
}

func NewHTTPTransport(ctx context.Context, itemService ItemService) error {
	router := chi.NewRouter()
	buildItemRoutes(router, itemService)

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
