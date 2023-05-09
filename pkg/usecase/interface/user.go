package interfaces

import (
	"context"

	domain "github.com/rganes5/go-gin-clean-arch/pkg/domain"
)

type UserUseCase interface {
	// FindAll(ctx context.Context) ([]domain.Users, error)
	FindByEmail(ctx context.Context, Email string) (domain.Users, error)
	SignUpUser(ctx context.Context, user domain.Users) error
	// FindByID(ctx context.Context, id uint) (domain.Users, error)
	// Save(ctx context.Context, user domain.Users) (domain.Users, error)
	// Delete(ctx context.Context, user domain.Users) error
}
