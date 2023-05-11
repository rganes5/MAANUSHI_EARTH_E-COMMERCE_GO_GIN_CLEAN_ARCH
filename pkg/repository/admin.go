package repository

import (
	"context"
	"errors"

	domain "github.com/rganes5/go-gin-clean-arch/pkg/domain"
	interfaces "github.com/rganes5/go-gin-clean-arch/pkg/repository/interface"
	"gorm.io/gorm"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB}
}

// Finds whether a email is already in the database or not and also checks whether a user is blocked or not
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
