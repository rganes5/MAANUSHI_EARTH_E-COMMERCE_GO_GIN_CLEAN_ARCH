package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type OrderRepository interface {
	FindCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error)
	SubmitOrder(ctx context.Context, order domain.Order, cartItems []domain.CartItem) error
	ListOrders(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseOrders, error)
	AdminListOrders(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseOrdersAdmin, error)
	ListOrderDetails(ctx context.Context, orderId uint, pagination utils.Pagination) ([]utils.ResponseOrderDetails, error)
}
