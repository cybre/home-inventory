package common

import (
	"fmt"
)

type UserID string

func NewUserID(id string) (UserID, error) {
	if id == "" {
		return "", fmt.Errorf("user id cannot be empty")
	}

	return UserID(id), nil
}

func (id UserID) String() string {
	return string(id)
}
