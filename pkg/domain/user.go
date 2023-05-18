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
