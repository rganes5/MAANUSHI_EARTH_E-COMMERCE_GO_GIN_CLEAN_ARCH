package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type CartUseCase interface {
	AddToCart(ctx context.Context, productId string, id uint) error
	ListCart(ctx context.Context, id uint, pagination utils.Pagination) (int, []utils.ResponseCart, error)
	RemoveFromCart(ctx context.Context, productId string, id uint) error
}
