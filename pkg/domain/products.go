package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	CategoryName string `json:"categoryname" gorm:"uniqueIndex;not null"`
}

type products struct {
	gorm.Model
}
