package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
)

type OrderRepository interface {
	FindCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error)
	SubmitOrder(ctx context.Context, order domain.Order, cartItems []domain.CartItem) error
}
