package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/checkout/placeorder [post]
func (cr *OrderHandler) PlaceNewOrder(c *gin.Context) {
	paymentId, _ := strconv.Atoi(c.Query("payment_id"))
	addressId, _ := strconv.Atoi(c.Query("address_id"))
	userId, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	if paymentId == 1 {
		fmt.Println("COD")
		fmt.Println("the address id and user id, payment id and user id is from handler are", uint(addressId), uint(paymentId), userId.(uint))
		if err := cr.orderUseCase.PlaceNewOrder(c.Request.Context(), uint(addressId), uint(paymentId), userId.(uint)); err != nil {
			response := utils.ErrorResponse(500, "Error: Failed to add to place order", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	} else if paymentId == 2 {
		fmt.Println("Razorpay not added yet")
		response := utils.ErrorResponse(400, "Error: Razorpay payment method is not yet available", "", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		fmt.Println("Please select correct payment method")
		response := utils.ErrorResponse(400, "Error: Please select correct payment method", "", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully placed the order", nil)
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
// @Param order_id query uint true "Enter the order id"
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/orders/cancel/{order_id} [post]
func (cr *OrderHandler) CancelOrder(c *gin.Context) {
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	fmt.Println("id is", id)
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
