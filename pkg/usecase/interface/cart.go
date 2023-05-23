package interfaces

import (
	"context"
)

type CartUseCase interface {
	AddToCart(ctx context.Context, productId string, id uint) error
}
