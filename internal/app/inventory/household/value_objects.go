package household

import "errors"

var (
	ErrInvalidHouseholdName = errors.New("invalid household name")
)

const (
	MinHouseholdNameLength = 3
	MaxHouseholdNameLength = 100
)

type Household struct {
	Name HouseholdName
}

func NewHousehold(name string) (Household, error) {
	householdName, err := NewHouseholdName(name)
	if err != nil {
		return Household{}, err
	}

	return Household{
		Name: householdName,
	}, nil
}

type HouseholdName string

func NewHouseholdName(name string) (HouseholdName, error) {
	if len(name) < MinHouseholdNameLength || len(name) > MaxHouseholdNameLength {
		return "", ErrInvalidHouseholdName
	}

	return HouseholdName(name), nil
}
