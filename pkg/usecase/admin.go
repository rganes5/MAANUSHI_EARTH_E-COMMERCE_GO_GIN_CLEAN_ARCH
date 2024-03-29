package usecase

import (
	"context"
	"errors"
	"time"

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

func (c *adminUseCase) SignUpAdmin(ctx context.Context, body utils.AdminSignUp) (domain.Admin, error) {
	//Check whether such email already exits
	if _, err := c.adminRepo.FindByEmail(ctx, body.Email); err == nil {
		return domain.Admin{}, errors.New("user already exists with the email")
	}
	newAdminOutput, err := c.adminRepo.SignUpAdmin(ctx, body)
	return newAdminOutput, err
}

func (c *adminUseCase) FindByEmail(ctx context.Context, Email string) (domain.Admin, error) {
	admin, err := c.adminRepo.FindByEmail(ctx, Email)
	return admin, err
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

func (c *adminUseCase) GetAllCoupons(ctx context.Context, pagination utils.Pagination) ([]domain.Coupon, error) {
	coupons, err := c.adminRepo.GetAllCoupons(ctx, pagination)
	return coupons, err
}

func (c *adminUseCase) DeleteCoupon(ctx context.Context, couponId string) error {
	err := c.adminRepo.DeleteCoupon(ctx, couponId)
	return err
}

func (c *adminUseCase) AddCoupon(ctx context.Context, couponBody utils.BodyAddCoupon) error {
	date, err1 := time.Parse("2006-01-02", couponBody.ExpirationDate)
	if err1 != nil {
		return err1
	}
	coupon := domain.Coupon{
		CouponCode:         couponBody.Code,
		CouponType:         couponBody.Type,
		Discount:           couponBody.Discount,
		UsageLimit:         couponBody.UsageLimit,
		ExpirationDate:     date,
		MinimumOrderAmount: couponBody.MinOrderAmount,
		ProductID:          couponBody.ProductID,
		CategoryID:         couponBody.CategoryID,
	}
	err2 := c.adminRepo.AddCoupon(ctx, coupon)
	return err2
}

func (c *adminUseCase) UpdateCoupon(ctx context.Context, couponBody utils.BodyAddCoupon, couponId string) error {
	date, err1 := time.Parse("2006-01-02", couponBody.ExpirationDate)
	if err1 != nil {
		return err1
	}
	coupon := domain.Coupon{
		CouponCode:         couponBody.Code,
		CouponType:         couponBody.Type,
		Discount:           couponBody.Discount,
		UsageLimit:         couponBody.UsageLimit,
		ExpirationDate:     date,
		MinimumOrderAmount: couponBody.MinOrderAmount,
		ProductID:          couponBody.ProductID,
		CategoryID:         couponBody.CategoryID,
	}
	err2 := c.adminRepo.UpdateCoupon(ctx, coupon, couponId)
	return err2
}
