package http

import (
	"encoding/json"
	"net/http"

	"github.com/cybre/home-inventory/internal/shared"
	"github.com/go-chi/chi/v5"
)

func buildHouseholdRoutes(router chi.Router, householdService HouseholdService) {
	router.Post("/household", createHouseholdHandler(householdService))
	router.Post("/household/{id}/rooms", addHouseholdRoomHandler(householdService))
	router.Post("/household/{id}/rooms/{roomID}/items", addItemHandler(householdService))
	router.Put("/household/{id}/rooms/{roomID}/items/{itemID}", updateItemHandler(householdService))
}

func createHouseholdHandler(householdService HouseholdService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data shared.CreateHouseholdCommandData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := householdService.CreateHousehold(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func addHouseholdRoomHandler(householdService HouseholdService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data shared.AddRoomCommandData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := householdService.AddRoom(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func addItemHandler(householdService HouseholdService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data shared.AddItemCommandData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := householdService.AddItem(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func updateItemHandler(householdService HouseholdService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data shared.UpdateItemCommandData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := householdService.UpdateItem(r.Context(), data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
