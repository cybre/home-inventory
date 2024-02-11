package user

import "fmt"

const (
	MinFirstNameLength = 1
	MaxFirstNameLength = 50

	MinLastNameLength = 1
	MaxLastNameLength = 50
)

type FirstName string

func NewFirstName(firstName string) (FirstName, error) {
	if len(firstName) < MinFirstNameLength || len(firstName) > MaxFirstNameLength {
		return "", fmt.Errorf("first name must be between %d and %d characters: %s", MinFirstNameLength, MaxFirstNameLength, firstName)
	}

	return FirstName(firstName), nil
}

func (n FirstName) String() string {
	return string(n)
}

type LastName string

func NewLastName(lastName string) (LastName, error) {
	if len(lastName) < MinLastNameLength || len(lastName) > MaxLastNameLength {
		return "", fmt.Errorf("last name must be between %d and %d characters: %s", MinLastNameLength, MaxLastNameLength, lastName)
	}

	return LastName(lastName), nil
}

func (n LastName) String() string {
	return string(n)
}
