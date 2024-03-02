package household

import (
	"context"
	"fmt"
	"slices"

	"github.com/cybre/home-inventory/internal/utils"
	"github.com/gocql/gocql"
)

type UserHouseholdRepository struct {
	db *gocql.Session
}

func NewUserHouseholdRepository(db *gocql.Session) *UserHouseholdRepository {
	return &UserHouseholdRepository{db: db}
}

func (r UserHouseholdRepository) InsertHousehold(ctx context.Context, model UserHouseholdModel) error {
	return r.db.Query("INSERT INTO user_households (user_id, household_id, name, location, description, tstamp, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)", model.UserID, model.HouseholdID, model.Name, model.Location, model.Description, model.Timestamp, model.Order).WithContext(ctx).Exec()
}

func (r UserHouseholdRepository) UpdateHousehold(ctx context.Context, model UserHouseholdModel) error {
	return r.db.Query("UPDATE user_households SET name = ?, location = ?, description = ?, tstamp = ? WHERE user_id = ? AND household_id = ?", model.Name, model.Location, model.Description, model.Timestamp, model.UserID, model.HouseholdID).WithContext(ctx).Exec()
}

func (r UserHouseholdRepository) GetUserHouseholds(ctx context.Context, userId string) ([]UserHouseholdModel, error) {
	var householdId gocql.UUID
	var name, location, description string
	var timestamp int64
	var order uint
	var rooms map[string]UserHouseholdRoomModel
	iter := r.db.Query("SELECT household_id, name, location, description, rooms, tstamp, sort_order FROM user_households WHERE user_id = ?", userId).WithContext(ctx).Iter()
	defer iter.Close()

	households := make([]UserHouseholdModel, 0)
	for iter.Scan(&householdId, &name, &location, &description, &rooms, &timestamp, &order) {
		roomList := utils.Values(rooms)
		slices.SortFunc(roomList, func(a, b UserHouseholdRoomModel) int {
			if a.Order < b.Order {
				return -1
			}
			if a.Order > b.Order {
				return 1
			}
			return 0
		})

		households = append(households, UserHouseholdModel{
			UserID:      userId,
			HouseholdID: householdId,
			Name:        name,
			Location:    location,
			Description: description,
			Rooms:       roomList,
			Timestamp:   timestamp,
			Order:       order,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get user households: %w", err)
	}

	slices.SortFunc(households, func(a, b UserHouseholdModel) int {
		if a.Order < b.Order {
			return -1
		}
		if a.Order > b.Order {
			return 1
		}
		return 0
	})

	return households, nil
}

func (r UserHouseholdRepository) GetUserHousehold(ctx context.Context, userId string, householdId string) (UserHouseholdModel, error) {
	householdUUID, err := gocql.ParseUUID(householdId)
	if err != nil {
		return UserHouseholdModel{}, fmt.Errorf("invalid household ID: %s", householdId)
	}

	var name, location, description string
	var rooms map[string]UserHouseholdRoomModel
	var timestamp int64
	var order uint
	if err := r.db.Query("SELECT name, location, description, rooms, tstamp, sort_order FROM user_households WHERE user_id = ? AND household_id = ?", userId, householdUUID).WithContext(ctx).Scan(&name, &location, &description, &rooms, &timestamp, &order); err != nil {
		return UserHouseholdModel{}, fmt.Errorf("failed to get user household: %w", err)
	}

	roomList := utils.Values(rooms)
	slices.SortFunc(roomList, func(a, b UserHouseholdRoomModel) int {
		if a.Order < b.Order {
			return -1
		}
		if a.Order > b.Order {
			return 1
		}
		return 0
	})

	return UserHouseholdModel{
		UserID:      userId,
		HouseholdID: householdUUID,
		Name:        name,
		Location:    location,
		Description: description,
		Rooms:       roomList,
		Timestamp:   timestamp,
		Order:       order,
	}, nil
}

func (r UserHouseholdRepository) UpsertRoom(ctx context.Context, userId string, model UserHouseholdRoomModel) error {
	return r.db.Query("UPDATE user_households SET rooms[?] = ? WHERE user_id = ? AND household_id = ?", model.RoomID.String(), model, userId, model.HouseholdID).WithContext(ctx).Exec()
}

func (r UserHouseholdRepository) GetRoom(ctx context.Context, userId string, householdId string, roomId string) (UserHouseholdRoomModel, error) {
	householdUUID, err := gocql.ParseUUID(householdId)
	if err != nil {
		return UserHouseholdRoomModel{}, fmt.Errorf("invalid household ID: %s", householdId)
	}

	room := UserHouseholdRoomModel{}
	if err := r.db.Query("SELECT rooms[?] FROM user_households WHERE user_id = ? AND household_id = ?", roomId, userId, householdUUID).WithContext(ctx).Scan(&room); err != nil {
		return UserHouseholdRoomModel{}, fmt.Errorf("failed to get room: %w", err)
	}

	return room, nil
}
