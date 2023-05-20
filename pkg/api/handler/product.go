package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type ProductHandler struct {
	productUseCase services.ProductUseCase
}

func NewProductHandler(service services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: service,
	}
}

// CATEGORY MANAGEMENT

// ADD
func (cr *ProductHandler) AddCategory(c *gin.Context) {
	var category domain.Category
	if err := c.BindJSON(&category); err != nil {
		response := utils.ErrorResponse(400, "invalid input", err.Error(), category)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err := cr.productUseCase.AddCategory(c.Request.Context(), category)
	if err != nil {
		response := utils.ErrorResponse(500, "Faild to add cateogy", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "successfully added a new category")
	c.JSON(http.StatusOK, response)
}

// DELETE
func (cr *ProductHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("categoryid")
	if err := cr.productUseCase.DeleteCategory(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Faild to Delete Cateogy", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "successfully deleted category")
	c.JSON(http.StatusOK, response)
}

// LIST
func (cr *ProductHandler) ListCategories(c *gin.Context) {
	categories, err := cr.productUseCase.ListCategories(c.Request.Context())
	if err != nil {
		// c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		// 	"error": err.Error(),
		// })
		// return
		response := utils.ErrorResponse(500, "Faild to List Cateogy", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"categories": categories,
	// })
	response := utils.SuccessResponse(200, "List of categories", categories)
	c.JSON(http.StatusOK, response)
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

// EDIT PRODUCTS
func (cr *ProductHandler) EditProduct(c *gin.Context) {
	// var product domain.Products
	var body utils.UpdateProducts
	id := c.Param("productid")
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error while binding json",
		})
		return
	}
	var product domain.Products
	copier.Copy(&product, &body)
	err := cr.productUseCase.EditProduct(c.Request.Context(), product, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Product updated",
	})
}
