package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	_ "github.com/jinzhu/gorm"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/auth"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	domain "github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	support "github.com/rganes5/maanushi_earth_e-commerce/pkg/support"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
	_ "gorm.io/gorm"
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

// var otp_user domain.Users

// @title MAANUSHI_EARTH_E-COMMERCE REST API
// @version 2.0
// @description MAANUSHI_EARTH_E-COMMERCE REST API built using Go, PSQL, REST API following Clean Architecture.

// @contact
// name: Ganesh R
// url: https://github.com/rganes5
// email: ganeshraveendranit@gmail.com

// @license
// name: MIT
// url: https://opensource.org/licenses/MIT

// @host localhost:3000

// @Basepath /
// @Accept json
// @Produce json
// @Router / [get]

// USER SIGN-UP WITH SENDING OTP
// @Summary API FOR NEW USER SIGN UP
// @ID SIGNUP-USER
// @Description CREATE A NEW USER WITH REQUIRED DETAILS
// @Tags USER
// @Accept json
// @Produce json
// @Param user_details body utils.UsersSignUp true "New user Details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/signup [post]
func (cr *UserHandler) UserSignUp(c *gin.Context) {
	var body utils.UsersSignUp

	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var signUp_user domain.Users
	copier.Copy(&signUp_user, &body)
	if err := support.Email_validator(signUp_user.Email); err != nil {
		response := utils.ErrorResponse(400, "Error: Enter a valid email. Email format is incorrect", err.Error(), signUp_user)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := support.MobileNum_validator(signUp_user.PhoneNum); err != nil {
		response := utils.ErrorResponse(400, "Error: Enter a valid Phone Number. Phone Number format is incorrect", err.Error(), signUp_user)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := cr.userUseCase.FindByEmail(c.Request.Context(), signUp_user.Email); err == nil {
		response := utils.ErrorResponse(401, "Error: User with the email already exits!", err.Error(), signUp_user)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	fmt.Println("config variable", config.GetCofig())

	signUp_user.Password, _ = support.HashPassword(signUp_user.Password)
	PhoneNum, err := cr.userUseCase.SignUpUser(c.Request.Context(), signUp_user)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to Add user, please try again", err.Error(), signUp_user)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: User data saved", signUp_user)
	c.JSON(http.StatusOK, response)

	respSid, err1 := cr.otpUseCase.TwilioSendOTP(c.Request.Context(), PhoneNum)

	if err1 != nil {
		response := utils.ErrorResponse(500, "Error: Failed to generate OTP!", err.Error(), signUp_user)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response1 := utils.SuccessResponse(200, "Success: Enter the otp and the response id", respSid)
	c.JSON(http.StatusOK, response1)
}

// USER SIGN-UP WITH VERIFICATION OF OTP
// @Summary API FOR NEW USER SIGN UP OTP VERIFICATION
// @ID SIGNUP-USER-OTP-VERIFY
// @Description VERIFY THE OTP AND UPDATE THE VERIFIED COLUMN
// @Tags USER
// @Accept json
// @Produce json
// @Param otp_details body utils.OtpSignUpVerify true "otp"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/signup/otp/verify [post]
func (cr *UserHandler) SignupOtpverify(c *gin.Context) {
	var SignUpOtpverify utils.OtpSignUpVerify
	if err := c.BindJSON(&SignUpOtpverify); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), SignUpOtpverify)
		c.JSON(http.StatusBadRequest, response)
	}
	var signUp_user domain.Users
	var otp utils.OtpVerify

	copier.Copy(&otp, &SignUpOtpverify)

	session, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), otp)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Verification failed", err.Error(), signUp_user)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	err1 := cr.userUseCase.UpdateVerify(c.Request.Context(), session.PhoneNum)
	if err1 != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update the verification status of user", err.Error(), signUp_user)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response1 := utils.SuccessResponse(200, "Success: Users Phone Number Successfully verified", session.PhoneNum)
	c.JSON(http.StatusOK, response1)
}

