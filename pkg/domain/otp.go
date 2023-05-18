package domain

type OtpSession struct {
	ID       uint   `json:"id" gorm:"primarykey;auto_increment"`
	OtpID    string `json:"otpid" gorm:"not null"`
	PhoneNum string `json:"phonenum" gorm:"not null"`
}
