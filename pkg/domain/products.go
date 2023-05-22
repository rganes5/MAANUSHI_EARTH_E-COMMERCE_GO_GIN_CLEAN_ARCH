package domain

import "gorm.io/gorm"

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
