package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

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
    order_statuses.status AS order_status,
    payment_statuses.status AS payment_status
FROM
    orders
INNER JOIN addresses ON orders.address_id = addresses.id
INNER JOIN payment_statuses ON orders.payment_status_id = payment_statuses.id
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
	var walletBalance int
	tx := c.DB.Begin()
	//creating a new order table
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	//Checking wheather the order is placed using wallet
	if order.PaymentID == 3 {
		err := tx.Model(&domain.Wallet{}).Select("COALESCE(sum(amount), 0)").Where("user_id = ?", order.UserID).Scan(&walletBalance).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if walletBalance < int(order.GrandTotal) {
			tx.Rollback()
			return errors.New("insufficient balance in wallet, choose a different payment option")
		}

		current := time.Now()
		wallet := domain.Wallet{
			UserID:       order.UserID,
			CreditedDate: nil,
			DebitedDate:  &current,
			Amount:       -1 * int(order.GrandTotal),
		}
		if err := tx.Create(&wallet).Error; err != nil {
			tx.Rollback()
			return err
		}
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

// Finding the orderItems by passing the order details id
func (c *orderDatabase) FindOrderItemsbyId(ctx context.Context, orderDetailsId uint) (domain.OrderDetails, time.Time, error) {
	var order domain.Order
	var orderItems domain.OrderDetails
	var date time.Time
	// if err := c.DB.Where("id=?", orderId).Find(&order).Error; err != nil {
	// 	return order, orderItems, date, err
	// }
	if err := c.DB.Where("id=?", orderDetailsId).Find(&orderItems).Error; err != nil {
		return orderItems, date, err
	}
	if err := c.DB.Model(&order).Select("placed_date").Where("id=?", orderItems.OrderID).Scan(&date).Error; err != nil {
		return orderItems, date, err
	}
	return orderItems, date, nil
}

// Order cancellation
func (c *orderDatabase) CancelOrder(ctx context.Context, userId uint, orderItems domain.OrderDetails) error {
	var coupon domain.Coupon
	TempProductDetails := struct {
		DiscountedPrice uint
		QtyInStock      uint
		PaymentMode     int
		ProductID       int
	}{
		DiscountedPrice: 0,
		QtyInStock:      0,
		PaymentMode:     0,
		ProductID:       0,
	}
	var grandTotal int
	tx := c.DB.Begin()

	//Retrivals
	//to find the stock and product id
	if err := tx.Model(&domain.ProductDetails{}).
		Where("id=?", orderItems.ProductDetailID).
		Select("in_stock,product_id").
		Scan(&TempProductDetails).Error; err != nil {
		tx.Rollback()
		return err
	}
	//to find the price of selected item from product table
	err := tx.Raw("SELECT p.discount_price FROM products p JOIN product_details pd ON p.id = pd.product_id WHERE pd.id = ?", TempProductDetails.ProductID).Scan(&TempProductDetails.DiscountedPrice).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// if err := tx.Model(&domain.Products{}).
	// 	Where("id=?", TempProductDetails.ProductID).Select("discount_price").
	// 	// Joins("JOIN product_details ON products.product_id = product_details.product_id").Where("product_details.id = ?", orderItems.ProductDetailID).
	// 	// Select("products.discounted_price").
	// 	Scan(&TempProductDetails).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	//we can retrive the payment id so that in case if it was not orderered by cash on delivery, we can refund the amount to wallet
	if err := tx.Model(&domain.Order{}).Where("id=?", orderItems.OrderID).Select("payment_id").Scan(&TempProductDetails.PaymentMode).Error; err != nil {
		tx.Rollback()
		return err
	}
	//we can retreive the actual grand total and store it in a temporary variable so that we can use it.
	if err := tx.Model(&domain.Order{}).Where("id=?", orderItems.OrderID).Select("grand_total").Scan(&grandTotal).Error; err != nil {
		tx.Rollback()
		return err
	}

	prodtotal := orderItems.Quantity * TempProductDetails.DiscountedPrice
	// prodtotal := orderItems.TotalPrice

	total := grandTotal - int(prodtotal)
	//to find the coupon
	result := tx.Model(&domain.Coupon{}).Joins("join orders on coupons.id=orders.coupon_id").Where("orders.id=?", orderItems.OrderID).Find(&coupon)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected != 0 {
		if total < int(*coupon.MinimumOrderAmount) || *coupon.ProductID == TempProductDetails.ProductID {
			if coupon.CouponType == 1 {
				coupondiscount := (coupon.Discount * TempProductDetails.DiscountedPrice) / 100
				total += int(coupondiscount)
				prodtotal -= coupondiscount
			} else if coupon.CouponType == 2 {
				total += int(coupon.Discount)
				prodtotal -= coupon.Discount
			}
		}
	}

	//Updations
	//to update the order details fields as cancelled and date and time of cancellation.
	if err := tx.Model(&domain.OrderDetails{}).Where("id=?", orderItems.ID).UpdateColumns(&orderItems).Error; err != nil {
		tx.Rollback()
		return err
	}
	//we can upate the stock on actual product since the item is now cancelled
	if err := tx.Model(&domain.ProductDetails{}).Where("id=?", orderItems.ProductDetailID).UpdateColumn("in_stock", (TempProductDetails.QtyInStock + orderItems.Quantity)).Error; err != nil {
		tx.Rollback()
		return err
	}
	//Calculate the total price of cancelled item and reduce the price from the grand total of the order.
	totalPrice := grandTotal - int(orderItems.TotalPrice)
	if err := tx.Model(&domain.Order{}).Where("id=?", orderItems.OrderID).UpdateColumn("grand_total", totalPrice).Error; err != nil {
		tx.Rollback()
		return err
	}
	if TempProductDetails.PaymentMode != 1 {
		current := time.Now()
		wallet := domain.Wallet{
			UserID:       userId,
			CreditedDate: &current,
			DebitedDate:  nil,
			Amount:       int(orderItems.TotalPrice),
		}
		if err := tx.Create(&wallet).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// admins end
//
// List Orders
func (c *orderDatabase) AdminListOrders(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseOrdersAdmin, error) {
	fmt.Println("Entered into list orders")
	var listUsers []utils.ResponseOrdersAdmin
	offset := pagination.Offset
	limit := pagination.Limit
	query := `SELECT DISTINCT ON (orders.id)
    orders.id,
    orders.placed_date,
    orders.grand_total,
	payment_modes.mode,
	users.first_name AS primary_user_name,
	users.phone_num AS primary_user_phone_number,
    addresses.name,
    addresses.phone_number,
    addresses.house,
    addresses.area,
    addresses.land_mark,
    addresses.city,
    addresses.state,
    addresses.country,
    addresses.pincode,
	order_statuses.status AS order_status,
    payment_statuses.status AS payment_status	FROM
		orders
	INNER JOIN addresses ON orders.address_id = addresses.id
	INNER JOIN users ON orders.user_id=users.id
	INNER JOIN payment_statuses ON orders.payment_status_id = payment_statuses.id
INNER JOIN order_details ON order_details.order_id = orders.id
INNER JOIN order_statuses ON order_details.order_status_id = order_statuses.id
INNER JOIN payment_modes ON orders.payment_id = payment_modes.id
INNER JOIN product_details ON order_details.product_detail_id = product_details.id
INNER JOIN products ON product_details.product_id = products.id
	LIMIT $1 OFFSET $2`
	err := c.DB.Raw(query, limit, offset).Scan(&listUsers).Error
	if err != nil {
		return listUsers, err
	}
	return listUsers, nil
}

func (c *orderDatabase) UpdateStatus(ctx context.Context, orderItem domain.OrderDetails) error {
	TempProductDetails := struct {
		DiscountedPrice uint
		QtyInStock      uint
		PaymentMode     int
		ProductID       int
	}{
		DiscountedPrice: 0,
		QtyInStock:      0,
		PaymentMode:     0,
		ProductID:       0,
	}
	var grandTotal int
	var userId uint
	tx := c.DB.Begin()
	//Updations
	//to update the order details fields as confirmed, shipped, out for delivery, cancelled, returned and date and time of delivery.
	if err := tx.Model(&domain.OrderDetails{}).Where("id=?", orderItem.ID).UpdateColumns(&orderItem).Error; err != nil {
		tx.Rollback()
		return err
	}
	if orderItem.OrderStatusID == 8 || orderItem.OrderStatusID == 10 {
		//Retrivals
		//To find the user id so that it can be used to update the wallet
		if err := tx.Model(&domain.Order{}).Where("id=?", orderItem.OrderID).Select("user_id").Scan(&userId).Error; err != nil {
			tx.Rollback()
			return err
		}
		//to find the stock and product id
		if err := tx.Model(&domain.ProductDetails{}).
			Where("id=?", orderItem.ProductDetailID).
			Select("in_stock,product_id").
			Scan(&TempProductDetails).Error; err != nil {
			tx.Rollback()
			return err
		}
		//we can retrive the payment id so that in case if it was not orderered by cash on delivery, we can refund the amount to wallet
		if err := tx.Model(&domain.Order{}).Where("id=?", orderItem.OrderID).Select("payment_id").Scan(&TempProductDetails.PaymentMode).Error; err != nil {
			tx.Rollback()
			return err
		}
		//we can retreive the actual grand total and store it in a temporary variable so that we can use it.
		if err := tx.Model(&domain.Order{}).Where("id=?", orderItem.OrderID).Select("grand_total").Scan(&grandTotal).Error; err != nil {
			tx.Rollback()
			return err
		}

		//Updations
		//we can upate the stock on actual product since the item is now cancelled
		if err := tx.Model(&domain.ProductDetails{}).Where("id=?", orderItem.ProductDetailID).UpdateColumn("in_stock", (TempProductDetails.QtyInStock + orderItem.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
		//Calculate the total price of cancelled item and reduce the price from the grand total of the order.
		totalPrice := grandTotal - int(orderItem.TotalPrice)
		if err := tx.Model(&domain.Order{}).Where("id=?", orderItem.OrderID).UpdateColumn("grand_total", totalPrice).Error; err != nil {
			tx.Rollback()
			return err
		}
		if TempProductDetails.PaymentMode != 1 {
			current := time.Now()
			wallet := domain.Wallet{
				UserID:       userId,
				CreditedDate: &current,
				DebitedDate:  nil,
				Amount:       int(orderItem.TotalPrice),
			}
			if err := tx.Create(&wallet).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (c *orderDatabase) ReturnOrder(ctx context.Context, orderItem domain.OrderDetails) error {
	if err := c.DB.Model(&domain.OrderDetails{}).Where("id=?", orderItem.ID).UpdateColumns(&orderItem).Error; err != nil {
		return err
	}
	return nil
}

func (c *orderDatabase) FindCoupon(ctx context.Context, code string) (domain.Coupon, error) {
	var coupon domain.Coupon
	result := c.DB.Model(&domain.Coupon{}).Where("coupon_code = ?", code).Find(&coupon)
	if result.Error != nil {
		return coupon, result.Error
	}
	if result.RowsAffected == 0 {
		return coupon, errors.New("coupon doesn't exist")
	}
	return coupon, nil
}

func (c *orderDatabase) ValidateCoupon(ctx context.Context, coupon domain.Coupon, cartItems []domain.CartItem, cart *domain.Cart) error {
	fmt.Println("Entered into the validate coupon repo")
	tx := c.DB.Begin()
	var useCount uint
	var found bool
	prodetail := struct {
		ProductId uint
		Price     uint
	}{
		ProductId: 0,
		Price:     0,
	}
	fmt.Println("1")
	if coupon.MinimumOrderAmount != nil && *coupon.MinimumOrderAmount > uint(cart.GrandTotal) {
		return errors.New("requires a minimum amount for the coupon to apply")
	}
	fmt.Println("2")
	if time.Now().After(coupon.ExpirationDate) {
		return errors.New("the coupon had expired")
	}
	fmt.Println("3")
	if err := tx.Model(&domain.CouponUsage{}).Where("coupon_id=? and user_id=?", coupon.ID, cart.UserID).Select("usage").Scan(&useCount); err.Error != nil {
		tx.Rollback()
		return err.Error
	} else if err.RowsAffected == 0 {
		couponusage := domain.CouponUsage{
			UserID:   cart.UserID,
			CouponID: coupon.ID,
			Usage:    0,
		}
		if err1 := tx.Create(&couponusage).Error; err1 != nil {
			tx.Rollback()
			return err1
		}
	}
	fmt.Println("4")
	if useCount >= coupon.UsageLimit {
		return errors.New("coupon usage limit exceeds")
	}
	fmt.Println("5")
	if coupon.ProductID != nil {
		for _, v := range cartItems {
			if useCount >= coupon.UsageLimit {
				break
			}
			if err := tx.Model(&domain.ProductDetails{}).Where("product_id=?", v.ProductId).Select("ID").Scan(&prodetail.ProductId).Error; err != nil {
				tx.Rollback()
				return err
			}
			err := tx.Raw("SELECT p.discount_price FROM products p JOIN product_details pd ON p.id = pd.product_id WHERE pd.id = ?", v.ProductId).Scan(&prodetail.Price).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			fmt.Println("the product id in the coupon is ", *coupon.ProductID)
			fmt.Println("the product id in the found product in cart is is ", int(prodetail.ProductId))

			if *coupon.ProductID == int(prodetail.ProductId) {
				found = true
				if coupon.CouponType == 2 {
					discount := (prodetail.Price * coupon.Discount) / 100
					cart.GrandTotal = cart.GrandTotal - int(discount)
				} else if coupon.CouponType == 1 {
					cart.GrandTotal = cart.GrandTotal - int(coupon.Discount)
				}
				useCount++
			}
		}
		if !found {
			return errors.New("this coupon can't be aplied for these products")
		}
	} else {
		for useCount < coupon.UsageLimit {
			if coupon.CouponType == 2 {
				discount := (cart.GrandTotal * int(coupon.Discount)) / 100
				cart.GrandTotal = cart.GrandTotal - int(discount)
			} else if coupon.CouponType == 1 {
				cart.GrandTotal = cart.GrandTotal - int(coupon.Discount)
			}
			useCount++
		}
	}
	if err := tx.Model(&domain.CouponUsage{}).Where("coupon_id=?", coupon.ID).UpdateColumn("usage", useCount).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&domain.Cart{}).Where("user_id=?", cart.UserID).Updates(cart).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
