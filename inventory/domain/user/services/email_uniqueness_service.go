package services

import (
	"context"

	"github.com/cybre/home-inventory/inventory/domain/user"
)

type EmailUniquenessRepository interface {
	IsUniqueEmail(ctx context.Context, email string) (bool, error)
}

type EmailUniquenessService struct {
	emailUniquenessRepository EmailUniquenessRepository
}

func NewEmailUniquenessService(emailUniquenessRepository EmailUniquenessRepository) *EmailUniquenessService {
	return &EmailUniquenessService{
		emailUniquenessRepository: emailUniquenessRepository,
	}
}

func (s *EmailUniquenessService) IsUnique(ctx context.Context, email user.Email) (bool, error) {
	return s.emailUniquenessRepository.IsUniqueEmail(ctx, email.String())
}
