package domain

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	FirstName string `json:"firstname" binding:"required" gorm:"not null"`
	LastName  string `json:"lastname" binding:"required" gorm:"not null"`
	Email     string `json:"email" binding:"required" gorm:"uniqueIndex;not null"`
	PhoneNum  string `json:"phonenum" binding:"required" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" binding:"required" gorm:"not null"`
	Block     bool   `json:"block" gorm:"default:false"`
	Verified  bool   `json:"verified" gorm:"default:false"`
}

type Address struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	House       string `json:"house" gorm:"not null" binding:"required"`
	Area        string `json:"area"`
	LandMark    string `json:"land_mark" gorm:"not null" binding:"required"`
	City        string `json:"city"  binding:"required"`
	Pincode     uint   `json:"pincode" gorm:"not null" binding:"required"`
	State       string `json:"state" binding:"required" gorm:"not null"`
	Country     string `json:"country" binding:"required" gorm:"not null"`
	Primary     bool   `json:"primary" gorm:"default:false"`
	UserID      uint   `json:"userid"`
}
