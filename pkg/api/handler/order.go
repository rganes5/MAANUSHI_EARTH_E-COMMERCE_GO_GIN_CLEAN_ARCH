package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/support"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type OrderHandler struct {
	orderUseCase services.OrderUseCase
}

func NewOrderHandler(service services.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: service,
	}
}

// USERS END
//
// PLACE A NEW ORDER
// @Summary API FOR PLACING A NEW ORDER
// @ID USER-PROCEED-ORDER
// @Description Users can place a new order with the cart items.
// @Tags ORDER
// @Accept json
// @Produce json
// @Param payment_id query string true "Enter the payment id"
// @Param address_id query string true "Enter the address id"
// @Param code query string false "Enter coupon code if available"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/checkout/placeorder [get]
func (cr *OrderHandler) PlaceNewOrder(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	paymentId, _ := strconv.Atoi(c.Query("payment_id"))
	addressId, _ := strconv.Atoi(c.Query("address_id"))
	userId, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	fmt.Println("The coupon code is", code)

	// var couponid *uint
	// var err error

	couponid, err := cr.orderUseCase.ValidateCoupon(c.Request.Context(), userId.(uint), code)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to read the coupon code", err.Error(), nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	fmt.Println("The address ID, user ID, and payment ID are", uint(addressId), uint(paymentId), userId.(uint))

	// Check the payment modes
	if paymentId == 1 || paymentId == 3 {
		if err := cr.orderUseCase.PlaceNewOrder(c.Request.Context(), uint(addressId), uint(paymentId), userId.(uint), couponid); err != nil {
			response := utils.ErrorResponse(500, "Error: Failed to place order", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := utils.SuccessResponse(200, "Success: Successfully placed the order", nil)
		c.JSON(http.StatusOK, response)

	} else if paymentId == 2 {
		body, err := cr.orderUseCase.RazorPayOrder(c.Request.Context(), userId.(uint), couponid)
		fmt.Println("body from the razorpay new order is", body)
		if err != nil {
			response := utils.ErrorResponse(500, "Error: Failed to load a razorpay payment", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		fmt.Println("loads to here", body.RazorpayOrderID, body.UserID, body.AmountToPay)
		c.HTML(200, "app.html", gin.H{
			"UserID":      body.UserID,
			"Orderid":     body.RazorpayOrderID,
			"Total_price": body.AmountToPay,
		})
	} else if paymentId == 5 {
		fmt.Println("Paypal not added yet")
		response := utils.ErrorResponse(400, "Error: Paypal not added yet", "", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
}

/*
// Razorpay success verification handler
//
// @Summary API FOR VERIFYING THE RAZORPAY STATUS
// @Description For checking and verifying the payment
// @Tags ORDER
// @Accept json
// @Produce json
// @Param razorpay_order_id query string true "Enter the razorpay_order_id"
// @Param razorpay_payment_id query string true "Enter the razorpay_payment_id"
// @Param razorpay_signature query string true "Enter the razorpay_signature"
// @Param payment_id query string true "Enter the payment id"
// @Param address_id query string true "Enter the address id"
// @Param code query string false "Enter coupon code if available"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/checkout/success [post]
*/
func (cr *OrderHandler) RazorPaySuccess(c *gin.Context) {
	userId, err3 := strconv.Atoi(c.Query("user_id"))
	razorpay_order_id := c.Query("order_id")
	razorpay_payment_id := c.Query("payment_ref")
	razorpay_signature := c.Query("signature")
	addressId, err2 := strconv.Atoi(c.Query("address_id"))
	paymentId, err1 := strconv.Atoi(c.Query("payment_id"))
	code := c.Query("code")
	fmt.Println("The coupon from the verificatin is ", code)
	response := gin.H{
		"data":    false,
		"message": "Payment failed",
	}
	// userId, ok := c.Get("user-id")
	// if !ok {
	// 	response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
	// 	c.JSON(http.StatusUnauthorized, response)
	// 	return
	// }
	err := errors.Join(err1, err2, err3)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to get the id's", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fmt.Println("Passed errors")

	fmt.Println("Here are the data from razorpaysuccess handler", razorpay_order_id, razorpay_payment_id, razorpay_signature)
	if err := support.VerifyRazorPayment(razorpay_order_id, razorpay_payment_id, razorpay_signature); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to verify the payment", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("Passed verification")

	couponid, er := cr.orderUseCase.FindCoupon(c.Request.Context(), code)
	fmt.Println("error from the find coupon is ", er)
	if er != nil {
		response := utils.ErrorResponse(500, "Error: Failed to find coupon", er.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("Passed coupon check")

	if err := cr.orderUseCase.PlaceNewOrder(c.Request.Context(), uint(addressId), uint(paymentId), uint(userId), couponid); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to place order", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// response := utils.SuccessResponse(200, "Success:Payment verified and order placed successfully.", nil)
	// c.JSON(http.StatusOK, response)
	response["data"] = true
	response["message"] = "Payment verified and order placed successfully."
	c.JSON(http.StatusOK, response)
}

// VIEW ORDERS
// @Summary API FOR VIEWING ORDERS
// @Description Users can view all orders.
// @Tags ORDER
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/orders/list/all [get]
func (cr *OrderHandler) ListOrders(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
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
	orderList, err := cr.orderUseCase.ListOrders(c.Request.Context(), id.(uint), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list order list", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Order list", orderList)
	c.JSON(http.StatusOK, response)

}

// VIEW ORDERS DETAILS
// @Summary API FOR VIEWING ORDERS DETAILS
// @Description Users can the selected order details.
// @Tags ORDER
// @Accept json
// @Produce json
// @Param order_id query uint true "Enter the order id"
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/orders/list/details/{order_id} [get]
// @Router /admin/orders/list/details/{order_id} [get]
func (cr *OrderHandler) ListOrderDetails(c *gin.Context) {
	fmt.Println("1")
	orderId, err := strconv.Atoi(c.Query("order_id"))
	fmt.Println("order id is", orderId)

	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to convert id", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
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
	orderDetails, err := cr.orderUseCase.ListOrderDetails(c.Request.Context(), uint(orderId), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list order details list", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Order list", orderDetails)
	c.JSON(http.StatusOK, response)
}

// CANCEL ORDER
// @Summary API FOR CANCELLING A ORDER
// @Description Users can cancel orders
// @Tags ORDER
// @Accept json
// @Produce json
// @Param order_details_id query uint true "Enter the order details id"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/orders/cancel/{order_details_id} [post]
func (cr *OrderHandler) CancelOrder(c *gin.Context) {
	userId, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	fmt.Println("id is", userId)
	orderDetailsId, err := strconv.Atoi(c.Query("order_details_id"))
	fmt.Println("order id is", orderDetailsId)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to convert id", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := cr.orderUseCase.CancelOrder(c.Request.Context(), userId.(uint), uint(orderDetailsId)); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to cancel the order with order id", err.Error(), orderDetailsId)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully cancelled the order with order id", orderDetailsId)
	c.JSON(http.StatusOK, response)
}

// RETURN ORDER
// @Summary API FOR RETURNING A ORDER
// @Description Users can return orders
// @Tags ORDER
// @Accept json
// @Produce json
// @Param order_details_id query uint true "Enter the order details id"
// @Param status_id query uint true "Enter the order status id"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/orders/return/{order_details_id} [post]
func (cr *OrderHandler) ReturnOrder(c *gin.Context) {
	orderDetailsId, err1 := strconv.Atoi(c.Query("order_details_id"))
	statusId, err2 := strconv.Atoi(c.Query("status_id"))
	err := errors.Join(err1, err2)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to convert id'ss", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := cr.orderUseCase.ReturnOrder(c.Request.Context(), uint(orderDetailsId), uint(statusId)); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to submit the return request", err.Error(), orderDetailsId)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully requested a return", orderDetailsId)
	c.JSON(http.StatusOK, response)

}

// ADMINS END
//
// VIEW ORDERS
// @Summary API FOR VIEWING ORDERS
// @Description Admin can view all orders.
// @Tags ORDER
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/orders/list/all [get]
func (cr *OrderHandler) AdminListOrders(c *gin.Context) {

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
	orderList, err := cr.orderUseCase.AdminListOrders(c.Request.Context(), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list order list", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Order list", orderList)
	c.JSON(http.StatusOK, response)

}

// CHANGE STATUS OF ORDER
// @Summary API FOR CHANGING THE STATUS OF A ORDER
// @Description Admin can change the ststus of orders
// @Tags ORDER
// @Accept json
// @Produce json
// @Param order_details_id query uint true "Enter the order details id"
// @Param status_id query uint true "Enter the order status id"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/orders/update/{order_details_id} [post]
func (cr *OrderHandler) UpdateStatus(c *gin.Context) {
	orderDetailsId, err1 := strconv.Atoi(c.Query("order_details_id"))
	statusId, err2 := strconv.Atoi(c.Query("status_id"))
	err := errors.Join(err1, err2)
	if err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to convert id'ss", err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err := cr.orderUseCase.UpdateStatus(c.Request.Context(), uint(orderDetailsId), uint(statusId)); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update the status of order id", err.Error(), orderDetailsId)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully updated the status of the order with order id", orderDetailsId)
	c.JSON(http.StatusOK, response)

}
