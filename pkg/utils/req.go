package utils

import "gorm.io/gorm"

type OtpVerify struct {
	Otp   string `json:"otp" binding:"required"`
	OtpID string `json:"otpid" binding:"required"`
}

type LoginBody struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type UpdateProducts struct {
	gorm.Model
	ProductName string `json:"productname" binding:"required" gorm:"uniqueIndex;not null"`
	CategoryID  uint   `json:"categoryid"`
	// Category    Category `gorm:"foreignKey:CategoryID"`
}
