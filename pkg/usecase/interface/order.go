package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type OrderUseCase interface {
	PlaceNewOrder(ctx context.Context, addressId uint, paymentId uint, userId uint) error
	ListOrders(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseOrders, error)
}
