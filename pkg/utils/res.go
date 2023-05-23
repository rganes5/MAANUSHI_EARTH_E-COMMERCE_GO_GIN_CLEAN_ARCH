package utils

import "strings"

// struct used to list all users from admins end
type ResponseUsers struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PhoneNum  string `json:"phonenum"`
	Block     bool   `json:"block"`
}

// struct used to list all users from admins end
type ResponseUsersDetails struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PhoneNum  string `json:"phonenum"`
}

// struct to list all categories from admins end
type ResponseCategory string

// struct to list all products from admins and users end
type ResponseProduct struct {
	ProductName   string `json:"productname"`
	Image         string `json:"image"`
	Details       string `json:"details"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discountprice"`
	// CategoryID    uint   `json:"categoryid"`
}

type ResponseProductDetails struct {
	// ProductID      uint   `json:"productid"`
	ProductDetails string `json:"productdetails"`
	InStock        uint   `json:"qty_in_stock"`
}

type ResponseProductAndDetails struct {
	// ProductID      uint   `json:"productid"`
	ProductName    string `json:"productname"`
	Image          string `json:"image"`
	Details        string `json:"details"`
	Price          uint   `json:"price"`
	DiscountPrice  uint   `json:"discountprice"`
	ProductDetails string `json:"productdetails"`
	InStock        uint   `json:"qty_in_stock"`
}

// type ResponseProductUser struct {
// 	ProductName   string `json:"productname"`
// 	Image         string `json:"image"`
// 	Details       string `json:"details"`
// 	Price         uint   `json:"price"`
// 	DiscountPrice uint   `json:"discountprice"`
// }

type ResponseAddress struct {
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
}

// ERROR MANAGEMENT
// Req,Res,Err coding standard
type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func SuccessResponse(statusCode int, message string, data ...interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Errors:     nil,
		Data:       data,
	}

}

func ErrorResponse(statusCode int, message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Errors:     splittedError,
		Data:       data,
	}
}
