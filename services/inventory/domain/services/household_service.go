package services

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/services/inventory/domain/common"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
	"github.com/cybre/home-inventory/services/inventory/shared"
)

type HouseholdsGetter interface {
	GetUserHouseholds(ctx context.Context, userID string) ([]shared.UserHousehold, error)
}

type HouseholdDomainService struct {
	householdsGetter HouseholdsGetter
}

func NewHouseholdDomainService(householdsGetter HouseholdsGetter) *HouseholdDomainService {
	return &HouseholdDomainService{
		householdsGetter: householdsGetter,
	}
}

func (s HouseholdDomainService) GetHouseholdCount(ctx context.Context, userID common.UserID) (uint, error) {
	households, err := s.householdsGetter.GetUserHouseholds(ctx, userID.String())
	if err != nil {
		return 0, fmt.Errorf("failed to get user households: %w", err)
	}

	return uint(len(households)), nil
}

func (s HouseholdDomainService) CheckHouseholdNameAvailability(ctx context.Context, userID common.UserID, name household.HouseholdName) (bool, error) {
	households, err := s.householdsGetter.GetUserHouseholds(ctx, userID.String())
	if err != nil {
		return false, fmt.Errorf("failed to get user households: %w", err)
	}

	for _, household := range households {
		if household.Name == name.String() {
			return false, nil
		}
	}

	return true, nil
}
