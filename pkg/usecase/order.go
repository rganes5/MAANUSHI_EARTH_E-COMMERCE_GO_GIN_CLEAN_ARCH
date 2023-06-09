package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/support"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type OrderUseCase struct {
	OrderRepo interfaces.OrderRepository
	CartRepo  interfaces.CartRepository
}

func NewOrderUseCase(repo interfaces.OrderRepository, CartRepo interfaces.CartRepository) services.OrderUseCase {
	return &OrderUseCase{
		OrderRepo: repo,
		CartRepo:  CartRepo,
	}
}

// Users end
func (c *OrderUseCase) PlaceNewOrder(ctx context.Context, addressId uint, paymentId uint, userId uint, couponid *uint) error {
	var psId uint
	cart, err := c.CartRepo.FindCartById(ctx, userId)
	//fmt.Println("user id passed from to find cart by id from place new order function from use case is", userId)
	//fmt.Println("Cart found from the find cart by id from place new order function from use case is", cart)
	if err != nil {
		return err
	}
	//fmt.Println("1 cartid which is used to find the cartitems table from the use case is", cart.ID)

	cartItems, err1 := c.OrderRepo.FindCartItems(ctx, cart.ID)
	//fmt.Println("2 cartid which is used to find the cartitems table from the use case is", cart.ID)
	if err1 != nil {
		return err
	}
	switch paymentId {
	case 1:
		psId = 1
	case 2:
		psId = 3
	case 3:
		psId = 4
	}
	Neworder := domain.Order{
		UserID:          cart.UserID,
		PlacedDate:      time.Now(),
		AddressID:       addressId,
		PaymentID:       paymentId,
		PaymentStatusID: psId,
		GrandTotal:      uint(cart.GrandTotal),
		CouponID:        couponid,
	}
	//fmt.Println("New order is", Neworder)
	if err := c.OrderRepo.SubmitOrder(ctx, Neworder, cartItems); err != nil {
		return err
	}
	return nil
}

func (c *OrderUseCase) RazorPayOrder(ctx context.Context, userId uint, couponid *uint) (utils.RazorpayOrder, error) {
	var razorPayOrder utils.RazorpayOrder
	cart, err := c.CartRepo.FindCartById(ctx, userId)
	//fmt.Println("user id passed from to find cart by id from place new order function from use case is", userId)
	//fmt.Println("Cart found from the find cart by id from place new order function from use case is", cart)
	if err != nil {
		return razorPayOrder, err
	}
	// Generate new razorpay order
	// Razorpay amount is calculated on paisa for india so convert the actual price into paisa
	fmt.Println("razorpayment amount is", cart.GrandTotal)
	razorPayAmount := cart.GrandTotal * 100
	fmt.Println("The razorpay amount from the usecase after convertion is ", razorPayAmount)
	razorPayReceipt := "test receipt"
	razorPayOrderId, err := support.GenerateNewRazorPayOrder(razorPayAmount, razorPayReceipt)
	if err != nil {
		return razorPayOrder, err
	}
	// set all details on razopay order
	razorPayOrder.AmountToPay = uint(cart.GrandTotal)
	razorPayOrder.RazorpayAmount = razorPayAmount
	fmt.Println("razorpayamount", razorPayAmount)
	razorPayOrder.RazorpayKey = config.GetCofig().RAZORPAYKEY
	razorPayOrder.RazorpayOrderID = razorPayOrderId
	razorPayOrder.UserID = userId

	return razorPayOrder, err
}

func (c *OrderUseCase) CancelOrder(ctx context.Context, userId uint, orderDetailsId uint) error {
	//Find the corresponding order item from the order.
	orderItem, date, err := c.OrderRepo.FindOrderItemsbyId(ctx, orderDetailsId)
	if err != nil {
		return err
	}
	//fmt.Println("This is the order", orderItem)
	if orderItem.DeliveredDate != nil {
		return errors.New("order is already delivered, Please submit a return request. If not delivered, please contact customer support")
	}
	if orderItem.CancelledDate != nil {
		return errors.New("order is already cancelled")
	}
	if time.Now().After(date.Add(24 * time.Hour)) {
		return errors.New("sorry unable to cancel the order since the order is placed 24 hours ago. Cancellation time exceeded! Please return the order once delivered")
	}
	current := time.Now()
	orderItem.CancelledDate = &current
	orderItem.OrderStatusID = 9
	if err := c.OrderRepo.CancelOrder(ctx, userId, orderItem); err != nil {
		return err
	}
	return nil
}

func (c *OrderUseCase) ListOrders(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseOrders, error) {
	listOrders, err := c.OrderRepo.ListOrders(ctx, id, pagination)
	return listOrders, err
}

