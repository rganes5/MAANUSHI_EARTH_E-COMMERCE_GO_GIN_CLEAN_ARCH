package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
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

func (c *OrderUseCase) PlaceNewOrder(ctx context.Context, addressId uint, paymentId uint, userId uint) error {
	cart, err := c.CartRepo.FindCartById(ctx, userId)
	fmt.Println("user id passed from to find cart by id from place new order function from use case is", userId)
	fmt.Println("Cart found from the find cart by id from place new order function from use case is", cart)
	if err != nil {
		return err
	}
	fmt.Println("1 cartid which is used to find the cartitems table from the use case is", cart.ID)

	cartItems, err1 := c.OrderRepo.FindCartItems(ctx, cart.ID)
	fmt.Println("2 cartid which is used to find the cartitems table from the use case is", cart.ID)
	if err1 != nil {
		return err
	}

	Neworder := domain.Order{
		UserID:     cart.UserID,
		PlacedDate: time.Now(),
		AddressID:  addressId,
		PaymentID:  paymentId,
		GrandTotal: uint(cart.GrandTotal),
	}
	fmt.Println("New order is", Neworder)
	if err := c.OrderRepo.SubmitOrder(ctx, Neworder, cartItems); err != nil {
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
