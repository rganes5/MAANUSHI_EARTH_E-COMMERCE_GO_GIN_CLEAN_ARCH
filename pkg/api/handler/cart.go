package handler

import (
	"net/http"
	"strconv"

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

// CART MANAGEMENT

// ADD TO CART
// @Summary API FOR ADDING PRODUCTS TO CART BY USER
// @ID USER-ADD-TO-CART
// @Description ADDING ITEMS TO CART FROM USERS END
// @Tags CART
// @Accept json
// @Produce json
// @Param product_id path string true "Enter the product id"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/cart/add/{product_id} [post]
func (cr *CartHandler) AddToCart(c *gin.Context) {
	productId := c.Param("product_id")
	id, ok := c.Get("user-id")
	if !ok {
		response := utils.ErrorResponse(401, "Error: Failed to get the id from the token string", "", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	if err := cr.cartUseCase.AddToCart(c.Request.Context(), productId, id.(uint)); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to add the item to cart", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully added the item to cart", nil)
	c.JSON(http.StatusOK, response)
}

// LIST CART_DETAILS
// @Summary API FOR DISPLAYING CART TO USER
// @ID USER-LIST-CART
// @Description LISTING CART AND ITEMS FROM USERS END
// @Tags CART
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/cart/list [get]
func (cr *CartHandler) ListCart(c *gin.Context) {
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
	grantTotal, viewCart, err := cr.cartUseCase.ListCart(c.Request.Context(), id.(uint), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to load cart", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	responseCart := utils.ResponseFullCart{
		CartItems: viewCart,
		SubTotal:  grantTotal,
	}
	response := utils.SuccessResponse(200, "Success: ", responseCart)
	c.JSON(http.StatusOK, response)
}
