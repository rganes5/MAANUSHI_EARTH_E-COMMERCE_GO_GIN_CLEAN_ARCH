package interfaces

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type ProductUseCase interface {
	AddCategory(ctx context.Context, category domain.Category) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]utils.ResponseCategory, error)
	AddProduct(ctx context.Context, products domain.Products) error
	DeleteProduct(ctx context.Context, id string) error
	EditProduct(ctx context.Context, product domain.Products, id string) error
}
