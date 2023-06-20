package repository

import (
	"context"
	"errors"
	"fmt"

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

// Find the user via email or phone number
func (c *userDatabase) FindByEmailOrNumber(ctx context.Context, body utils.OtpLogin) (domain.Users, error) {
	var user domain.Users
	_ = c.DB.Where("email=? or phone_num=?", body.Email, body.PhoneNum).Find(&user)
	fmt.Println("The user is", user)
	if user.ID == 0 {
		return domain.Users{}, errors.New("user with such email or phone number does not exist in database")
	}
	if user.Block {
		return user, errors.New("you are blocked")

	}
	if !user.Verified {
		return user, errors.New("your phone number is not verified")
	}
	return user, nil
}

// UserSign-up
func (c *userDatabase) SignUpUser(ctx context.Context, user domain.Users) (string, error) {
	tx := c.DB.Begin()
	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return user.PhoneNum, err
	}
	query := `insert into carts(user_id)values($1)`
	if err := tx.Exec(query, user.ID).Error; err != nil {
		tx.Rollback()
		return user.PhoneNum, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
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
// func (c *userDatabase) ListProducts(ctx context.Context) ([]utils.ResponseProductUser, error) {
// 	var products []utils.ResponseProductUser
// 	query := `select product_name,image,details,price,discount_price from products where deleted_at is null`
// 	err := c.DB.Raw(query).Scan(&products).Error
// 	if err != nil {
// 		return products, err
// 	}
// 	return products, nil
// }

// User home handler
func (c *userDatabase) HomeHandler(ctx context.Context, id uint) (utils.ResponseUsersDetails, error) {
	var user utils.ResponseUsersDetails
	query := `SELECT id,first_name,last_name,email,phone_num from users where id=?`
	err := c.DB.Raw(query, id).Scan(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// Update profile
func (c *userDatabase) UpdateProfile(ctx context.Context, updateProfile domain.Users, id uint) error {
	err := c.DB.Model(&domain.Users{}).Where("id=?", id).Updates(&updateProfile).Error
	if err != nil {
		return err
	}
	return nil
}

// Add address
func (c *userDatabase) AddAddress(ctx context.Context, address domain.Address) error {
	err := c.DB.Create(&address).Error
	if err != nil {
		return err
	}
	return nil
}

// List addresses
func (c *userDatabase) ListAddress(ctx context.Context, id uint, pagination utils.Pagination) ([]utils.ResponseAddress, error) {
	offset := pagination.Offset
	limit := pagination.Limit
	var address []utils.ResponseAddress
	query := `SELECT id, name, phone_number, house, area, land_mark, city, pincode, state, country, "primary" FROM addresses WHERE deleted_at IS NULL AND user_id = ? LIMIT ? OFFSET ?`
	err := c.DB.Raw(query, id, limit, offset).Scan(&address).Error
	if err != nil {
		return address, err
	}
	return address, nil
}

// Edit address
func (c *userDatabase) UpdateAddress(ctx context.Context, updateAddress domain.Address, id string) error {
	err := c.DB.Model(&domain.Address{}).Where("id=?", id).Updates(&updateAddress).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete address
func (c *userDatabase) DeleteAddress(ctx context.Context, id string) error {
	err := c.DB.Where("id=?", id).Delete(&domain.Address{}).Error
	if err != nil {
		return err
	}
	return nil
}

// Update new password
func (c *userDatabase) ChangePassword(ctx context.Context, NewHashedPassword string, PhoneNum string) error {
	err := c.DB.Model(&domain.Users{}).Where("phone_num=?", PhoneNum).UpdateColumn("password", NewHashedPassword)
	if err.RowsAffected == 0 {
		return errors.New("no row updated")
	} else if err.Error != nil {
		return err.Error
	}
	return nil
}

// List the wallet
func (c *userDatabase) ViewWallet(ctx context.Context, userId uint) ([]utils.Wallet, int, error) {
	var wallet []utils.Wallet
	var totalBalance int
	if err := c.DB.Model(&domain.Wallet{}).Where("user_id=?", userId).Scan(&wallet).Error; err != nil {
		return wallet, totalBalance, err
	}
	if err := c.DB.Model(&domain.Wallet{}).Select("sum(amount) as balance").Where("user_id=?", userId).Scan(&totalBalance).Error; err != nil {
		return wallet, totalBalance, err
	}

	return wallet, totalBalance, nil
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
