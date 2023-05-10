package utils

type Otpverify struct {
	Otp string `binding:"required"`
}

type LoginBody struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}
