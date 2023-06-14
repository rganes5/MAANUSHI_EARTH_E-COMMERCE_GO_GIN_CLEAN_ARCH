package interfaces

import (
	"context"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type UserRepository interface {
	// FindAll(ctx context.Context) ([]domain.Users, error)
	FindByEmail(ctx context.Context, Email string) (domain.Users, error)
	FindByEmailOrNumber(ctx context.Context, body utils.OtpLogin) (domain.Users, error)
	SignUpUser(ctx context.Context, user domain.Users) (string, error)
	UpdateVerify(ctx context.Context, PhoneNum string) error
	// ListProducts(ctx context.Context) ([]utils.ResponseProductUser, error)
	AddAddress(ctx context.Context, address domain.Address) error
	ListAddress(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseAddress, error)
	DeleteAddress(ctx context.Context, id string) error
	UpdateAddress(ctx context.Context, updateAddress domain.Address, id string) error
	HomeHandler(ctx context.Context, id uint) (utils.ResponseUsersDetails, error)
	UpdateProfile(ctx context.Context, updateProfile domain.Users, id uint) error
	ChangePassword(ctx context.Context, NewHashedPassword string, PhoneNum string) error
	ViewWallet(ctx context.Context, userId uint) ([]utils.Wallet, int, error)

	// FindByID(ctx context.Context, id uint) (domain.Users, error)
	// Save(ctx context.Context, user domain.Users) (domain.Users, error)
	// Delete(ctx context.Context, user domain.Users) error
}
