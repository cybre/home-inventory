package http

import (
	"context"
	"net/http"

	"github.com/cybre/home-inventory/internal/app"
	"github.com/go-chi/chi/v5"
)

type ItemService interface {
	AddItem(context.Context, app.AddItemCommandData) error
}

func NewHTTPTransport(itemService ItemService) http.Handler {
	router := chi.NewRouter()
	buildItemRoutes(router, itemService)

	return router
}
