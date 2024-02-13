package common

import (
	"fmt"

	"github.com/google/uuid"
)

type UserID string

func NewUserID(id string) (UserID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("invalid user ID. must be valid UUID: %s", id)
	}

	return UserID(uuid.String()), nil
}

func (id UserID) String() string {
	return string(id)
}
