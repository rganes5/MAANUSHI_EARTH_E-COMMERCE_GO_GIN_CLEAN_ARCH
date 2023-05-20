package usecase

import (
	"context"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

// func (c *userUseCase) FindAll(ctx context.Context) ([]domain.Users, error) {
// 	users, err := c.userRepo.FindAll(ctx)
// 	return users, err
// }

func (c *userUseCase) FindByEmail(ctx context.Context, Email string) (domain.Users, error) {
	users, err := c.userRepo.FindByEmail(ctx, Email)
	return users, err
}

func (c *userUseCase) SignUpUser(ctx context.Context, user domain.Users) (string, error) {
	PhoneNum, err := c.userRepo.SignUpUser(ctx, user)
	return PhoneNum, err
}

func (c userUseCase) UpdateVerify(ctx context.Context, PhoneNum string) error {
	err := c.userRepo.UpdateVerify(ctx, PhoneNum)
	return err
}

func (c userUseCase) ListProducts(ctx context.Context) ([]utils.ResponseProductUser, error) {
	products, err := c.userRepo.ListProducts(ctx)
	return products, err
}

func (c userUseCase) AddAddress(ctx context.Context, address domain.Address) error {
	err := c.userRepo.AddAddress(ctx, address)
	return err
}

// func (c *userUseCase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
// 	user, err := c.userRepo.FindByID(ctx, id)
// 	return user, err
// }

// func (c *userUseCase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
// 	user, err := c.userRepo.Save(ctx, user)

// 	return user, err
// }

// func (c *userUseCase) Delete(ctx context.Context, user domain.Users) error {
// 	err := c.userRepo.Delete(ctx, user)

// 	return err
// }
