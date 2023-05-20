package repository

import (
	"context"
	"errors"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
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
func (c *userDatabase) SignUpUser(ctx context.Context, user domain.Users) (string, error) {
	err := c.DB.Create(&user).Error
	if err != nil {
		return user.PhoneNum, err
	}
	return user.PhoneNum, nil
}

// Verify the otp column
func (c *userDatabase) UpdateVerify(ctx context.Context, PhoneNum string) error {
	err := c.DB.Model(&domain.Users{}).Where("phone_num=?", PhoneNum).UpdateColumn("Verified", true).Error
	if err != nil {
		return err
	}
	return nil
}

// List products
func (c *userDatabase) ListProducts(ctx context.Context) ([]utils.ResponseProductUser, error) {
	var products []utils.ResponseProductUser
	query := `select product_name,image,details,price,discount_price from products where deleted_at is null`
	err := c.DB.Raw(query).Scan(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
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
