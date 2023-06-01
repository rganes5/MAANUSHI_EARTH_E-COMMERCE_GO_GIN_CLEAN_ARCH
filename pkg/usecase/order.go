package usecase

import (
	"context"
	"time"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
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
	if err != nil {
		return err
	}

	cartItems, err1 := c.OrderRepo.FindCartItems(ctx, cart.ID)
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

	if err := c.OrderRepo.SubmitOrder(ctx, Neworder, cartItems); err != nil {
		return err
	}
	return nil
}
