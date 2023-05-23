package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type CartHandler struct {
	cartUseCase services.CartUseCase
}

func NewCartHandler(service services.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: service,
	}
}

func (cr *CartHandler) AddToCart(c *gin.Context) {
	productId := c.Param("id")
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	if err := cr.cartUseCase.AddToCart(c.Request.Context(), productId, id.(uint)); err != nil {
		response := utils.ErrorResponse(500, "Failed to add the item to cart", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Successfully added the item to cart", nil)
	c.JSON(http.StatusOK, response)
}
