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
		WithSetCacheFn(c.cache, fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, household.UserID), func() any {
			if !household.IsFromOnboarding {
				return nil
			}

			return []shared.UserHousehold{
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
