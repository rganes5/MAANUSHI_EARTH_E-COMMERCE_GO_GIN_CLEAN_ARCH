package utils

type OtpVerify struct {
	Otp   string `json:"otp" binding:"required"`
	OtpID string `json:"otpid" binding:"required"`
}

type LoginBody struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}
