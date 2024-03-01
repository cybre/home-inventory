package household

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

type UserHouseholdRepository struct {
	db *gocql.Session
}

func NewUserHouseholdRepository(db *gocql.Session) *UserHouseholdRepository {
	return &UserHouseholdRepository{db: db}
}

func (r UserHouseholdRepository) InsertHousehold(ctx context.Context, userId string, householdId string, name string, location string, description string) error {
	householdUUID, err := gocql.ParseUUID(householdId)
	if err != nil {
		return fmt.Errorf("invalid household ID: %s", householdId)
	}

	return r.db.Query("INSERT INTO user_households (user_id, household_id, name, location, description, item_count, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)", userId, householdUUID, name, location, description, 0, time.Now().UnixMilli()).WithContext(ctx).Exec()
}

func (r UserHouseholdRepository) UpdateHousehold(ctx context.Context, userId string, householdId string, name string, location string, description string) error {
	householdUUID, err := gocql.ParseUUID(householdId)
	if err != nil {
		return fmt.Errorf("invalid household ID: %s", householdId)
	}

	return r.db.Query("UPDATE user_households SET name = ?, location = ?, description = ? WHERE user_id = ? AND household_id = ?", name, location, description, userId, householdUUID).WithContext(ctx).Exec()
}

func (r UserHouseholdRepository) GetUserHouseholds(ctx context.Context, userId string) ([]UserHouseholdModel, error) {
	var householdId gocql.UUID
	var name, location, description string
	var item_count int
	var timestamp int64
	var rooms []UserHouseholdRoomModel
	iter := r.db.Query("SELECT household_id, name, location, description, item_count, rooms, timestamp FROM user_households WHERE user_id = ?", userId).WithContext(ctx).Iter()
	defer iter.Close()

	households := make([]UserHouseholdModel, 0)
	for iter.Scan(&householdId, &name, &location, &description, &item_count, &rooms, &timestamp) {
		households = append(households, UserHouseholdModel{
			UserID:      userId,
			HouseholdID: householdId.String(),
			Name:        name,
			Location:    location,
			Description: description,
			ItemCount:   item_count,
			Rooms:       rooms,
			Timestamp:   timestamp,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get user households: %w", err)
	}

	return households, nil
}

func (r UserHouseholdRepository) GetUserHousehold(ctx context.Context, userId string, householdId string) (UserHouseholdModel, error) {
	householdUUID, err := gocql.ParseUUID(householdId)
	if err != nil {
		return UserHouseholdModel{}, fmt.Errorf("invalid household ID: %s", householdId)
	}

	var name, location, description string
	var item_count int
	var rooms []UserHouseholdRoomModel
	var timestamp int64
	if err := r.db.Query("SELECT name, location, description, item_count, rooms, timestamp FROM user_households WHERE user_id = ? AND household_id = ?", userId, householdUUID).WithContext(ctx).Scan(&name, &location, &description, &item_count, &rooms, &timestamp); err != nil {
		return UserHouseholdModel{}, fmt.Errorf("failed to get user household: %w", err)
	}

	return UserHouseholdModel{
		UserID:      userId,
		HouseholdID: householdId,
		Name:        name,
		Location:    location,
		Description: description,
		ItemCount:   item_count,
		Rooms:       rooms,
		Timestamp:   timestamp,
	}, nil
}
