package usecase

import (
	"context"
	"fmt"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type ProductUseCase struct {
	ProductRepo interfaces.ProductRepository
}

func NewProductUseCase(repo interfaces.ProductRepository) services.ProductUseCase {
	return &ProductUseCase{
		ProductRepo: repo,
	}
}

func (c *ProductUseCase) AddCategory(ctx context.Context, category domain.Category) error {
	err := c.ProductRepo.AddCategory(ctx, category)
	return err
}

func (c *ProductUseCase) UpdateCategory(ctx context.Context, categories domain.Category, id string) error {
	err := c.ProductRepo.UpdateCategory(ctx, categories, id)
	return err
}

func (c *ProductUseCase) DeleteCategory(ctx context.Context, id string) error {
	product, err1 := c.ProductRepo.CheckItemsPresent(ctx, id)
	if err1 != nil {
		return err1
	}
	if product.ID != 0 {
		return fmt.Errorf("items exist under the category")
	}
	err := c.ProductRepo.DeleteCategory(ctx, id)
	return err
}

func (c *ProductUseCase) ListCategories(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseCategory, error) {
	categories, err := c.ProductRepo.ListCategories(ctx, pagination)
	return categories, err
}

func (c *ProductUseCase) AddProduct(ctx context.Context, product domain.Products) error {
	err := c.ProductRepo.AddProduct(ctx, product)
	return err
}

func (c *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	err := c.ProductRepo.DeleteProduct(ctx, id)
	return err
}

func (c *ProductUseCase) EditProduct(ctx context.Context, product domain.Products, id string) error {
	err := c.ProductRepo.EditProduct(ctx, product, id)
	return err
}

func (c *ProductUseCase) ListProducts(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseProduct, error) {
	products, err := c.ProductRepo.ListProducts(ctx, pagination)
	return products, err
}

func (c *ProductUseCase) AddProductDetails(ctx context.Context, productDetails domain.ProductDetails) error {
	err := c.ProductRepo.AddProductDetails(ctx, productDetails)
	return err
}

func (c *ProductUseCase) ListProductDetailsById(ctx context.Context, id string) ([]utils.ResponseProductDetails, error) {
	productDetails, err := c.ProductRepo.ListProductDetailsById(ctx, id)
	return productDetails, err
}

func (c *ProductUseCase) ListProductAndDetailsById(ctx context.Context, id string) ([]utils.ResponseProductAndDetails, error) {
	productAndDetails, err := c.ProductRepo.ListProductAndDetailsById(ctx, id)
	return productAndDetails, err
}
