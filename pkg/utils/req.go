package utils

type Pagination struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}

type OtpVerify struct {
	Otp         string `json:"otp" binding:"required"`
	OtpID       string `json:"otpid" binding:"required"`
	NewPassword string `json:"newpassword"`
}

type LoginBody struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type OtpLogin struct {
	Email    string `json:"email"`
	PhoneNum string `json:"phonenum"`
}

type Products struct {
	ProductName   string `json:"productname" binding:"required" gorm:"uniqueIndex;not null"`
	Image         string `json:"image" gorm:"not null"`
	Details       string `json:"details"`
	Price         uint   `json:"price" gorm:"not null" binding:"required,numeric"`
	DiscountPrice uint   `json:"discountprice"`
	CategoryID    uint   `json:"categoryid"`
	// Category    Category `gorm:"foreignKey:CategoryID"`
}

type ProductDetails struct {
	ProductID      uint   `json:"productid"`
	ProductDetails string `json:"productdetails"`
	InStock        uint   `json:"qty_in_stock"`
}

type UpdateCategory struct {
	CategoryName string `json:"categoryname" binding:"required" gorm:"uniqueIndex;not null" `
}

type Address struct {
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

type UpdateAddress struct {
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

type UpdateProfile struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}
