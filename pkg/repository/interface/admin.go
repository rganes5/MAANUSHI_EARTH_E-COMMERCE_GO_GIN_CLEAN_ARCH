package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type AdminRepository interface {
	// FindAll(ctx context.Context) ([]domain.Users, error)
	FindByEmail(ctx context.Context, Email string) (domain.Admin, error)
	SignUpAdmin(ctx context.Context, admin domain.Admin) error
	ListUsers(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseUsers, error)
	AccessHandler(ctx context.Context, id string, access bool) error
	Dashboard(ctx context.Context) (utils.ResponseWidgets, error)
	SalesReport(utils.SalesReport) ([]utils.ResponseSalesReport, error)
	AddCoupon(ctx context.Context, coupon domain.Coupon) error
	GetAllCoupons(ctx context.Context, pagination utils.Pagination) ([]domain.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon domain.Coupon, couponId string) error
	DeleteCoupon(ctx context.Context, couponId string) error

	// FindByID(ctx context.Context, id uint) (domain.Users, error)
	// Save(ctx context.Context, user domain.Users) (domain.Users, error)
	// Delete(ctx context.Context, user domain.Users) error
}
