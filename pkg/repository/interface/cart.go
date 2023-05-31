package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type CartRepository interface {
	FindCartById(ctx context.Context, id uint) (domain.Cart, error)
	ListCart(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseCart, error)
	FindProductDetailsById(ctx context.Context, productId string) (domain.ProductDetails, error)
	FindProductById(ctc context.Context, productId string) (domain.Products, error)
	FindDuplicateProduct(ctx context.Context, productId string, cartID uint) (domain.CartItem, error)
	UpdateCartItem(ctx context.Context, existingItem domain.CartItem) error
	DeleteFromCart(ctx context.Context, existingItem domain.CartItem) error
	AddNewItem(ctx context.Context, newItem domain.CartItem) error
}
