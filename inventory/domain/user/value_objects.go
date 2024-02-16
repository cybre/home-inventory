package user

import (
	"fmt"
	"net/mail"
)

const (
	MinFirstNameLength = 2
	MaxFirstNameLength = 50

	MinLastNameLength = 2
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

type Email string

func NewEmail(email string) (Email, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}
