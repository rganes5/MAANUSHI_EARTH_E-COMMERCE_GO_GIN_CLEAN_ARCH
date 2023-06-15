package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	CategoryName string `json:"categoryname" binding:"required" gorm:"uniqueIndex;not null" `
}

type Products struct {
	gorm.Model
	ProductName   string `json:"productname" binding:"required" gorm:"uniqueIndex;not null"`
	Image         string `json:"image" gorm:"not null"`
	Details       string `json:"details" binding:"required" gorm:"not null"`
	Price         uint   `json:"price" gorm:"not null" binding:"required,numeric"`
	DiscountPrice uint   `json:"discountprice"`
	CategoryID    uint   `json:"categoryid"`
	// Category Category `gorm:"foreignKey:CategoryID"`
}

type ProductDetails struct {
	gorm.Model
	ProductID      uint   `json:"productid" binding:"required,numeric"`
	ProductDetails string `json:"productdetails"`
	InStock        uint   `json:"qty_in_stock" gorm:"not null" binding:"required"`
}

type Coupon struct {
	ID                 uint      `json:"id" gorm:"primarykey;auto_increment"`
	CouponCode         string    `json:"couponcode" gorm:"uniqueIndex;not null"`
	CouponType         uint      `json:"coupontype" gorm:"not null"`
	Discount           uint      `json:"discount" gorm:"not null"`
	UsageLimit         uint      `json:"usagelimit" gorm:"default:1"`
	ExpirationDate     time.Time `json:"expirationdate" gorm:"not null"`
	MinimumOrderAmount *uint     `json:"minimumorderamount"`
	ProductID          *int      `json:"productid"`
	CategoryID         *int      `json:"categoryid"`
}

type CouponType struct {
	ID   uint   `json:"id" gorm:"primarykey;auto_increment"`
	Type string `json:"type" gorm:"not null"`
}

type CouponUsage struct {
	ID       uint `json:"id" gorm:"primarykey;auto_increment"`
	UserID   uint `json:"userid" gorm:"not null"`
	CouponID uint `json:"couponid" gorm:"not null"`
	Usage    uint `json:"usage" gorm:"not null"`
}

type Discount struct {
	ID         uint `json:"id" gorm:"primarykey;auto_increment"`
	Percentage uint `json:"percentage" gorm:"not null"`
}
