package usecase

import (
	"context"

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
	err := c.ProductRepo.DeleteCategory(ctx, id)
	return err
}

func (c *ProductUseCase) ListCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	categories, err := c.ProductRepo.ListCategories(ctx)
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

func (c *ProductUseCase) ListProducts(ctx context.Context) ([]utils.ResponseProductAdmin, error) {
	products, err := c.ProductRepo.ListProducts(ctx)
	return products, err
}