// USER USERLOGIN
// @Summary API FOR USER LOGIN
// @ID USER-LOGIN
// @Description VERIFY THE EMAIL,PASSWORD, HASH THE PASSWORD AND GENERATE A JWT TOKEN AND SET IT TO A COOKIE
// @Tags USER
// @Accept json
// @Produce json
// @Param login_details body utils.LoginBody true "Enter the email and password"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/login [post]
func (cr *UserHandler) LoginHandler(c *gin.Context) {
	//Cookie check
	_, err1 := c.Cookie("user-token")
	if err1 == nil {
		c.Redirect(http.StatusFound, "/user/profile/home")
		// c.AbortWithStatusJSON(http.StatusFound, gin.H{
		// 	"alert": "User already logged in and cookie present",
		// })
		return
	}
	//Login logic
	var loginBody utils.LoginBody
	if err := c.BindJSON(&loginBody); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), loginBody)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//Checks whether such user email exits or not and also returns back the user details of that specific user related to the email and stores in user.
	user, err := cr.userUseCase.FindByEmail(c.Request.Context(), loginBody.Email)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Please check the email", err.Error(), loginBody)
		c.JSON(http.StatusUnauthorized, response)
		return

	}
	//Checks the given password with retreived password to that specific email from the database(user variable)
	if err := support.CheckPasswordHash(loginBody.Password, user.Password); err != nil {
		response := utils.ErrorResponse(401, "Error: Please check the password", err.Error(), loginBody)
		c.JSON(http.StatusUnauthorized, response)
		return

	}

	//GenerateJWT function from the auth package, passing user.Email and User.ID as an argument. It assigns the generated JWT to the tokenString variable
	tokenString, err := auth.GenerateJWT(user.Email, user.ID)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Error generating the jwt token", err.Error(), loginBody)
		c.JSON(http.StatusInternalServerError, response)
		return

	}

	//Sets a cookie named "user-token" with the value tokenString. The cookie has an expiration time of 60 minutes from the current time.
	c.SetCookie("user-token", tokenString, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	c.Set("user-email", user.Email)
	response1 := utils.SuccessResponse(200, "Success: Login Successful")
	c.JSON(http.StatusOK, response1)
}

// USER FORGOT PASSWORD
// @Summary API FOR USER FORGOT PASSWORD OPTION
// @ID USER-FORGOT-PASSWORD
// @Description VERIFY THE EMAIL AND NUMBER AND FIND THE DATA. SEND THE OTP AND VERIFY WITH NEW PASSWORD AND OTP.
// @Tags USER
// @Accept json
// @Produce json
// @Param login_details body utils.OtpLogin true "Enter the email and phoneNumber"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/forgot/password [post]
func (cr *UserHandler) ForgotPassword(c *gin.Context) {
	var body utils.OtpLogin
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	user, err := cr.userUseCase.FindByEmailOrNumber(c.Request.Context(), body)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Incorrect email or password", err.Error(), user)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	respSid, err := cr.otpUseCase.TwilioSendOTP(c.Request.Context(), user.PhoneNum)
	fmt.Println("Send otp")
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to send otp", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("This is the response id", respSid)
	response := utils.SuccessResponse(200, "Success: Successfully sent the otp. Now enter the otp,response id and new password", "", respSid)
	c.JSON(http.StatusOK, response)
}

// USER FORGOT PASSWORD OTP VERIFY
// @Summary API FOR USER FORGOT PASSWORD OTP VERIFICATION
// @ID USER-FORGOT-PASSWORD-OTP-VERIFY
// @Description VERIFY THE OTP AND ENTER A NEW PASSWORD
// @Tags USER
// @Accept json
// @Produce json
// @Param verify_details body utils.OtpVerify true "Enter the Otp and New Password"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/forgot/password/otp/verify [patch]
func (cr *UserHandler) ForgotPasswordOtpVerify(c *gin.Context) {
	var body utils.OtpVerify
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	session, err := cr.otpUseCase.TwilioVerifyOTP(c.Request.Context(), body)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to verify otp", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	NewHashedPassword, err := support.HashPassword(body.NewPassword)
	if err != nil {
		response := utils.ErrorResponse(500, "Error:  Failed to Hash New password", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	if err := cr.userUseCase.ChangePassword(c.Request.Context(), NewHashedPassword, session.PhoneNum); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update new password", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully updated new password", nil)
	c.JSON(http.StatusOK, response)

}

// USERLOGOUT
// @Summary API FOR USER LOGOUT
// @ID USER-LOGOUT
// @Description LOGOUT USER AND ALSO CLEAR COOKIES
// @Tags USER
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/logout [post]
func (cr *UserHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("user-token", "", -1, "/", "localhost", false, true)
	response := utils.SuccessResponse(200, "Success: Logout Successful", nil)
	c.JSON(http.StatusOK, response)
}

// USER PROFILE
// @Summary API FOR USER PROFILE
// @ID USER-PROFILE
// @Description DISPLAY USER PROFILE
// @Tags USER
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile/home [get]
func (cr *UserHandler) HomeHandler(c *gin.Context) {
	// email, ok := c.Get(("user-email"))
	// if !ok {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"error": "Invalid user",
	// 	})
	// }
	// user, err := cr.userUseCase.FindByEmail(c.Request.Context(), email.(string))
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"error": "Invalid user",
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusOK, user)
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	user, err := cr.userUseCase.HomeHandler(c.Request.Context(), id.(uint))
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to fetch user details", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Personal details", user)
	c.JSON(http.StatusOK, response)

}

