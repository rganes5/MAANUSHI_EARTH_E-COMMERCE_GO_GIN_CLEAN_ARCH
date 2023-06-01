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
