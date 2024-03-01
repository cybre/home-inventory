package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/requestbuilder"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/labstack/echo/v4"
)

const (
	GetUserHouseholdsCacheKeyFormat = "GetUserHouseholds_%s"
	GetUserHouseholdCacheKeyFormat  = "GetUserHousehold_%s_%s"
)

type InventoryClient struct {
	address string
	cache   *cache.Cache[string]
}

func New(address string, cache *cache.Cache[string]) *InventoryClient {
	return &InventoryClient{
		address: address,
		cache:   cache,
	}
}

func (c InventoryClient) GetUserHouseholds(ctx context.Context, userId string) ([]shared.UserHousehold, error) {
	resp, err := requestbuilder.
		New(http.MethodGet, c.address+shared.UserHouseholdsRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, userId).
		WithHeader("Accept", "application/json").
		WithCache(c.cache, fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, userId)).
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
	UserID           string `json:"-"`
	HouseholdID      string `json:"householdId"`
	Name             string `json:"name"`
	Location         string `json:"location"`
	Description      string `json:"description"`
	IsFromOnboarding bool   `json:"-"`
}

func (c InventoryClient) CreateHousehold(ctx context.Context, household CreateHouseholdRequest) error {
	timestamp := time.Now().UnixMilli()
	resp, err := requestbuilder.New(http.MethodPost, c.address+shared.UserHouseholdsRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, household.UserID).
		WithBody(household).
		WithSetCacheFn(c.cache, func() (string, any) {
			if !household.IsFromOnboarding {
				return fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, household.UserID), nil
			}

			return fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, household.UserID), []shared.UserHousehold{
				{
					UserID:      household.UserID,
					HouseholdID: household.HouseholdID,
					Name:        household.Name,
					Location:    household.Location,
					Description: household.Description,
					ItemCount:   0,
					Rooms:       []shared.UserHouseholdRoom{},
					Timestamp:   timestamp,
				},
			}
		}).
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

type UpdateHouseholdRequest struct {
	UserID      string `json:"-"`
	HouseholdID string `json:"-"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func (c InventoryClient) UpdateHousehold(ctx context.Context, household UpdateHouseholdRequest) error {
	resp, err := requestbuilder.New(http.MethodPut, c.address+shared.UserHouseholdRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, household.UserID).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, household.HouseholdID).
		WithBody(household).
		WithInvalidateCache(
			c.cache,
			fmt.Sprintf(GetUserHouseholdCacheKeyFormat, household.UserID, household.HouseholdID),
			fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, household.UserID),
		).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return echo.NewHTTPError(resp.StatusCode, "failed to update household")
	}

	return nil
}

func (c InventoryClient) GetUserHousehold(ctx context.Context, userId, householdId string) (shared.UserHousehold, error) {
	resp, err := requestbuilder.
		New(http.MethodGet, c.address+shared.UserHouseholdRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, userId).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, householdId).
		WithHeader("Accept", "application/json").
		WithCache(c.cache, fmt.Sprintf(GetUserHouseholdCacheKeyFormat, userId, householdId)).
		WithRetry().
		Do(ctx)
	if err != nil {
		return shared.UserHousehold{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return shared.UserHousehold{}, echo.NewHTTPError(resp.StatusCode, "failed to get household")
	}

	defer resp.Body.Close()

	var household shared.UserHousehold
	if err := json.NewDecoder(resp.Body).Decode(&household); err != nil {
		return shared.UserHousehold{}, err
	}

	return household, nil
}
