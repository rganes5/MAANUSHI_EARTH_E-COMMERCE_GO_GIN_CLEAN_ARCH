package repository

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB}
}

// Add category

func (c *productDatabase) AddCategory(ctx context.Context, category domain.Category) error {
	err := c.DB.Create(&category).Error
	if err != nil {
		// return errors.New("failed to add the category")
		return err
	}
	return nil
}

// Edit category
func (c *productDatabase) UpdateCategory(ctx context.Context, categories domain.Category, id string) error {
	err := c.DB.Model(&domain.Category{}).Where("id=?", id).Updates(&categories).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete category

func (c *productDatabase) DeleteCategory(ctx context.Context, id string) error {
	err := c.DB.Where("id=?", id).Delete(&domain.Category{}).Error
	if err != nil {
		return err
	}
	return nil
}

//List categories

func (c *productDatabase) ListCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	var categories []utils.ResponseCategory
	query := `select category_name from categories where deleted_at is null`
	err := c.DB.Raw(query).Scan(&categories).Error
	if err != nil {
		return categories, err
	}
	return categories, nil

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
		return err
	}
	return nil
}

//Edit product

func (c *productDatabase) EditProduct(ctx context.Context, product domain.Products, id string) error {
	err := c.DB.Model(&domain.Products{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return err
	}
	return nil
}

//List products on admins end

func (c *productDatabase) ListProducts(ctx context.Context) ([]utils.ResponseProductAdmin, error) {
	var products []utils.ResponseProductAdmin
	query := `select product_name,image,details,price,discount_price,category_id from products where deleted_at is null`
	err := c.DB.Raw(query).Scan(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}
