package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/internal/cache"
	"github.com/cybre/home-inventory/internal/requestbuilder"
	"github.com/cybre/home-inventory/services/inventory/shared"
	"github.com/labstack/echo/v4"
)

const (
	GetUserHouseholdsCacheKeyFormat    = "GetUserHouseholds_%s"
	GetUserHouseholdCacheKeyFormat     = "GetUserHousehold_%s_%s"
	GetUserHouseholdRoomCacheKeyFormat = "GetUserHouseholdRoom_%s_%s_%s"
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
		return nil, propagateError(resp)
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
					Rooms:       []shared.UserHouseholdRoom{},
					Timestamp:   timestamp,
					Order:       1,
				},
			}
		}, func() (string, any) {
			return fmt.Sprintf(GetUserHouseholdCacheKeyFormat, household.UserID, household.HouseholdID), shared.UserHousehold{
				UserID:      household.UserID,
				HouseholdID: household.HouseholdID,
				Name:        household.Name,
				Location:    household.Location,
				Description: household.Description,
				Rooms:       []shared.UserHouseholdRoom{},
				Timestamp:   timestamp,
				Order:       0,
			}
		}).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return propagateError(resp)
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
		return propagateError(resp)
	}

	return nil
}

func (c InventoryClient) DeleteHousehold(ctx context.Context, userId, householdId string) error {
	resp, err := requestbuilder.New(http.MethodDelete, c.address+shared.UserHouseholdRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, userId).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, householdId).
		WithInvalidateCache(
			c.cache,
			fmt.Sprintf(GetUserHouseholdCacheKeyFormat, userId, householdId),
			fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, userId),
		).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return propagateError(resp)
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
		return shared.UserHousehold{}, propagateError(resp)
	}

	defer resp.Body.Close()

	var household shared.UserHousehold
	if err := json.NewDecoder(resp.Body).Decode(&household); err != nil {
		return shared.UserHousehold{}, err
	}

	return household, nil
}

type AddRoomRequest struct {
	UserID      string `json:"-"`
	HouseholdID string `json:"-"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
}

func (c InventoryClient) AddRoom(ctx context.Context, room AddRoomRequest) error {
	resp, err := requestbuilder.New(http.MethodPost, c.address+shared.UserHouseholdRoomsRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, room.UserID).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, room.HouseholdID).
		WithBody(room).
		WithInvalidateCache(
			c.cache,
			fmt.Sprintf(GetUserHouseholdCacheKeyFormat, room.UserID, room.HouseholdID),
			fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, room.UserID),
		).
		WithSetCacheFn(c.cache, func() (string, any) {
			uhr := shared.UserHouseholdRoom{
				HouseholdID: room.HouseholdID,
				RoomID:      room.RoomID,
				Name:        room.Name,
				Order:       0,
				Timestamp:   0,
			}

			return fmt.Sprintf(GetUserHouseholdRoomCacheKeyFormat, room.UserID, room.HouseholdID, room.RoomID), uhr
		}).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return propagateError(resp)
	}

	return nil
}

type UpdateRoomRequest struct {
	UserID      string `json:"-"`
	HouseholdID string `json:"-"`
	RoomID      string `json:"-"`
	Name        string `json:"name"`
}

func (c InventoryClient) UpdateRoom(ctx context.Context, room UpdateRoomRequest) error {
	resp, err := requestbuilder.New(http.MethodPut, c.address+shared.UserHouseholdRoomRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, room.UserID).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, room.HouseholdID).
		WithPathParam(shared.UserHouseholdsRoomIDParam, room.RoomID).
		WithBody(room).
		WithInvalidateCache(
			c.cache,
			fmt.Sprintf(GetUserHouseholdCacheKeyFormat, room.UserID, room.HouseholdID),
			fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, room.UserID),
		).
		WithSetCacheFn(c.cache, func() (string, any) {
			uhr := shared.UserHouseholdRoom{
				HouseholdID: room.HouseholdID,
				RoomID:      room.RoomID,
				Name:        room.Name,
				Order:       0,
				Timestamp:   0,
			}

			return fmt.Sprintf(GetUserHouseholdRoomCacheKeyFormat, room.UserID, room.HouseholdID, room.RoomID), uhr
		}).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return propagateError(resp)
	}

	return nil
}

func (c InventoryClient) DeleteRoom(ctx context.Context, userId, householdId, roomId string) error {
	resp, err := requestbuilder.New(http.MethodDelete, c.address+shared.UserHouseholdRoomRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, userId).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, householdId).
		WithPathParam(shared.UserHouseholdsRoomIDParam, roomId).
		WithInvalidateCache(
			c.cache,
			fmt.Sprintf(GetUserHouseholdCacheKeyFormat, userId, householdId),
			fmt.Sprintf(GetUserHouseholdsCacheKeyFormat, userId),
			fmt.Sprintf(GetUserHouseholdRoomCacheKeyFormat, userId, householdId, roomId),
		).
		WithRetry().
		Do(ctx)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return propagateError(resp)
	}

	return nil
}

func (c InventoryClient) GetUserHouseholdRoom(ctx context.Context, userID, householdID, roomID string) (shared.UserHouseholdRoom, error) {
	resp, err := requestbuilder.
		New(http.MethodGet, c.address+shared.UserHouseholdRoomRoute).
		WithPathParam(shared.UserHouseholdsUserIDParam, userID).
		WithPathParam(shared.UserHouseholdsHouseholdIDParam, householdID).
		WithPathParam(shared.UserHouseholdsRoomIDParam, roomID).
		WithHeader("Accept", "application/json").
		WithCache(c.cache, fmt.Sprintf(GetUserHouseholdRoomCacheKeyFormat, userID, householdID, roomID)).
		WithRetry().
		Do(ctx)
	if err != nil {
		return shared.UserHouseholdRoom{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return shared.UserHouseholdRoom{}, propagateError(resp)
	}

	defer resp.Body.Close()

	var room shared.UserHouseholdRoom
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		return shared.UserHouseholdRoom{}, err
	}

	return room, nil
}

func propagateError(resp *http.Response) error {
	defer resp.Body.Close()

	messageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(resp.StatusCode)
	}

	return echo.NewHTTPError(resp.StatusCode, string(messageBytes))
}
