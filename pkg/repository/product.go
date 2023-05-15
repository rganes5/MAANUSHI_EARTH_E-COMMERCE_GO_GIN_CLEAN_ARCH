package repository

import (
	"context"
	"errors"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB}
}

//Add product

func (c *productDatabase) AddProduct(ctx context.Context, product domain.Products) error {
	err := c.DB.Create(&product).Error
	if err != nil {
		return err
	}
	return nil
}

//Delete Product

func (c *productDatabase) DeleteProduct(ctx context.Context, id string) error {
	err := c.DB.Where("id=?", id).Delete(&domain.Products{}).Error
	if err != nil {
		return errors.New("failed to delete the category")
	}
	return nil
}
