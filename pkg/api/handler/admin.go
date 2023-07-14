package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/auth"
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

// ADMIN SIGN-UP WITHOUT SENDING OTP
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
	// var signUp_admin domain.Admin
	// copier.Copy(&signUp_admin, &body)

	//Check the email format
	if err := support.Email_validator(body.Email); err != nil {
		response := utils.ErrorResponse(400, "Error: Enter a valid email. Email format is incorrect", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return

	}

	//Check the phone number format
	if err := support.MobileNum_validator(body.PhoneNum); err != nil {
		response := utils.ErrorResponse(400, "Error: Enter a valid Phone Number. Phone Number format is incorrect", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//Check whether such email already exits
	if _, err := cr.adminUseCase.FindByEmail(c.Request.Context(), body.Email); err == nil {
		response := utils.ErrorResponse(401, "Error: Admin with the email already exits!", err.Error(), body)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	//Hash the password and sign up
	body.Password, _ = support.HashPassword(body.Password)
	if _, err := cr.adminUseCase.SignUpAdmin(c.Request.Context(), body); err != nil {
		response := utils.ErrorResponse(401, "Error: Failed to Add Admin, please try again", err.Error(), body)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.SuccessResponse(200, "success: Admin sign-up Successful", body)
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
	c.SetCookie("admin-token", tokenstring, int(time.Now().Add(60*time.Minute).Unix()), "/", "maanushiearth.shop", true, false)
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
		response := utils.ErrorResponse(500, "Error: Failed to update the access", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: User access updated", access)
	c.JSON(http.StatusOK, response)

}

// ADMIN DASHBOARD WITH WIDGETS
// @Summary API FOR LISTING WIDGETS
// @ID ADMIN-LIST-WIDGETS
// @Description ADMIN DASHBOARD AND LISTING WIDGETS FOR ADMIN
// @Tags ADMIN
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/dashboard [get]
func (cr *AdminHandler) Dashboard(c *gin.Context) {
	responseWidgets, err := cr.adminUseCase.Dashboard(c.Request.Context())
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to fetch the data", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Here are the widgets", responseWidgets)
	c.JSON(http.StatusOK, response)
}

// ADMIN SALES REPORT
// @Summary API FOR GETTING SALES REPORT
// @ID ADMIN-SALES-REPORT
// @Description ADMIN SALES REPORT, VIA MONTHLY AND YEARLY
// @Tags ADMIN
// @Accept json
// @Produce json
// @Param frequency query string false "Enter frequency"
// @Param month query int false "Enter the month"
// @Param year query int false "Enter the year"
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/dashboard/salesreport [get]
func (cr *AdminHandler) SalesReport(c *gin.Context) {
	// time
	monthInt, err1 := strconv.Atoi(c.DefaultQuery("month", "1"))
	month := time.Month(monthInt)
	year, err2 := strconv.Atoi(c.Query("year"))
	frequency := c.Query("frequency")
	err := errors.Join(err1, err2)
	if err != nil {
		response := utils.ErrorResponse(400, "Failed to parse the month and year", err1.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	page, err3 := strconv.Atoi(c.Query("page"))
	if err3 != nil {
		page = 1 // Default page number is 1 if not provided
	}

	limit, err4 := strconv.Atoi(c.Query("limit"))
	if err4 != nil || limit <= 0 {
		limit = 3 // Default limit is 3 items if not provided or if an invalid value is entered
	}

	offset := (page - 1) * limit
	reqData := utils.SalesReport{
		Month:     month,
		Year:      year,
		Frequency: frequency,
		Pagination: utils.Pagination{
			Offset: uint(offset),
			Limit:  uint(limit),
		},
	}
	salesreport, err5 := cr.adminUseCase.SalesReport(reqData)
	if err5 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err5.Error(),
		})
		return
	}
	if salesreport == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "there is no sales report for this selected period",
		})
	} else {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment;filename=maanushi_earth_e-commerce_salesreport.csv")

		csvWriter := csv.NewWriter(c.Writer)
		headers := []string{
			"UserID", "FirstName", "Email",
			"ProductDetailID", "ProductName", "Price",
			"DiscountPrice", "Quantity", "OrderID",
			"PlacedDate", "PaymentMode", "OrderStatus", "Total",
		}

		if err := csvWriter.Write(headers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		grandtotal := 0
		for _, sales := range salesreport {
			total := sales.Quantity * sales.DiscountPrice
			fmt.Println("This is the total", total)
			row := []string{
				fmt.Sprintf("%v", sales.UserID),
				sales.FirstName,
				sales.Email,
				fmt.Sprintf("%v", sales.ProductDetailID),
				sales.ProductName,
				fmt.Sprintf("%v", sales.Price),
				fmt.Sprintf("%v", sales.DiscountPrice),
				fmt.Sprintf("%v", sales.Quantity),
				fmt.Sprintf("%v", sales.OrderID),
				sales.PlacedDate.Format("2006-01-02 15:04:05"),
				sales.PaymentMode,
				sales.OrderStatus,
				fmt.Sprintf("%v", total),
			}

			if err := csvWriter.Write(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			fmt.Println("sales payment mode is", sales.PaymentMode)
			fmt.Println("sales status  is", sales.OrderStatus)

			if sales.PaymentMode == "Razorpay" || sales.PaymentMode == "Wallet" || (sales.PaymentMode == "cash on delivery" && sales.OrderStatus == "delivered") {
				grandtotal += int(total)
				fmt.Println("This is the grandtotal", grandtotal)
			}
		}
		rowtotal := []string{
			fmt.Sprintf("Grand Total=%v", grandtotal),
		}
		if err := csvWriter.Write(rowtotal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		csvWriter.Flush()
	}
}

// ADDING COUPON
// @Summary API FOR ADDING COUPON
// @ID ADMIN-ADD-COUPON
// @Description ADDING COUPON FROM ADMINS END
// @Tags COUPON
// @Accept json
// @Produce json
// @Param couponBody body utils.BodyAddCoupon false "Enter the coupon details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/coupon/add [post]
func (cr *AdminHandler) AddCoupon(c *gin.Context) {
	var couponBody utils.BodyAddCoupon
	if err := c.BindJSON(&couponBody); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json", err.Error(), couponBody)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := cr.adminUseCase.AddCoupon(c.Request.Context(), couponBody); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to add coupon", err.Error(), couponBody)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully added the coupon", couponBody)
	c.JSON(http.StatusOK, response)
}

// LIST COUPONS
// @Summary API FOR LISTING COUPONS
// @Description LISTING ALL COUPOUNS
// @Tags COUPON
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/coupon/list [get]
func (cr *AdminHandler) GetAllCoupons(c *gin.Context) {
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
	coupons, err := cr.adminUseCase.GetAllCoupons(c.Request.Context(), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list coupons", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: ", coupons)
	c.JSON(http.StatusOK, response)
}

// EDITING COUPON
// @Summary API FOR UPDATING COUPON
// @Description UPDATING COUPON FROM ADMINS END
// @Tags COUPON
// @Accept json
// @Produce json
// @Param coupon_id query string	true "Enter the coupon id to update"
// @Param couponBody body utils.BodyAddCoupon false "Enter the coupon details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/coupon/update [patch]
func (cr *AdminHandler) UpdateCoupon(c *gin.Context) {
	couponId := c.Query("coupon_id")
	var couponBody utils.BodyAddCoupon
	if err := c.BindJSON(&couponBody); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json", err.Error(), couponBody)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := cr.adminUseCase.UpdateCoupon(c.Request.Context(), couponBody, couponId); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update coupon", err.Error(), couponBody)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully updated the coupon", couponBody)
	c.JSON(http.StatusOK, response)
}

// DELETING COUPON
// @Summary API FOR DELETING COUPON
// @Description DELETING COUPON FROM ADMINS END
// @Tags COUPON
// @Accept json
// @Produce json
// @Param coupon_id query string	true "Enter the coupon id to delete"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/coupon/delete [delete]
func (cr *AdminHandler) DeleteCoupon(c *gin.Context) {
	couponId := c.Query("coupon_id")
	fmt.Println("this is the id", couponId)
	if err := cr.adminUseCase.DeleteCoupon(c.Request.Context(), couponId); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to delete coupon", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully deleted the coupon with id", couponId)
	c.JSON(http.StatusOK, response)
}