func (c *OrderUseCase) ListOrderDetails(ctx context.Context, orderId uint, pagination utils.Pagination) ([]utils.ResponseOrderDetails, error) {
	listOrderDetails, err := c.OrderRepo.ListOrderDetails(ctx, orderId, pagination)
	return listOrderDetails, err
}

// Admins End
func (c *OrderUseCase) AdminListOrders(ctx context.Context, pagination utils.Pagination) ([]utils.ResponseOrdersAdmin, error) {
	listOrders, err := c.OrderRepo.AdminListOrders(ctx, pagination)
	return listOrders, err
}

func (c *OrderUseCase) UpdateStatus(ctx context.Context, orderDetailsId uint, statusId uint) error {
	//Find the corresponding order item from the order.
	orderItem, _, err := c.OrderRepo.FindOrderItemsbyId(ctx, orderDetailsId)
	if err != nil {
		return err
	}
	if orderItem.CancelledDate != nil {
		return errors.New("order is already cancelled")
	}
	if statusId == 3 {
		orderItem.OrderStatusID = statusId
		if err := c.OrderRepo.UpdateStatus(ctx, orderItem); err != nil {
			return err
		}
	} else if statusId == 4 {
		orderItem.OrderStatusID = statusId
		if err := c.OrderRepo.UpdateStatus(ctx, orderItem); err != nil {
			return err
		}
	} else if statusId == 5 {
		orderItem.OrderStatusID = statusId
		if err := c.OrderRepo.UpdateStatus(ctx, orderItem); err != nil {
			return err
		}
	} else if statusId == 6 {
		if orderItem.DeliveredDate != nil {
			return errors.New("order is already delivered")
		}
		current := time.Now()
		orderItem.DeliveredDate = &current
		orderItem.OrderStatusID = statusId
		if err := c.OrderRepo.UpdateStatus(ctx, orderItem); err != nil {
			return err
		}
	} else if statusId == 8 {
		if orderItem.ReturnSubmitDate == nil {
			return errors.New("user has not requested a return for this particular item")
		}
		orderItem.OrderStatusID = statusId
		if err := c.OrderRepo.UpdateStatus(ctx, orderItem); err != nil {
			return err
		}
	} else if statusId == 10 {
		if orderItem.DeliveredDate != nil {
			return errors.New("order is already delivered. So cannot be cancelled")
		}
		current := time.Now()
		orderItem.CancelledDate = &current
		orderItem.OrderStatusID = statusId
		if err := c.OrderRepo.UpdateStatus(ctx, orderItem); err != nil {
			return err
		}
	}
	return nil
}

func (c *OrderUseCase) ReturnOrder(ctx context.Context, orderDetailsId uint, statusId uint) error {
	// Find the corresponding order item from the order.
	orderItem, _, err := c.OrderRepo.FindOrderItemsbyId(ctx, orderDetailsId)
	if err != nil {
		return err
	}
	if orderItem.CancelledDate != nil {
		return errors.New("order is already cancelled")
	}
	if orderItem.DeliveredDate == nil {
		return errors.New("order is not delivered yet")
	}
	if orderItem.ReturnSubmitDate != nil {
		return errors.New("return is already submitted for this order. Please contact customer support")
	}
	if time.Now().After(orderItem.DeliveredDate.Add(168 * time.Hour)) {
		return errors.New("returning time exceeds")
	}
	current := time.Now()
	orderItem.OrderStatusID = statusId
	orderItem.ReturnSubmitDate = &current
	if err := c.OrderRepo.ReturnOrder(ctx, orderItem); err != nil {
		return err
	}
	return nil
}

func (c *OrderUseCase) ValidateCoupon(ctx context.Context, userId uint, code string) (*uint, error) {
	cart, err := c.CartRepo.FindCartById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if code != "" {
		coupon, err2 := c.OrderRepo.FindCoupon(ctx, code)
		if err2 != nil {
			return nil, fmt.Errorf("failed to find coupon: %w", err2)
		}
		cartItems, err1 := c.OrderRepo.FindCartItems(ctx, userId)
		if err1 != nil {
			return nil, fmt.Errorf("failed to find cart items: %w", err1)
		}
		if err = c.OrderRepo.ValidateCoupon(ctx, coupon, cartItems, &cart); err != nil {
			return nil, fmt.Errorf("failed to validate coupon: %w", err)
		}
		return &coupon.ID, nil
	}
	return nil, nil
}

func (c *OrderUseCase) FindCoupon(ctx context.Context, code string) (*uint, error) {
	if code != "null" {
		fmt.Println("enters into find coupon in the usecase")
		coupon, err := c.OrderRepo.FindCoupon(ctx, code)
		if err != nil {
			return nil, err
		}
		return &coupon.ID, err
	}
	return nil, nil
}
