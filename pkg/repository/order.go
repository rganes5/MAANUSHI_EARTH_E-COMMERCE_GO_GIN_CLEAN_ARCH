package repository

import (
	"context"
	"errors"
	"strconv"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"gorm.io/gorm"
)

type orderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) interfaces.OrderRepository {
	return &orderDatabase{DB}
}

// finding eacg cart items rows for each user and storing in a slice
func (c *orderDatabase) FindCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	var cartItems []domain.CartItem
	err := c.DB.Where("cart_id=?", cartID).Find(cartItems).Error
	if err != nil {
		return cartItems, err
	}
	return cartItems, nil
}

// create a order table with address,payment mode and also creating a order details table , checking and updating stock and also deleting the cart details tables too.
func (c *orderDatabase) SubmitOrder(ctx context.Context, order domain.Order, cartItems []domain.CartItem) error {
	var stock uint
	tx := c.DB.Begin()
	//creating a new order table
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	//adding each item one by one to order details table from products and cart items table
	for _, value := range cartItems {
		orderDetails := domain.OrderDetails{
			OrderID:         order.ID,
			OrderStatusID:   1,
			DeliveredDate:   nil,
			CancelledDate:   nil,
			ProductDetailID: value.ProductId,
			Quantity:        value.Quantity,
		}
		//creating a order detail table
		if err := tx.Create(&orderDetails).Error; err != nil {
			tx.Rollback()
			return err
		}
		//getting the stock details of each item in cart details table
		if err := tx.Model(&domain.ProductDetails{}).Where("product_id=?", value.ProductId).Select("qty_in_stock").Scan(&stock).Error; err != nil {
			tx.Rollback()
			return err
		}
		//checking the added quantity with existing quantity from actual product details table
		if int(stock-value.Quantity) < 0 {
			tx.Rollback()
			errorMessage := "Insufficient stock: not enough quantity available for product ID. Please reduce the quantity and try again." + strconv.Itoa(int(value.ProductId))
			return errors.New(errorMessage)
		}
		//updating the remaining stock after placing order from the product details table
		updatedStock := stock - value.Quantity
		if err := tx.Model(&domain.ProductDetails{}).Where("product_id=?", value.ProductId).UpdateColumn("qty_in_stock", updatedStock).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	//Deleting each cart_item when the required cart_item details are added to order details.
	query := `DELETE FROM cart_items
			WHERE cart_id=$1`
	if err := tx.Exec(query, cartItems[0].CartID).Error; err != nil {
		tx.Rollback()
		return err
	}
	//updating the grand_total of main cart to 0 for the user
	if err := tx.Model(&domain.Cart{}).Where("id=?", cartItems[0].CartID).UpdateColumn("grand_total", 0).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
