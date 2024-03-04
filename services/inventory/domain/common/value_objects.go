package common

import (
	"github.com/bnkamalesh/errors"
)

type UserID string

func NewUserID(id string) (UserID, error) {
	if id == "" {
		return "", errors.InputBody("user id is required")
	}

	return UserID(id), nil
}

func (id UserID) String() string {
	return string(id)
}
