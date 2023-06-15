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
// Checking whether the category contains items
func (c *productDatabase) CheckItemsPresent(ctx context.Context, id string) (domain.Products, error) {
	var product domain.Products
	err1 := c.DB.Where("category_id=?", id).Find(&product).Error
	if err1 != nil {
		return product, err1
	}
	return product, nil
}

func (c *productDatabase) DeleteCategory(ctx context.Context, id string) error {
	err := c.DB.Where("id=?", id).Delete(&domain.Category{}).Error
	if err != nil {
		return err
	}
	return nil
}

//List categories

func (c *productDatabase) ListCategories(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseCategory, error) {
	var categories []utils.ResponseCategory
	offset := pagination.Offset
	limit := pagination.Limit
	query := `select id,category_name from categories where deleted_at is null LIMIT $1 OFFSET $2`
	err := c.DB.Raw(query, limit, offset).Scan(&categories).Error
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
	err := c.DB.Model(&domain.Products{}).Where("id = ?", id).UpdateColumns(&product).Error
	if err != nil {
		return err
	}
	return nil
}

//List products on admins end

func (c *productDatabase) ListProducts(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseProduct, error) {
	var products []utils.ResponseProduct
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT
    products.id,
    products.product_name,
    products.image,
    products.details,
	products.price,
	products.discount_price,
    products.category_id,
	categories.category_name
	FROM
    products
	INNER JOIN
    categories ON products.category_id = categories.id
	WHERE
    products.deleted_at IS NULL
	LIMIT $1 OFFSET $2`
	err := c.DB.Raw(query, limit, offset).Scan(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}

// List products based on category
func (c *productDatabase) ListProductsBasedOnCategory(ctx context.Context, id string, pagination utils.Pagination) ([]utils.ResponseProduct, error) {
	var products []utils.ResponseProduct
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT
    products.id,
    products.product_name,
    products.image,
    products.details,
	products.price,
	products.discount_price,
    products.category_id,
	categories.category_name
	FROM
    products
	INNER JOIN
    categories ON products.category_id = categories.id
	WHERE
    categories.id = $1
    AND products.deleted_at IS NULL
	LIMIT $2 OFFSET $3`
	err := c.DB.Raw(query, id, limit, offset).Scan(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}

// Add product details
func (c *productDatabase) AddProductDetails(ctx context.Context, productDetails domain.ProductDetails) error {
	err := c.DB.Create(&productDetails).Error
	if err != nil {
		return err
	}
	return nil
}

// Edit product details
func (c *productDatabase) EditProductDetailsById(ctx context.Context, product_details domain.ProductDetails, id string) error {
	err := c.DB.Model(&domain.ProductDetails{}).Where("id = ?", id).UpdateColumns(&product_details).Error
	if err != nil {
		return err
	}
	return nil
}

// List  product details
func (c *productDatabase) ListProductDetailsById(ctx context.Context, id string) ([]utils.ResponseProductDetails, error) {
	var productDetails []utils.ResponseProductDetails
	query := `SELECT id,product_id,product_details,in_stock FROM product_details WHERE product_id = ? AND deleted_at IS NULL`
	err := c.DB.Raw(query, id).Scan(&productDetails).Error
	if err != nil {
		return productDetails, err
	}
	return productDetails, nil
}

// List product and details
func (c *productDatabase) ListProductAndDetailsById(ctx context.Context, id string) ([]utils.ResponseProductAndDetails, error) {
	var productAndDetails []utils.ResponseProductAndDetails
	query := `SELECT product_id,products.product_name, products.image, products.details, products.price, 
	products.discount_price, product_details.product_details, product_details.in_stock FROM products 
	JOIN product_details ON products.id = product_details.product_id WHERE products.id = ? AND products.deleted_at IS NULL`
	err := c.DB.Raw(query, id).Scan(&productAndDetails).Error
	if err != nil {
		return productAndDetails, err
	}
	return productAndDetails, nil
}
