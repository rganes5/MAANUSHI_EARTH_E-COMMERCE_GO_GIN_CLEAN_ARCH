package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/auth"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/support"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type UserHandler struct {
	userUseCase services.UserUseCase
	otpUseCase  services.OtpUseCase
}

// type Response struct {
// 	ID   uint   `copier:"must"`
// 	Name string `copier:"must"`
// }

func NewUserHandler(usecase services.UserUseCase, otpusecase services.OtpUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
		otpUseCase:  otpusecase,
	}
}

// Variable declared contataining type as users which is already initialiazed in domain folder.
var signUp_user domain.Users

// var otp_user domain.Users

// USER SIGN-UP WITH OTP SENDING
func (cr *UserHandler) UserSignUp(c *gin.Context) {
	if err := c.BindJSON(&signUp_user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := support.Email_validator(signUp_user.Email); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := support.MobileNum_validator(signUp_user.PhoneNum); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _, err := cr.userUseCase.FindByEmail(c.Request.Context(), signUp_user.Email); err == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User with the email already exits!",
		})
		return
	}

	fmt.Println("config variable", config.GetCofig())

	signUp_user.Password, _ = support.HashPassword(signUp_user.Password)
	PhoneNum, err := cr.userUseCase.SignUpUser(c.Request.Context(), signUp_user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":         "failed to add user",
			"error_details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User data": "Saved",
	})

	respSid, err1 := cr.otpUseCase.TwilioSendOTP(c.Request.Context(), PhoneNum)

	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":         "Failed generating otp",
			"error_details": err1.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success":    "Enter the otp",
		"responseid": respSid,
	})
}

// SIGN UP OTP VERIFICATION

func (cr *UserHandler) SignupOtpverify(c *gin.Context) {
	var otp utils.OtpVerify
	if err := c.BindJSON(&otp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error binding json",
		})
		return
	}

	session, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), otp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"error1": "verification failed",
		})
		return
	}

	err1 := cr.userUseCase.UpdateVerify(c.Request.Context(), session.PhoneNum)
	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":  err1.Error(),
			"error1": "updation end fails",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User registration": "Success",
	})
}

// USERLOGIN
func (cr *UserHandler) LoginHandler(c *gin.Context) {
	//Cookie check
	_, err1 := c.Cookie("user-token")
	if err1 == nil {
		c.Redirect(http.StatusFound, "/user/home")
		// c.AbortWithStatusJSON(http.StatusFound, gin.H{
		// 	"alert": "User already logged in and cookie present",
		// })
		return
	}
	//Login logic
	var loginBody utils.LoginBody
	if err := c.BindJSON(&loginBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//Checks whether such user email exits or not and also returns back the user details of that specific user related to the email and stores in user.
	user, err := cr.userUseCase.FindByEmail(c.Request.Context(), loginBody.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	//Checks the given password with retreived password to that specific email from the database(user variable)
	if err := support.CheckPasswordHash(loginBody.Password, user.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	//GenerateJWT function from the auth package, passing user.Email and User.ID as an argument. It assigns the generated JWT to the tokenString variable
	tokenString, err := auth.GenerateJWT(user.Email, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error generating the jwt token",
		})
	}

	//Sets a cookie named "user-token" with the value tokenString. The cookie has an expiration time of 60 minutes from the current time.
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.Set("user-email", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"Login": "Success",
	})
}

// USERLOGOUT
func (cr *UserHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("user-token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"logout": "Success",
	})
}

// HomeHandler
func (cr *UserHandler) Homehandler(c *gin.Context) {
	email, ok := c.Get(("user-email"))
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
	}
	user, err := cr.userUseCase.FindByEmail(c.Request.Context(), email.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}
	c.JSON(http.StatusOK, user)
	// user, err := cr.userUseCase.Homehandler(c.Request.Context())
	// if err != nil {
	// 	response := utils.ErrorResponse(400, "Failed to fetch user details", err.Error(), nil)
	// 	c.JSON(http.StatusInternalServerError, response)
	// }
}

// ListProduct
func (cr *UserHandler) ListProducts(c *gin.Context) {
	products, err := cr.userUseCase.ListProducts(c.Request.Context())

	if err != nil {
		response := utils.ErrorResponse(500, "Failed to list products", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "All product details", products)
	c.JSON(http.StatusOK, response)
}

// Add Address
func (cr *UserHandler) AddAddress(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	var body utils.Address
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(http.StatusBadRequest, "Failed to bind JSON body", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	body.UserID = id.(uint)
	var address domain.Address
	copier.Copy(&address, &body)
	if err := cr.userUseCase.AddAddress(c.Request.Context(), address); err != nil {
		response := utils.ErrorResponse(500, "Failed to add address", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Successfully added the address", nil)
	c.JSON(http.StatusOK, response)

}

//
//
//
//
//
//
//
//
//
//
//
//

// FindAll godoc
// @summary Get all users
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /api/users [get]
// @response 200 {object} []Response "OK"
// func (cr *UserHandler) FindAll(c *gin.Context) {
// 	users, err := cr.userUseCase.FindAll(c.Request.Context())

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusNotFound)
// 	} else {
// 		response := []Response{}
// 		copier.Copy(&response, &users)

// 		c.JSON(http.StatusOK, response)
// 	}
// }
