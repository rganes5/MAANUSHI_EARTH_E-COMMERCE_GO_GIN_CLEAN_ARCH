package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	"gorm.io/gorm"
)

type cartDatabase struct {
	DB *gorm.DB
}

func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
	return &cartDatabase{DB}
}

// To find the cart by user id
func (c *cartDatabase) FindCartById(ctx context.Context, id uint) (domain.Cart, error) {
	var cart domain.Cart
	if err := c.DB.Where("user_id=?", id).Find(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

// To find the products by product id
func (c *cartDatabase) FindProductById(ctx context.Context, productId string) (domain.Products, error) {
	var product domain.Products
	// Convert the string productId to a uint
	pID, err := strconv.ParseUint(productId, 10, 64)
	if err != nil {
		return product, err
	}
	if err := c.DB.Where("id=?", pID).Find(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}

// To find the product details by id
func (c *cartDatabase) FindProductDetailsById(ctx context.Context, productId string) (domain.ProductDetails, error) {
	var productDetails domain.ProductDetails
	// Convert the string productId to a uint
	pDID, err := strconv.ParseUint(productId, 10, 64)
	if err != nil {
		return productDetails, err
	}
	if err := c.DB.Where("product_id = ?", pDID).Find(&productDetails).Error; err != nil {
		return productDetails, err
	}

	return productDetails, nil
}

// To find the duplicate product so that we can update the quantity
func (c *cartDatabase) FindDuplicateProduct(ctx context.Context, productId string, cartID uint) (domain.CartItem, error) {
	var duplicateItem domain.CartItem
	// Convert the string productId to a uint
	pID, err := strconv.ParseUint(productId, 10, 64)
	if err != nil {
		return duplicateItem, err
	}
	fmt.Println("this is the productId from the repository of duplicate function", pID)
	fmt.Println("this is the cartId from the repository of duplicate function", cartID)
	if err := c.DB.Where("product_id=$1 and cart_id=$2", pID, cartID).Find(&duplicateItem).Error; err != nil {
		return duplicateItem, err
	}
	return duplicateItem, nil
}

func (c *cartDatabase) ListCart(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseCart, error) {
	var cartDetails []utils.ResponseCart
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT
			 products.id,
			 products.product_name,
			 products.image,
			 products.details,
			 categories.category_name,
			 products.price,
			 products.discount_price,
			 cart_items.quantity,
			 cart_items.total_price
			 FROM
			 products
			 INNER JOIN cart_items ON cart_items.product_id=products.id
			 INNER JOIN categories ON products.category_id=categories.id
			 INNER JOIN carts ON cart_items.cart_id=carts.id
			 WHERE carts.user_id=?
			 LIMIT $2 OFFSET $3`
	err := c.DB.Raw(query, id, limit, offset).Scan(&cartDetails).Error
	if err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}

//result := c.DB.Where("product_detail_id=$1 and cart_id=$2", id, cartid).Find(&exsistitem)

func (c *cartDatabase) UpdateCartItem(ctx context.Context, existingItem domain.CartItem) error {
	var grantTotal int
	tx := c.DB.Begin()
	if err := tx.Model(&domain.CartItem{}).Where("id=?", existingItem.ID).UpdateColumns(&existingItem).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.CartItem{}).Where("cart_id=?", existingItem.CartID).Select("SUM(total_price)").Scan(&grantTotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("id=?", existingItem.CartID).UpdateColumn("grand_total", grantTotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("duplicate item updated from update function is", existingItem.ProductId)
	return nil
}

func (c *cartDatabase) AddNewItem(ctx context.Context, newItem domain.CartItem) error {
	var grantTotal int
	tx := c.DB.Begin()
	if err := tx.Create(&newItem).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.CartItem{}).Where("cart_id=?", newItem.CartID).Select("SUM(total_price)").Scan(&grantTotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("id=?", newItem.CartID).UpdateColumn("grand_total", grantTotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("new item updated from add item is", newItem.ProductId)
	return nil
}

func (c *cartDatabase) DeleteFromCart(ctx context.Context, existingItem domain.CartItem) error {
	tx := c.DB.Begin()
	if err := tx.Model(&domain.CartItem{}).Where("id=?", existingItem.ID).Delete(&existingItem).Error; err != nil {
		tx.Rollback()
		return err
	}
	var grantTotal sql.NullInt64
	if err := tx.Model(&domain.CartItem{}).Where("cart_id=?", existingItem.CartID).Select("SUM(total_price)").Scan(&grantTotal).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&domain.Cart{}).Where("id=?", existingItem.CartID).UpdateColumn("grand_total", grantTotal).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
