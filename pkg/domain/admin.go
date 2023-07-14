package domain

import "time"

type Admin struct {
	// gorm.Model
	ID        uint      `json:"id" gorm:"primarykey;auto_increment"`
	FirstName string    `json:"firstname" gorm:"not null"`
	LastName  string    `json:"lastname" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNum  string    `json:"Phonenum" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
