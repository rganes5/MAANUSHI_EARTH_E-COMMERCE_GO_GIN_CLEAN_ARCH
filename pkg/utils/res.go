package utils

import "strings"

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

// struct used to list all users from admins end
type ResponseUsers struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PhoneNum  string `json:"phonenum"`
	Block     bool   `json:"block"`
}

// struct used to list all users from admins end
type ResponseUsersDetails struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PhoneNum  string `json:"phonenum"`
}

// struct to list all categories from admins end
type ResponseCategory struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"category_name"`
}

// struct to list all products from admins and users end
type ResponseProduct struct {
	CategoryName  string `json:"category_name"`
	ID            uint   `json:"id"`
	ProductName   string `json:"productname"`
	Image         string `json:"image"`
	Details       string `json:"details"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discountprice"`
	CategoryID    uint   `json:"category_id"`
}

type ResponseProductDetails struct {
	ID             uint   `json:"id"`
	ProductID      uint   `json:"productid"`
	ProductDetails string `json:"productdetails"`
	InStock        uint   `json:"qty_in_stock"`
}

type ResponseProductAndDetails struct {
	ProductID      uint   `json:"productid"`
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
	ID          uint   `json:"id"`
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
}

type ResponseCart struct {
	ProductName   string `json:"productname"`
	CategoryName  string `json:"categoryname"`
	Image         string `json:"image"`
	Details       string `json:"details"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discountprice"`
	Quantity      uint   `json:"quantity"`
	TotalPrice    uint   `json:"totalprice"`
}

type ResponseFullCart struct {
	CartItems []ResponseCart
	// AppliedCouponID uint `json:"applied_coupon_id"`
	SubTotal int `json:"subtotal"`
	// DiscountAmount  uint `json:"discount_amount"`
}
