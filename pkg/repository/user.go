package repository

import (
	"context"
	"errors"

	domain "github.com/rganes5/go-gin-clean-arch/pkg/domain"
	interfaces "github.com/rganes5/go-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

// func (c *userDatabase) FindAll(ctx context.Context) ([]domain.Users, error) {
// 	var users []domain.Users
// 	err := c.DB.Find(&users).Error

// 	return users, err
// }

// Finds whether a email is already in the database or not and also checks whether a user is blocked or not
func (c *userDatabase) FindByEmail(ctx context.Context, Email string) (domain.Users, error) {
	var user domain.Users
	_ = c.DB.Where("Email=?", Email).Find(&user)
	if user.ID == 0 {
		return domain.Users{}, errors.New("invalid Email")
	}
	if user.Block {
		return user, errors.New("you are blocked")
	}
	return user, nil
}

// UserSign-up
func (c *userDatabase) SignUpUser(ctx context.Context, user domain.Users) error {
	err := c.DB.Create(&user).Error
	return err
}

// func (c *userDatabase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
// 	var user domain.Users
// 	err := c.DB.First(&user, id).Error

// 	return user, err
// }

// func (c *userDatabase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
// 	err := c.DB.Save(&user).Error

// 	return user, err
// }

// func (c *userDatabase) Delete(ctx context.Context, user domain.Users) error {
// 	err := c.DB.Delete(&user).Error

// 	return err
// }
