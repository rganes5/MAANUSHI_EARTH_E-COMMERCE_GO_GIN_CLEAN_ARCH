package utils

import (
	"strings"
	"time"
)

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
	ID            uint   `json:"productid"`
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

// struct used to view orders
type ResponseOrders struct {
	ID            uint      `json:"id"`
	PlacedDate    time.Time `json:"placeddate"`
	Name          string    `json:"name"`
	PhoneNumber   string    `json:"phonenumber"`
	House         string    `json:"house"`
	Area          string    `json:"area"`
	LandMark      string    `json:"landmark"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	Pincode       string    `json:"pincode"`
	Mode          string    `json:"mode"`
	OrderStatus   string    `json:"orderstatus"`
	PaymentStatus string    `json:"paymentstatus"`
	GrandTotal    uint      `json:"grandtotal"`
}

// struct used to view the order_details
type ResponseOrderDetails struct {
	ID               uint       `json:"id "`
	ProductName      string     `json:"productname"`
	Details          string     `json:"details"`
	Image            string     `json:"image"`
	CategoryName     string     `json:"categoryname"`
	Price            uint       `json:"price"`
	DiscountPrice    uint       `json:"discountprice"`
	Quantity         uint       `json:"quantity"`
	TotalPrice       uint       `json:"totalprice"`
	Status           string     `json:"status"`
	DeliveredDate    *time.Time `json:"delivereddate"`
	CancelledDate    *time.Time `json:"cancelleddate"`
	ReturnSubmitDate *time.Time `json:"returnsubmitdate"`
	// Percentage    int        `json:"discountpercentage"`
}

// struct used to view orders
type ResponseOrdersAdmin struct {
	ID                     uint      `json:"id"`
	PlacedDate             time.Time `json:"placeddate"`
	PrimaryUserName        string    `json:"primaryusername"`
	PrimaryUserPhoneNumber string    `json:"primaryuserphonenumber"`
	Name                   string    `json:"name"`
	PhoneNumber            string    `json:"phonenumber"`
	House                  string    `json:"house"`
	Area                   string    `json:"area"`
	LandMark               string    `json:"landmark"`
	City                   string    `json:"city"`
	State                  string    `json:"state"`
	Country                string    `json:"country"`
	Pincode                string    `json:"pincode"`
	Mode                   string    `json:"mode"`
	OrderStatus            string    `json:"orderstatus"`
	PaymentStatus          string    `json:"paymentstatus"`
	GrandTotal             uint      `json:"grandtotal"`
}

// struct for admin to view the widgets
type ResponseWidgets struct {
	ActiveUsers    int `json:"activeusers"`
	BlockedUsers   int `json:"blockedusers"`
	Products       int `json:"products"`
	Pendingorders  int `json:"pendingorders"`
	ReturnRequests int `json:"returnrequests"`
}

// salesreport
type ResponseSalesReport struct {
	UserID          uint      `json:"userid" gorm:"column:userid"`
	FirstName       string    `json:"firstname"`
	Email           string    `json:"email"`
	ProductDetailID uint      `json:"productdetailid" gorm:"column:productdetailid"`
	ProductName     string    `json:"productname" gorm:"column:productname"`
	Price           uint      `json:"price"`
	DiscountPrice   uint      `json:"discountprice" gorm:"column:discountprice"`
	Quantity        uint      `json:"quantity"`
	OrderID         uint      `json:"orderid" gorm:"column:orderid"`
	PlacedDate      time.Time `json:"placeddate"`
	PaymentMode     string    `json:"paymentmode" gorm:"column:paymentmode"`
	OrderStatus     string    `json:"orderstatus" gorm:"column:orderstatus"`
}

// razorpay
type RazorpayOrder struct {
	RazorpayKey     string `json:"razorpaykey"`
	AmountToPay     uint   `json:"amounttopay"`
	RazorpayAmount  int    `json:"razorpayamount"`
	RazorpayOrderID string `json:"razorpayorderid"`
	UserID          uint   `json:"userid"`
}

// Wallet
type Wallet struct {
	ID           uint       `json:"id"`
	CreditedDate *time.Time `json:"crediteddate"`
	DebitedDate  *time.Time `json:"debiteddate"`
	Amount       int        `json:"amount"`
}
