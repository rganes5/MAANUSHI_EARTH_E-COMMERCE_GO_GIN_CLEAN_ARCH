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
		response := utils.ErrorResponse(500, "Faild to add cateogy", err.Error(), category)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "successfully added a new category", category)
	c.JSON(http.StatusOK, response)
}

// UpdateCategory
func (cr *ProductHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("categoryid")
	var categoryname utils.UpdateCategory
	if err := c.BindJSON(&categoryname); err != nil {
		response := utils.ErrorResponse(400, "Failed to bind JSON", err.Error(), categoryname)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var categories domain.Category
	copier.Copy(&categories, &categoryname)
	if err := cr.productUseCase.UpdateCategory(c.Request.Context(), categories, id); err != nil {
		response := utils.ErrorResponse(500, "Failed to update category", err.Error(), categories)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Successfully updated", categories)
	c.JSON(http.StatusOK, response)
}

// DELETE
func (cr *ProductHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("categoryid")
	if err := cr.productUseCase.DeleteCategory(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Faild to Delete Cateogy with id", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "successfully deleted category with id", id)
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
		response := utils.ErrorResponse(500, "Faild to List Categories", err.Error(), categories)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"categories": categories,
	// })
	response := utils.SuccessResponse(200, "List of categories:", categories)
	c.JSON(http.StatusOK, response)
}

// PRODUCT MANAGEMENT
// ADD
func (cr *ProductHandler) AddProduct(c *gin.Context) {
	var body utils.Products
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Failed to bind json", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var products domain.Products
	copier.Copy(&products, &body)
	if err := cr.productUseCase.AddProduct(c.Request.Context(), products); err != nil {
		response := utils.ErrorResponse(500, "Failed to add product", err.Error(), products)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Product Added Successfully", products)
	c.JSON(http.StatusOK, response)

}

// DELETE
func (cr *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("productid")
	if err := cr.productUseCase.DeleteProduct(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Failed to delete the product", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Product successfully deleted with product id", id)
	c.JSON(http.StatusOK, response)
}

// EDIT PRODUCTS
func (cr *ProductHandler) EditProduct(c *gin.Context) {
	// var product domain.Products
	var body utils.Products
	id := c.Param("productid")
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, " Error while binding json", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var product domain.Products
	copier.Copy(&product, &body)
	err := cr.productUseCase.EditProduct(c.Request.Context(), product, id)
	if err != nil {
		response := utils.ErrorResponse(500, "Error while editing the product", err.Error(), product)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Successfully edited the product", product)
	c.JSON(http.StatusOK, response)
}

// LIST PRODUCTS
func (cr *ProductHandler) ListProducts(c *gin.Context) {
	products, err := cr.productUseCase.ListProducts(c.Request.Context())
	if err != nil {
		response := utils.ErrorResponse(500, "Failed to list the products", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Product list", products)
	c.JSON(http.StatusOK, response)

}
