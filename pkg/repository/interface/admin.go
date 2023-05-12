package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
)

type AdminRepository interface {
	// FindAll(ctx context.Context) ([]domain.Users, error)
	FindByEmail(ctx context.Context, Email string) (domain.Admin, error)
	SignUpAdmin(ctx context.Context, admin domain.Admin) error
	// FindByID(ctx context.Context, id uint) (domain.Users, error)
	// Save(ctx context.Context, user domain.Users) (domain.Users, error)
	// Delete(ctx context.Context, user domain.Users) error
}
