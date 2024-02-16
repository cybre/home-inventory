package services

import (
	"context"

	"github.com/cybre/home-inventory/inventory/domain/user"
)

type EmailUniquenessChecker interface {
	IsUniqueEmail(ctx context.Context, email string) (bool, error)
}

type EmailUniquenessService struct {
	emailUniquenessChecker EmailUniquenessChecker
}

func NewEmailUniquenessService(emailUniquenessRepository EmailUniquenessChecker) *EmailUniquenessService {
	return &EmailUniquenessService{
		emailUniquenessChecker: emailUniquenessRepository,
	}
}

func (s *EmailUniquenessService) IsUnique(ctx context.Context, email user.Email) (bool, error) {
	return s.emailUniquenessChecker.IsUniqueEmail(ctx, email.String())
}
