package repository

import (
	"context"
	"errors"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB}
}

// Finds whether a email is already in the database or not
func (c *adminDatabase) FindByEmail(ctx context.Context, Email string) (domain.Admin, error) {
	var admin domain.Admin
	_ = c.DB.Where("Email=?", Email).Find(&admin)
	if admin.ID == 0 {
		return domain.Admin{}, errors.New("invalid Email")
	}
	return admin, nil
}

// UserSign-up
func (c *adminDatabase) SignUpAdmin(ctx context.Context, admin domain.Admin) error {
	err := c.DB.Create(&admin).Error
	return err
}

// List all users
func (c *adminDatabase) ListUsers(ctx context.Context) ([]utils.ResponseUsers, error) {
	var users []utils.ResponseUsers
	query := `SELECT first_name,last_name,email,phone_num,block from users`
	err := c.DB.Raw(query).Scan(&users).Error
	if err != nil {
		return users, errors.New("failed to retrieve all the users")
	}
	return users, nil
}

// Manage the access of users

func (c *adminDatabase) AccessHandler(ctx context.Context, id string, access bool) error {
	err := c.DB.Model(domain.Users{}).Where("id=?", id).UpdateColumn("block", access).Error
	if err != nil {
		return errors.New("failed to update")
	}
	return nil
}

// Add category

func (c *adminDatabase) AddCategory(ctx context.Context, category domain.Category) error {
	err := c.DB.Create(&category).Error
	if err != nil {
		// return errors.New("failed to add the category")
		return err
	}
	return nil
}

// Delete category

func (c *adminDatabase) DeleteCategory(ctx context.Context, id string) error {
	err := c.DB.Where("id=?", id).Delete(&domain.Category{}).Error
	if err != nil {
		return errors.New("failed to delete the category")
	}
	return nil
}

//List categories

func (c *adminDatabase) ListCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	var categories []utils.ResponseCategory
	query := `select category_name from categories where deleted_at is null`
	err := c.DB.Raw(query).Scan(&categories).Error
	if err != nil {
		return categories, errors.New("failed to retrieve all the categories")
	}
	return categories, nil

}
