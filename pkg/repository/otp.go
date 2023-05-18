package repository

import (
	"context"
	"fmt"

	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	"gorm.io/gorm"
)

type OtpDatabase struct {
	DB *gorm.DB
}

func NewOtpRepository(Db *gorm.DB) interfaces.OtpRepository {
	return &OtpDatabase{DB: Db}
}
func (c OtpDatabase) SaveOtp(ctx context.Context, otpsession domain.OtpSession) error {
	err := c.DB.Create(&otpsession).Error
	return err
}
func (c OtpDatabase) RetrieveSession(ctx context.Context, otp string) (domain.OtpSession, error) {
	fmt.Println(otp)
	var session domain.OtpSession
	err := c.DB.Where("otp_id=?", otp).Find(&session).Error
	if err != nil {
		return session, err
	}
	return session, nil
}
