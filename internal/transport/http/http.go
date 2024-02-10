package http

import (
	"context"
	"net/http"

	"github.com/cybre/home-inventory/internal/app/item"
	"github.com/go-chi/chi/v5"
)

type ItemService interface {
	AddItem(context.Context, item.AddItemCommandData) error
}

func NewHTTPTransport(itemService ItemService) error {
	router := chi.NewRouter()
	buildItemRoutes(router, itemService)

	return http.ListenAndServe(":8080", router)
}
