package interfaces

import (
	"context"
	"time"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type OrderRepository interface {
	FindCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error)
	FindOrderItemsbyId(ctx context.Context, orderDetailsId uint) (domain.OrderDetails, time.Time, error)
	SubmitOrder(ctx context.Context, order domain.Order, cartItems []domain.CartItem) error
	CancelOrder(ctx context.Context, userId uint, orderItems domain.OrderDetails) error
	UpdateStatus(ctx context.Context, orderItem domain.OrderDetails) error
	ReturnOrder(ctx context.Context, orderItem domain.OrderDetails) error
	ListOrders(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseOrders, error)
	AdminListOrders(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseOrdersAdmin, error)
	ListOrderDetails(ctx context.Context, orderId uint, pagination utils.Pagination) ([]utils.ResponseOrderDetails, error)
	FindCoupon(ctx context.Context, code string) (domain.Coupon, error)
	ValidateCoupon(ctx context.Context, coupon domain.Coupon, cartItems []domain.CartItem, cart *domain.Cart) error
}
