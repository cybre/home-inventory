package user

import (
	"fmt"
	"net"
	"net/mail"
	"time"
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

type LastLogin struct {
	Time      Time
	UserAgent UserAgent
	IP        IPAddress
}

func NewLastLogin(timestamp int64, userAgent, ipAddress string) (LastLogin, error) {
	time, err := NewTime(timestamp)
	if err != nil {
		return LastLogin{}, err
	}

	agent, err := NewUserAgent(userAgent)
	if err != nil {
		return LastLogin{}, err
	}

	ip, err := NewIPAddress(ipAddress)
	if err != nil {
		return LastLogin{}, err
	}

	return LastLogin{
		Time:      time,
		UserAgent: agent,
		IP:        ip,
	}, nil
}

type UserAgent string

func NewUserAgent(userAgent string) (UserAgent, error) {
	if len(userAgent) == 0 {
		return "", fmt.Errorf("user agent cannot be empty")
	}

	return UserAgent(userAgent), nil
}

func (u UserAgent) String() string {
	return string(u)
}

type Time time.Time

func NewTime(unix int64) (Time, error) {
	if unix <= 0 {
		return Time{}, fmt.Errorf("invalid unix timestamp: %d", unix)
	}

	return Time(time.Unix(unix, 0)), nil
}

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

type IPAddress string

func NewIPAddress(ip string) (IPAddress, error) {
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid ip address format: %s", ip)
	}

	return IPAddress(ip), nil
}

func (i IPAddress) String() string {
	return string(i)
}
