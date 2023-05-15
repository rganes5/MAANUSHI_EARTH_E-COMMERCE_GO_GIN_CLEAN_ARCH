package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	CategoryName string `json:"categoryname" gorm:"uniqueIndex;not null"`
}

type Products struct {
	gorm.Model
	ProductName string   `json:"productname" gorm:"uniqueIndex;not null"`
	CategoryID  uint     `json:"categoryid"`
	Category    Category `gorm:"foreignKey:CategoryID"`
}
