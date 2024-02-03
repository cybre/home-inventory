package http

import (
	"encoding/json"
	"net/http"

	"github.com/cybre/home-inventory/internal/app"
	"github.com/go-chi/chi/v5"
)

func buildItemRoutes(router chi.Router, itemService ItemService) {
	router.Post("/items", addItemHandler(itemService))
}

func addItemHandler(itemService ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data app.AddItemCommandData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := itemService.AddItem(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
