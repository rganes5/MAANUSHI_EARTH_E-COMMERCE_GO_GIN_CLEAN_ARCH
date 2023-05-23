package domain

type Cart struct {
	ID         uint `json:"id" gorm:"primarykey;auto_increment"`
	UserID     uint `json:"userid" gorm:"not null"`
	GrandTotal int  `json:"grandtotal" gorm:"default:0"`
}

type CartItem struct {
	ID         uint `json:"id" gorm:"primarykey;auto_increment"`
	CartID     uint `json:"cartid" gorm:"not null"`
	ProductId  uint `json:"productid" gorm:"not null"`
	Quantity   uint `json:"quantity" gorm:"not null"`
	TotalPrice uint `json:"totalprice" gorm:"not null"`
}
