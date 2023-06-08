package utils

import "time"

type Pagination struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}

type UsersSignUp struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	PhoneNum  string `json:"phonenum" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type AdminSignUp struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	PhoneNum  string `json:"Phonenum" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type OtpVerify struct {
	Otp         string `json:"otp" binding:"required"`
	OtpID       string `json:"otpid" binding:"required"`
	NewPassword string `json:"newpassword"`
}

type OtpSignUpVerify struct {
	Otp   string `json:"otp" binding:"required"`
	OtpID string `json:"otpid" binding:"required"`
}

type LoginBody struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type OtpLogin struct {
	Email    string `json:"email"`
	PhoneNum string `json:"phonenum"`
}

type AddCategory struct {
	CategoryName string `json:"categoryname"`
}

type Products struct {
	ProductName   string `json:"productname"`
	Image         string `json:"image"`
	Details       string `json:"details"`
	Price         uint   `json:"price"`
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
	// UserID      uint   `json:"userid"`
}

type UpdateAddress struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	House       string `json:"house"`
	Area        string `json:"area"`
	LandMark    string `json:"land_mark"`
	City        string `json:"city"`
	Pincode     uint   `json:"pincode"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Primary     bool   `json:"primary"`
	// UserID      uint   `json:"userid"`
}

type UpdateProfile struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type ReqCart struct {
	ID         uint `json:"id"`
	UserID     uint `json:"userid"`
	GrandTotal int  `json:"grandtotal"`
}

type ReqCartItem struct {
	ID         uint `json:"id"`
	CartID     uint `json:"cartid"`
	ProductId  uint `json:"productid"`
	Quantity   uint `json:"quantity"`
	TotalPrice uint `json:"totalprice"`
}

type SalesReport struct {
	Month     time.Month `json:"startdate"`
	Year      int        `json:"year"`
	Frequency string     `json:"frequency"`
	EndDate   time.Time  `json:"enddate"`
	Pagination
}
