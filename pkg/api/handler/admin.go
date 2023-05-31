package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/auth"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/support"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

// ADMIN SIGN-UP WITH SENDING OTP
// @Summary API FOR NEW USER SIGN UP
// @ID SIGNUP-ADMIN
// @Description CREATE A NEW ADMIN WITH REQUIRED DETAILS
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param admin_details body utils.AdminSignUp true "New Admin Details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/signup [post]
func (cr *AdminHandler) AdminSignUp(c *gin.Context) {
	var body utils.AdminSignUp

	//Binding
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var signUp_admin domain.Admin
	copier.Copy(&signUp_admin, &body)

	//Check the email format
	if err := support.Email_validator(signUp_admin.Email); err != nil {
		response := utils.ErrorResponse(400, "Error: Enter a valid email. Email format is incorrect", err.Error(), signUp_admin)
		c.JSON(http.StatusBadRequest, response)
		return

	}

	//Check the phone number format
	if err := support.MobileNum_validator(signUp_admin.PhoneNum); err != nil {
		response := utils.ErrorResponse(400, "Error: Enter a valid Phone Number. Phone Number format is incorrect", err.Error(), signUp_admin)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//Check whether such email already exits
	if _, err := cr.adminUseCase.FindByEmail(c.Request.Context(), signUp_admin.Email); err == nil {
		response := utils.ErrorResponse(401, "Error: Admin with the email already exits!", err.Error(), signUp_admin)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	//Hash the password and sign up
	signUp_admin.Password, _ = support.HashPassword(signUp_admin.Password)
	if err := cr.adminUseCase.SignUpAdmin(c.Request.Context(), signUp_admin); err != nil {
		response := utils.ErrorResponse(401, "Error: Failed to Add Admin, please try again", err.Error(), signUp_admin)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.SuccessResponse(200, "success: Admin sign-up Successful", signUp_admin)
	c.JSON(http.StatusOK, response)
}

// ADMIN LOGIN
// @Summary API FOR ADMIN LOGIN
// @ID ADMIN-LOGIN
// @Description VERIFY THE EMAIL,PASSWORD, HASH THE PASSWORD AND GENERATE A JWT TOKEN AND SET IT TO A COOKIE
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param login_details body utils.LoginBody true "Enter the email and password"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/login [post]
func (cr *AdminHandler) AdminLogin(c *gin.Context) {
	//Checks whether a jwt token already exists
	_, err := c.Cookie("admin-token")
	if err == nil {
		// c.JSON(http.StatusAlreadyReported, gin.H{
		// 	"admin": "Already logged in",
		// })
		c.Redirect(http.StatusFound, "/admin/home")
		return
	}
	//binding
	var Login_admin utils.LoginBody
	if err := c.BindJSON(&Login_admin); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), Login_admin)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//check email
	admin, err := cr.adminUseCase.FindByEmail(c.Request.Context(), Login_admin.Email)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Please check the email", err.Error(), Login_admin)
		c.JSON(http.StatusUnauthorized, response)
	}

	//check the password
	if err := support.CheckPasswordHash(Login_admin.Password, admin.Password); err != nil {
		response := utils.ErrorResponse(401, "Error: Please check the password", err.Error(), Login_admin)
		c.JSON(http.StatusUnauthorized, response)
	}

	//Create a jwt token and store it in cookie
	tokenstring, err := auth.GenerateJWT(admin.Email, admin.ID)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Error generating the jwt token", err.Error(), Login_admin)
		c.JSON(http.StatusInternalServerError, response)
	}
	c.SetCookie("admin-token", tokenstring, int(time.Now().Add(60*time.Minute).Unix()), "/", "localhost", false, true)
	response1 := utils.SuccessResponse(200, "Success: Login Successful")
	c.JSON(http.StatusOK, response1)
}

// ADMIN LOGOUT
// @Summary API FOR ADMIN LOGOUT
// @ID ADMIN-LOGOUT
// @Description ADMIN LOGOUT
// @Tags ADMIN
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/logout [post]
func (cr *AdminHandler) Logout(c *gin.Context) {
	c.SetCookie("admin-token", "", -1, "/", "localhost", false, true)
	response := utils.SuccessResponse(200, "Success: Logout Successful", nil)
	c.JSON(http.StatusOK, response)

}

// ADMIN PROFILE
// @Summary API FOR ADMIN PROFILE
// @ID ADMIN-PROFILE
// @Description DISPLAY ADMIN PROFILE
// @Tags ADMIN
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/home [get]
func (cr *AdminHandler) HomeHandler(c *gin.Context) {
	email, ok := c.Get(("admin-email"))
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	admin, err := cr.adminUseCase.FindByEmail(c.Request.Context(), email.(string))
	if err != nil {
		response := utils.ErrorResponse(400, "Error:Failed to fetch user details", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Personal details", admin)
	c.JSON(http.StatusOK, response)
}

// LIST USERS
// @Summary API FOR LISTING USERS
// @ID ADMIN-LIST-USERS
// @Description LISTING ALL EXISTING USERS
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/user [get]
func (cr *AdminHandler) ListUsers(c *gin.Context) {
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
	users, err := cr.adminUseCase.ListUsers(c.Request.Context(), pagination)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Failed to List users", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: ", users)
	c.JSON(http.StatusOK, response)
}

// ACCESS HANDLER
// @Summary API FOR BLOCKING/UNBLOCKING USERS
// @ID ADMIN-ACCESS
// @Description GRANTING ACCESS FOR INDIVIDUAL USERS.
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param user_id path string true "Enter the specific user id"
// @Param access query string false "Enter true/false"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/user/{user_id}/make [patch]
func (cr *AdminHandler) AccessHandler(c *gin.Context) {
	id := c.Param("user_id")
	str := c.Query("access")
	access, err1 := strconv.ParseBool(str)
	if err1 != nil {
		response := utils.ErrorResponse(400, "Failed to parse the access query", err1.Error(), str)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err := cr.adminUseCase.AccessHandler(c.Request.Context(), id, access)
	if err != nil {
		response := utils.ErrorResponse(401, "Error: Failed to update the access", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: User access updated", access)
	c.JSON(http.StatusOK, response)

}

// // CATEGORY MANAGEMENT

// // ADD
// func (cr *AdminHandler) AddCategory(c *gin.Context) {
// 	var category domain.Category
// 	if err := c.BindJSON(&category); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": "error while binding json",
// 		})
// 		return
// 	}
// 	if err := cr.adminUseCase.AddCategory(c.Request.Context(), category); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"Category": "Successfully added",
// 	})
// }

// // DELETE
// func (cr *AdminHandler) DeleteCategory(c *gin.Context) {
// 	id := c.Param("categoryid")
// 	if err := cr.adminUseCase.DeleteCategory(c.Request.Context(), id); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"category": "Deleted successfully",
// 	})
// }

// // LIST
// func (cr *AdminHandler) ListCategories(c *gin.Context) {
// 	categories, err := cr.adminUseCase.ListCategories(c.Request.Context())
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"categories": categories,
// 	})
// }
