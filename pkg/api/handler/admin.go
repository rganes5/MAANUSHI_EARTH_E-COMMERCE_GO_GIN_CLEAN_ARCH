package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// Variable declared containing type as Admin which is already initialiazed in domain folder.

func (cr *AdminHandler) AdminSignUp(c *gin.Context) {
	var signUp_admin domain.Admin

	//Binding
	if err := c.BindJSON(&signUp_admin); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to read json body", err.Error(), signUp_admin)
		c.JSON(http.StatusBadRequest, response)
		return
	}

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

// admin login
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

//admin logout

func (cr *AdminHandler) Logout(c *gin.Context) {
	c.SetCookie("admin-token", "", -1, "/", "localhost", false, true)
	response := utils.SuccessResponse(200, "Success: Logout Successful", nil)
	c.JSON(http.StatusOK, response)

}

//home handler

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

//list users

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

//Block and unblock

func (cr *AdminHandler) AccessHandler(c *gin.Context) {
	id := c.Param("userid")
	str := c.Query("access")
	access, _ := strconv.ParseBool(str)
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
