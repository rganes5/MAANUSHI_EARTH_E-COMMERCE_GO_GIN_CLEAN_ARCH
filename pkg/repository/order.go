package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
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
	fmt.Println("Entered into find cart items function from repository with cartid", cartID)
	var cartItems []domain.CartItem
	err := c.DB.Where("cart_id=?", cartID).Find(&cartItems).Error
	fmt.Println("cartitems found from the find cart items function from order repository is", cartItems)
	if err != nil {
		fmt.Println("err", err)
		return cartItems, err
	}
	return cartItems, nil
}

// List Orders
func (c *orderDatabase) ListOrders(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseOrders, error) {
	fmt.Println("Entered into list orders")
	var listUsers []utils.ResponseOrders
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT DISTINCT ON (orders.id)
    orders.id,
    orders.placed_date,
    orders.grand_total,
	payment_modes.mode,
    addresses.name,
    addresses.phone_number,
    addresses.house,
    addresses.area,
    addresses.land_mark,
    addresses.city,
    addresses.state,
    addresses.country,
    addresses.pincode,
    order_statuses.status
FROM
    orders
INNER JOIN addresses ON orders.address_id = addresses.id
INNER JOIN order_details ON order_details.order_id = orders.id
INNER JOIN order_statuses ON order_details.order_status_id = order_statuses.id
INNER JOIN payment_modes ON orders.payment_id = payment_modes.id
INNER JOIN product_details ON order_details.product_detail_id = product_details.id
INNER JOIN products ON product_details.product_id = products.id
WHERE
    orders.user_id = $1
	LIMIT $2 OFFSET $3`
	err := c.DB.Raw(query, id, limit, offset).Scan(&listUsers).Error
	if err != nil {
		return listUsers, err
	}
	return listUsers, nil
}

// List order details
func (c *orderDatabase) ListOrderDetails(ctx context.Context, orderId uint, pagination utils.Pagination) ([]utils.ResponseOrderDetails, error) {
	var responseOrderDetails []utils.ResponseOrderDetails
	fmt.Println("order id from the repo is", orderId)
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT
    order_details.id,
    order_details.quantity,
    order_details.total_price,
    order_details.delivered_date,
    order_details.cancelled_date,
    order_details.return_submit_date,
    order_statuses.status,
    products.product_name,
    products.image,
    products.details,
    products.price,
    products.discount_price,
    categories.category_name
FROM
    order_details
    INNER JOIN orders ON order_details.order_id = orders.id
    INNER JOIN order_statuses ON order_details.order_status_id = order_statuses.id
    INNER JOIN product_details ON order_details.product_detail_id = product_details.id
    INNER JOIN products ON product_details.product_id = products.id
    INNER JOIN categories ON products.category_id = categories.id
WHERE
    order_details.order_id = $1
	LIMIT $2 OFFSET $3`
	err := c.DB.Raw(query, orderId, limit, offset).Scan(&responseOrderDetails).Error
	if err != nil {
		return responseOrderDetails, err
	}
	return responseOrderDetails, nil
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
			OrderID:          order.ID,
			OrderStatusID:    1,
			DeliveredDate:    nil,
			CancelledDate:    nil,
			ReturnSubmitDate: nil,
			ProductDetailID:  value.ProductId,
			Quantity:         value.Quantity,
			TotalPrice:       value.TotalPrice,
		}
		//creating a order detail table
		if err := tx.Create(&orderDetails).Error; err != nil {
			tx.Rollback()
			return err
		}
		//getting the stock details of each item in cart details table
		if err := tx.Model(&domain.ProductDetails{}).Where("product_id=?", value.ProductId).Select("in_stock").Scan(&stock).Error; err != nil {
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
		if err := tx.Model(&domain.ProductDetails{}).Where("product_id=?", value.ProductId).UpdateColumn("in_stock", updatedStock).Error; err != nil {
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
