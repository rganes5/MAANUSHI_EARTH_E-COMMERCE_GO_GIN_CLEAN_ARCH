package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, product domain.Products) error
	DeleteProduct(ctx context.Context, id string) error
	EditProduct(ctx context.Context, product domain.Products, id string) error
}
