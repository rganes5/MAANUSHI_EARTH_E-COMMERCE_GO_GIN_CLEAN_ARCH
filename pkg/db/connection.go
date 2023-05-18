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

	return db, dbErr
}
