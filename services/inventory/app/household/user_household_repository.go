package household

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/services/inventory/shared"
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

	return r.db.Query("INSERT INTO user_households (user_id, household_id, name, location, description) VALUES (?, ?, ?, ?, ?)", userId, householdUUID, name, location, description).WithContext(ctx).Exec()
}

func (r UserHouseholdRepository) GetUserHouseholds(ctx context.Context, userId string) ([]shared.UserHousehold, error) {
	var householdId gocql.UUID
	var name, location, description string
	iter := r.db.Query("SELECT household_id, name, location, description FROM user_households WHERE user_id = ?", userId).WithContext(ctx).Iter()
	defer iter.Close()

	households := make([]shared.UserHousehold, 0)
	for iter.Scan(&householdId, &name, &location, &description) {
		households = append(households, shared.UserHousehold{
			UserID:      userId,
			HouseholdID: householdId.String(),
			Name:        name,
			Location:    location,
			Description: description,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get user households: %w", err)
	}

	return households, nil
}
