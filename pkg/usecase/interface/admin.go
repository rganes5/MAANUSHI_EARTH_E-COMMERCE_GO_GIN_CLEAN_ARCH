package interfaces

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type AdminUseCase interface {
	// FindAll(ctx context.Context) ([]domain.Users, error)
	FindByEmail(ctx context.Context, Email string) (domain.Admin, error)
	SignUpAdmin(ctx context.Context, admin domain.Admin) error
	ListUsers(ctx context.Context) ([]utils.ResponseUsers, error)
	AccessHandler(ctx context.Context, id string, access bool) error
	AddCategory(ctx context.Context, category domain.Category) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]utils.ResponseCategory, error)

	// FindByID(ctx context.Context, id uint) (domain.Users, error)
	// Save(ctx context.Context, user domain.Users) (domain.Users, error)
	// Delete(ctx context.Context, user domain.Users) error
}
