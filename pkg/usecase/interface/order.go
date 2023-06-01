package interfaces

import "context"

type OrderUseCase interface {
	PlaceNewOrder(ctx context.Context, addressId uint, paymentId uint, userId uint) error
}
