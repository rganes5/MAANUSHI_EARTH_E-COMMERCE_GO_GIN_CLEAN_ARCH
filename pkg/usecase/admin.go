package usecase

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
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

func (c *adminUseCase) ListUsers(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseUsers, error) {
	users, err := c.adminRepo.ListUsers(ctx, pagination)
	return users, err
}

func (c *adminUseCase) AccessHandler(ctx context.Context, id string, email bool) error {
	err := c.adminRepo.AccessHandler(ctx, id, email)
	return err
}

func (c *adminUseCase) Dashboard(ctx context.Context) (utils.ResponseWidgets, error) {
	responseWidgets, err := c.adminRepo.Dashboard(ctx)
	return responseWidgets, err
}

func (c *adminUseCase) SalesReport(reqData utils.SalesReport) ([]utils.ResponseSalesReport, error) {
	return c.adminRepo.SalesReport(reqData)
}
