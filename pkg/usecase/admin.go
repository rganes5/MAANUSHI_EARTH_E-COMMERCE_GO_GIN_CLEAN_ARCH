package usecase

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
)

type adminUseCase struct {
	adminRepo interfaces.AdminRepository
}

func NewAdminUseCase(repo interfaces.AdminRepository) services.AdminUseCase {
	return &adminUseCase{
		adminRepo: repo,
	}
}

func (c *adminUseCase) FindByEmail(ctx context.Context, Email string) (domain.Admin, error) {
	admin, err := c.adminRepo.FindByEmail(ctx, Email)
	return admin, err
}

func (c *adminUseCase) SignUpAdmin(ctx context.Context, admin domain.Admin) error {
	err := c.adminRepo.SignUpAdmin(ctx, admin)
	return err
}

// func (c *adminUseCase) ListUsers(ctx context.Context) error {
// 	err := c.adminRepo.ListUsers()
// 	return err
// }
