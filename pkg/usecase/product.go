package usecase

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
)

type ProductUseCase struct {
	ProductRepo interfaces.ProductRepository
}

func NewProductUseCase(repo interfaces.ProductRepository) services.ProductUseCase {
	return &ProductUseCase{
		ProductRepo: repo,
	}
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
