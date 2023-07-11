package domain

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	// ID        uint   `json:"id" gorm:"primarykey;auto_increment"`
	FirstName string `json:"firstname" gorm:"not null"`
	LastName  string `json:"lastname" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNum  string `json:"Phonenum" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" gorm:"not null"`
}
