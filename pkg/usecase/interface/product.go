package interfaces

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type ProductUseCase interface {
	AddCategory(ctx context.Context, category domain.Category) error
	UpdateCategory(ctx context.Context, categories domain.Category, id string) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]utils.ResponseCategory, error)
	AddProduct(ctx context.Context, products domain.Products) error
	DeleteProduct(ctx context.Context, id string) error
	EditProduct(ctx context.Context, product domain.Products, id string) error
	ListProducts(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseProduct, error)
	AddProductDetails(ctx context.Context, productDetails domain.ProductDetails) error
	ListProductDetailsById(ctx context.Context, id string) ([]utils.ResponseProductDetails, error)
}
