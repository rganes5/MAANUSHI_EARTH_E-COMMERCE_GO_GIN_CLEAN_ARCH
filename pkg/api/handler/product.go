package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
)

type ProductHandler struct {
	productUseCase services.ProductUseCase
}

func NewProductHandler(service services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: service,
	}
}

// PRODUCT MANAGEMENT
// ADD
func (cr *ProductHandler) AddProduct(c *gin.Context) {
	var products domain.Products
	if err := c.BindJSON(&products); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error while binding json",
		})
		return
	}
	if err := cr.productUseCase.AddProduct(c.Request.Context(), products); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Product": "Added Successfully",
	})
}

// DELETE
func (cr *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("productid")
	if err := cr.productUseCase.DeleteProduct(c.Request.Context(), id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Product": "Successfully deleted",
	})
}

// LIST PRODUCTS
// func (cr *ProductHandler) ListProducts(c *gin.Context) {
// 	products, err := cr.productUseCase.ListProducts(c.Request.Context())
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"Products_list": products,
// 	})
// }
