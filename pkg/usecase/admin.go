package usecase

import (
	"context"

	domain "github.com/rganes5/go-gin-clean-arch/pkg/domain"
	interfaces "github.com/rganes5/go-gin-clean-arch/pkg/repository/interface"
	services "github.com/rganes5/go-gin-clean-arch/pkg/usecase/interface"
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
	users, err := c.adminRepo.FindByEmail(ctx, Email)
	return users, err
}

func (c *adminUseCase) SignUpAdmin(ctx context.Context, admin domain.Admin) error {
	err := c.adminRepo.SignUpAdmin(ctx, admin)
	return err
}
