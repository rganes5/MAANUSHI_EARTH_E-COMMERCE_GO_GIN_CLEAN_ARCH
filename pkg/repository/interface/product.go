package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type ProductRepository interface {
	AddCategory(ctx context.Context, category domain.Category) error
	UpdateCategory(ctx context.Context, categories domain.Category, id string) error
	DeleteCategory(ctx context.Context, id string) error
	CheckItemsPresent(ctx context.Context, id string) (domain.Products, error)
	ListCategories(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseCategory, error)
	AddProduct(ctx context.Context, product domain.Products) error
	DeleteProduct(ctx context.Context, id string) error
	EditProduct(ctx context.Context, product domain.Products, id string) error
	ListProducts(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseProduct, error)
	ListProductsBasedOnCategory(ctx context.Context, id string, pagination utils.Pagination) ([]utils.ResponseProduct, error)
	AddProductDetails(ctx context.Context, productDetails domain.ProductDetails) error
	EditProductDetailsById(ctx context.Context, product_details domain.ProductDetails, id string) error
	ListProductDetailsById(ctx context.Context, id string) ([]utils.ResponseProductDetails, error)
	ListProductAndDetailsById(ctx context.Context, id string) ([]utils.ResponseProductAndDetails, error)
}
