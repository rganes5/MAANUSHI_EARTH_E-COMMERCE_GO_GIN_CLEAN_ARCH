package handler

import (
	"net/http"
	"strconv"

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
		response := utils.ErrorResponse(400, "Error: invalid input", err.Error(), category)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err := cr.productUseCase.AddCategory(c.Request.Context(), category)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Faild to add cateogy", err.Error(), category)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: successfully added a new category", category)
	c.JSON(http.StatusOK, response)
}

// UpdateCategory
func (cr *ProductHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("categoryid")
	var categoryname utils.UpdateCategory
	if err := c.BindJSON(&categoryname); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind JSON", err.Error(), categoryname)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var categories domain.Category
	copier.Copy(&categories, &categoryname)
	if err := cr.productUseCase.UpdateCategory(c.Request.Context(), categories, id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update category", err.Error(), categories)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success:Successfully updated", categories)
	c.JSON(http.StatusOK, response)
}

// DELETE
func (cr *ProductHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("categoryid")
	if err := cr.productUseCase.DeleteCategory(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Error: Faild to Delete Cateogy with id", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully deleted category with id", id)
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
		response := utils.ErrorResponse(500, "Error: Faild to List Categories", err.Error(), categories)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"categories": categories,
	// })
	response := utils.SuccessResponse(200, "Success: List of categories:", categories)
	c.JSON(http.StatusOK, response)
}

// PRODUCT MANAGEMENT
// ADD
func (cr *ProductHandler) AddProduct(c *gin.Context) {
	var body utils.Products
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var products domain.Products
	copier.Copy(&products, &body)
	if err := cr.productUseCase.AddProduct(c.Request.Context(), products); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to add product", err.Error(), products)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success:Product Added Successfully", products)
	c.JSON(http.StatusOK, response)

}

// DELETE
func (cr *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("productid")
	if err := cr.productUseCase.DeleteProduct(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to delete the product", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product successfully deleted with product id", id)
	c.JSON(http.StatusOK, response)
}

// EDIT PRODUCTS
func (cr *ProductHandler) EditProduct(c *gin.Context) {
	// var product domain.Products
	var body utils.Products
	id := c.Param("productid")
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error:  Error while binding json", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var product domain.Products
	copier.Copy(&product, &body)
	err := cr.productUseCase.EditProduct(c.Request.Context(), product, id)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Error while editing the product", err.Error(), product)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully edited the product", product)
	c.JSON(http.StatusOK, response)
}

// LIST PRODUCTS
func (cr *ProductHandler) ListProducts(c *gin.Context) {
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
	products, err := cr.productUseCase.ListProducts(c.Request.Context(), pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list the products", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product list", products)
	c.JSON(http.StatusOK, response)

}

// ADD PRODUCT DETAILS
func (cr *ProductHandler) AddProductDetails(c *gin.Context) {
	var body utils.ProductDetails
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var ProductDetails domain.ProductDetails
	copier.Copy(&ProductDetails, &body)
	if err := cr.productUseCase.AddProductDetails(c.Request.Context(), ProductDetails); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to add product details", err.Error(), body)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully added the product details", ProductDetails.ProductDetails, ProductDetails.InStock)
	c.JSON(http.StatusOK, response)
}

// LIST product details
func (cr *ProductHandler) ListProductDetailsById(c *gin.Context) {
	id := c.Param("productid")

	productDetails, err := cr.productUseCase.ListProductDetailsById(c.Request.Context(), id)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list product details", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product details", productDetails)
	c.JSON(http.StatusOK, response)
}

// List Product and product details
func (cr *ProductHandler) ListProductAndDetailsById(c *gin.Context) {
	id := c.Param("productid")

	productAndDetails, err := cr.productUseCase.ListProductAndDetailsById(c.Request.Context(), id)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list product details", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product details", productAndDetails)
	c.JSON(http.StatusOK, response)
}

// EDIT product details
// func (cr *ProductHandler) UpdateProductDetailsbyId(c *gin.Context) {
// 	id:= c.Param("productdetailsid")
// }
