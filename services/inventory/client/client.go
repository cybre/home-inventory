package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cybre/home-inventory/internal/requestbuilder"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/labstack/echo/v4"
)

type InventoryClient struct {
	address string
}

func New(address string) *InventoryClient {
	return &InventoryClient{
		address: address,
	}
}

func (c InventoryClient) GetUserHouseholds(ctx context.Context, userId string) ([]shared.UserHousehold, error) {
	resp, err := requestbuilder.
		New(http.MethodGet, c.address+shared.UserHouseholdsRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, userId).
		WithHeader("Accept", "application/json").
		WithRetry().
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, echo.NewHTTPError(resp.StatusCode, "failed to get user households")
	}

	defer resp.Body.Close()

	var households []shared.UserHousehold
	if err := json.NewDecoder(resp.Body).Decode(&households); err != nil {
		return nil, err
	}

	return households, nil
}

type CreateHouseholdRequest struct {
	UserID      string `json:"-"`
	HouseholdID string `json:"householdId"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func (c InventoryClient) CreateHousehold(ctx context.Context, household CreateHouseholdRequest) error {
	resp, err := requestbuilder.New(http.MethodPost, shared.UserHouseholdsRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, household.UserID).
		WithBody(household).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return echo.NewHTTPError(resp.StatusCode, "failed to create household")
	}

	return nil
}
