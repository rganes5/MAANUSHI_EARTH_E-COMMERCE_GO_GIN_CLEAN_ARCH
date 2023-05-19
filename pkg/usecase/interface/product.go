package interfaces

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
)

type ProductUseCase interface {
	AddProduct(ctx context.Context, products domain.Products) error
	DeleteProduct(ctx context.Context, id string) error
	EditProduct(ctx context.Context, product domain.Products, id string) error
}
