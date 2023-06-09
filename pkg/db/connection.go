package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	db.AutoMigrate(&domain.Users{})
	db.AutoMigrate(&domain.Admin{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Products{})
	db.AutoMigrate(&domain.OtpSession{})
	db.AutoMigrate(&domain.Address{})
	db.AutoMigrate(&domain.ProductDetails{})
	db.AutoMigrate(&domain.Cart{})
	db.AutoMigrate(&domain.CartItem{})
	db.AutoMigrate(&domain.Order{})
	db.AutoMigrate(&domain.OrderDetails{})
	db.AutoMigrate(&domain.OrderStatus{})
	db.AutoMigrate(&domain.PaymentModes{})
	db.AutoMigrate(&domain.PaymentStatus{})
	db.AutoMigrate(&domain.Wallet{})
	db.AutoMigrate(&domain.Coupon{})
	db.AutoMigrate(&domain.CouponType{})
	db.AutoMigrate(&domain.CouponUsage{})

	return db, dbErr
}