// EDIT PROFILE
// @Summary API FOR EDIT PROFILE
// @ID USER-PROFILE EDIT
// @Description EDIT/UPDATE USER PROFILE
// @Tags USER
// @Accept json
// @Produce json
// @Param update_details body utils.UpdateProfile true "Edit the details as per wish"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile/edit/profile [patch]
func (cr *UserHandler) UpdateProfile(c *gin.Context) {
	id, ok := c.Get("user-id")
	var body utils.UpdateProfile
	if !ok {
		reponse := utils.ErrorResponse(401, "Error: failed to get id from token strin", "", nil)
		c.JSON(http.StatusUnauthorized, reponse)
		return
	}
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind Json", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var updateProfile domain.Users
	copier.Copy(&updateProfile, &body)
	if err := cr.userUseCase.UpdateProfile(c.Request.Context(), updateProfile, id.(uint)); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update profile", err.Error(), updateProfile)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully updated profile details", nil)
	c.JSON(http.StatusOK, response)
}

// ListProduct
// func (cr *UserHandler) ListProducts(c *gin.Context) {
// 	products, err := cr.userUseCase.ListProducts(c.Request.Context())

// 	if err != nil {
// 		response := utils.ErrorResponse(500, "Failed to list products", err.Error(), nil)
// 		c.JSON(http.StatusInternalServerError, response)
// 		return
// 	}
// 	response := utils.SuccessResponse(200, "All product details", products)
// 	c.JSON(http.StatusOK, response)
// }

// ADD ADDRESS
// @Summary API FOR ADDING ADDRESS
// @ID USER-ADD-ADDRESS
// @Description ADDING NEW ADDRESS TO USER PROFILE
// @Tags USER
// @Accept json
// @Produce json
// @Param address_details body utils.Address true "Add the address details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile/add/address [post]
func (cr *UserHandler) AddAddress(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	var body utils.Address
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind JSON body", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	body.UserID = id.(uint)
	var address domain.Address
	copier.Copy(&address, &body)
	if err := cr.userUseCase.AddAddress(c.Request.Context(), address); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to add address", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully added the address", nil)
	c.JSON(http.StatusOK, response)

}

// LIST ADDRESS
// @Summary API FOR LISTING ADDRESSES
// @ID USER-LIST-ADDRESS
// @Description LISTING ALL ADDRESSES FOR THE PARTICULAR USER
// @Tags USER
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile/list/address [get]
func (cr *UserHandler) ListAddress(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1 // Default page number is 1 if not provided
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 3 // Default limit is 3 items if not provided or if an invalid value is entered
	}

	offset := (page - 1) * limit
	pagination := utils.Pagination{
		Offset: uint(offset),
		Limit:  uint(limit),
	}
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	address, err := cr.userUseCase.ListAddress(c.Request.Context(), id.(uint), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: failed to retreive the addresses", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success:", address)
	c.JSON(http.StatusOK, response)
}

// EDIT ADDRESS
// @Summary API FOR EDITING/UPDATING ADDRESS
// @ID USER-EDIT-ADDRESS
// @Description EDITING EXISTING ADDRESS ON USER PROFILE
// @Tags USER
// @Accept json
// @Produce json
// @Param address_id path string true "Enter the address id that you need to update"
// @Param address_details body utils.UpdateAddress true "edit the address details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile/edit/address/{addressid} [patch]
func (cr *UserHandler) UpdateAddress(c *gin.Context) {
	id := c.Param("addressid")
	var body utils.UpdateAddress
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind Json", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// body.UserID = id.(uint)
	var updateAddress domain.Address
	copier.Copy(&updateAddress, &body)
	if err := cr.userUseCase.UpdateAddress(c.Request.Context(), updateAddress, id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update the address", err.Error(), updateAddress)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Updated successfully", updateAddress)
	c.JSON(http.StatusOK, response)
}

// DELETE ADDRESS
// @Summary API FOR DELETING ADDRESS
// @ID USER-DELETE-ADDRESS
// @Description DELETING EXISTING ADDRESS ON USER PROFILE
// @Tags USER
// @Accept json
// @Produce json
// @Param address_id path string true "Enter the address id that you need to delete"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile/delete/address/{addressid} [post]
func (cr *UserHandler) DeleteAddress(c *gin.Context) {
	id := c.Param("addressid")
	if err := cr.userUseCase.DeleteAddress(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to delete the address", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully deleted the address", nil)
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
