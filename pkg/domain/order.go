package domain

import "time"

type Order struct {
	ID              uint      `json:"id" gorm:"primarykey;auto_increment"`
	UserID          uint      `json:"userid" gorm:"not null"`
	PlacedDate      time.Time `json:"placeddate" gorm:"not null"`
	AddressID       uint      `json:"addressid" gorm:"not null"`
	Address         Address   `gorm:"foreignkey:AddressID"`
	PaymentID       uint      `json:"paymentid" gorm:"not null"`
	PaymentStatusID uint      `json:"paymentstatus" gorm:"not null"`
	GrandTotal      uint      `json:"grandtotal" gorm:"not null"`
}

type OrderDetails struct {
	ID               uint       `json:"id" gorm:"primarykey;auto_increment"`
	OrderID          uint       `json:"orderid" gorm:"not null"`
	OrderStatusID    uint       `json:"orderstatusid" gorm:"not null"`
	DeliveredDate    *time.Time `json:"delivereddate"`
	CancelledDate    *time.Time `json:"cancelleddate"`
	ReturnSubmitDate *time.Time `json:"returnsubmitdate"`
	ProductDetailID  uint       `json:"productdetailsid"`
	Quantity         uint       `json:"quantity" gorm:"not null"`
	TotalPrice       uint       `json:"totalprice" gorm:"not null"`
}

type OrderStatus struct {
	ID     uint   `json:"id" gorm:"primarykey;auto_increment"`
	Status string `json:"status" gorm:"not null"`
}

type PaymentModes struct {
	ID   uint   `json:"id" gorm:"primarykey;auto_increment"`
	Mode string `json:"mode" gorm:"not null"`
}

type PaymentStatus struct {
	ID     uint   `json:"id" gorm:"primarykey;auto_increment"`
	Status string `json:"status" gorm:"not null"`
}
