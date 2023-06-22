package utils

import "time"

// input pagination
type Pagination struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}

// input sign up user sign up details
type UsersSignUp struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	PhoneNum  string `json:"phonenum" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// input admin sign up details

type AdminSignUp struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	PhoneNum  string `json:"Phonenum" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// input verification of otp in case of forgot password
type OtpVerify struct {
	Otp         string `json:"otp" binding:"required"`
	OtpID       string `json:"otpid" binding:"required"`
	NewPassword string `json:"newpassword"`
}

// input verification of otp
type OtpSignUpVerify struct {
	Otp   string `json:"otp" binding:"required"`
	OtpID string `json:"otpid" binding:"required"`
}

// input login credentials

type LoginBody struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

// input verification of otp at time of email or phone verification
type OtpLogin struct {
	Email    string `json:"email"`
	PhoneNum string `json:"phonenum"`
}

// input category
type AddCategory struct {
	CategoryName string `json:"categoryname"`
}

// input products
type Products struct {
	ProductName   string `json:"productname"`
	Image         string `json:"image"`
	Details       string `json:"details"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discountprice"`
	CategoryID    uint   `json:"categoryid"`
	// Category    Category `gorm:"foreignKey:CategoryID"`
}

// input products details
type ProductDetails struct {
	ProductID      uint   `json:"product_id"`
	ProductDetails string `json:"productdetails"`
	InStock        uint   `json:"qty_in_stock"`
}

// input update category

type UpdateCategory struct {
	CategoryName string `json:"categoryname" binding:"required" gorm:"uniqueIndex;not null" `
}

// input address
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

// input update category

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

// input update profile details
type UpdateProfile struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

// input cart
type ReqCart struct {
	ID         uint `json:"id"`
	UserID     uint `json:"userid"`
	GrandTotal int  `json:"grandtotal"`
}

// input cart items
type ReqCartItem struct {
	ID         uint `json:"id"`
	CartID     uint `json:"cartid"`
	ProductId  uint `json:"productid"`
	Quantity   uint `json:"quantity"`
	TotalPrice uint `json:"totalprice"`
}

// input salesreport
type SalesReport struct {
	Month     time.Month `json:"startdate"`
	Year      int        `json:"year"`
	Frequency string     `json:"frequency"`
	EndDate   time.Time  `json:"enddate"`
	Pagination
}

// input coupon
type BodyAddCoupon struct {
	Code           string `json:"code" binding:"required"`
	Type           uint   `json:"type" binding:"required"`
	Discount       uint   `json:"discount" binding:"required"`
	UsageLimit     uint   `json:"usagelimit" binding:"required"`
	ExpirationDate string `json:"expdate" binding:"required"`
	MinOrderAmount *uint  `json:"minorderamount"`
	ProductID      *int   `json:"productid"`
	CategoryID     *int   `json:"categoryid"`
}
