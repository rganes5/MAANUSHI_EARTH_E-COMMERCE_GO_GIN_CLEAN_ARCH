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

// ADD CATEGORY
// @Summary API FOR ADDING CATEGORY
// @ID ADMIN-ADD-CATEGORY
// @Description ADDING CATEGORY FROM ADMINS END
// @Tags PRODUCT CATEGORY
// @Accept json
// @Produce json
// @Param category_details body utils.AddCategory true "Enter the category name"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/category/add [post]
func (cr *ProductHandler) AddCategory(c *gin.Context) {
	var body utils.AddCategory
	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: invalid input", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var category domain.Category
	copier.Copy(&category, &body)
	err := cr.productUseCase.AddCategory(c.Request.Context(), category)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Faild to add cateogy", err.Error(), body)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: successfully added a new category", body)
	c.JSON(http.StatusOK, response)
}

// EDIT CATEGORY
// @Summary API FOR EDITING CATEGORY
// @ID ADMIN-EDIT-CATEGORY
// @Description UPDATING CATEGORY NAME FROM ADMINS END
// @Tags PRODUCT CATEGORY
// @Accept json
// @Produce json
// @Param category_id path string true "Enter the category id that you would like to make the change"
// @Param category_details body utils.UpdateCategory true "Enter the category details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/category/update/{category_id} [patch]
func (cr *ProductHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("category_id")
	var categoryname utils.UpdateCategory
	if err := c.BindJSON(&categoryname); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind JSON", err.Error(), categoryname)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var categories domain.Category
	copier.Copy(&categories, &categoryname)
	if err := cr.productUseCase.UpdateCategory(c.Request.Context(), categories, id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update category", err.Error(), categoryname)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success:Successfully updated", categoryname)
	c.JSON(http.StatusOK, response)
}

// DELETE CATEGORY
// @Summary API FOR DELETING A CATEGORY
// @ID ADMIN-DELETE-CATEGORY
// @Description DELETING CATEGORY AND ALSO CHECKING WHETHER IT HAS A EXISTING PRODUCT
// @Tags PRODUCT CATEGORY
// @Accept json
// @Produce json
// @Param category_id path string true "Enter the category id that you would like to delete"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/category/delete/{category_id} [post]
func (cr *ProductHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("category_id")
	if err := cr.productUseCase.DeleteCategory(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Error: Faild to Delete Cateogy with id", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully deleted category with id", id)
	c.JSON(http.StatusOK, response)
}

// LIST CATEGORY
// @Summary API FOR LISTING ALL CATEGORIES
// @Description LISTING ALL CATEGORIES FROM ADMINS AND USERS END
// @Tags PRODUCT CATEGORY
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/category/listall [get]
// @Router /admin/category/listall [get]
func (cr *ProductHandler) ListCategories(c *gin.Context) {
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
	categories, err := cr.productUseCase.ListCategories(c.Request.Context(), pagination)
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

// ADD PRODUCT
// @Summary API FOR ADDING PRODUCT
// @ID ADMIN-ADD-PRODUCT
// @Description ADDING PRODUCT FROM ADMINS END
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param product_details body utils.Products true "Enter the product details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/products/add [post]
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

// DELETE PRODUCT
// @Summary API FOR DELETING A PRODUCT
// @ID ADMIN-DELETE-PRODUCT
// @Description DELETING PRODUCT BASED ON PRODUCT ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param product_id path string true "Enter the product id that you would like to delete"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/products/delete/{product_id} [post]
func (cr *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("product_id")
	if err := cr.productUseCase.DeleteProduct(c.Request.Context(), id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to delete the product", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product successfully deleted with product id", id)
	c.JSON(http.StatusOK, response)
}

// EDIT PRODUCT
// @Summary API FOR EDITING PRODUCT
// @ID ADMIN-EDIT-PRODUCT
// @Description UPDATING PRODUCT DETAILS FROM ADMINS END
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param product_id path string true "Enter the product id that you would like to make the change"
// @Param product_details body utils.Products true "Enter the category details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/products/update/{product_id} [patch]
func (cr *ProductHandler) EditProduct(c *gin.Context) {
	// var product domain.Products
	var body utils.Products
	id := c.Param("product_id")
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
// @Summary API FOR LISTING ALL PRODUCTS
// @Description LISTING ALL PRODUCTS FROM ADMINS AND USERS END
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/products/listall [get]
// @Router /admin/products/listall [get]
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

// LIST PRODUCTS BASED ON CATEGORY
// @Summary API FOR LISTING ALL PRODUCTS BASED ON CATEGORY
// @Description LISTING ALL PRODUCTS FROM ADMINS AND USERS END BASED ON CATEGORY
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param category_id path string true "Enter the category id"
// @Param page query int false "Enter the page number to display"
// @Param limit query int false "Number of items to retrieve per page"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/products/listbasedoncategory/{category_id} [get]
// @Router /admin/products/listbasedoncategory/{category_id} [get]
func (cr *ProductHandler) ListProductsBasedOnCategory(c *gin.Context) {
	id := c.Param("category_id")
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
	products, err := cr.productUseCase.ListProductsBasedOnCategory(c.Request.Context(), id, pagination)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: failed to retreive products based on the given category", err.Error(), id)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: ", products)
	c.JSON(http.StatusOK, response)
}

// ADD PRODUCT DETAILS
// @Summary API FOR ADDING PRODUCT DETAILS
// @ID ADMIN-ADD-PRODUCT-DETAILS
// @Description ADDING PRODUCT DETAILS FROM ADMINS END
// @Tags PRODUCT DETAILS
// @Accept json
// @Produce json
// @Param product_details body utils.ProductDetails true "Enter the product details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/productsDetails/add [post]
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

// LIST PRODUCTS DETAILS
// @Summary API FOR LISTING PRODUCTS DETAILS BY ID
// @Description LISTING ALL PRODUCTS DETAILS FROM ADMINS AND USERS END
// @Tags PRODUCT DETAILS
// @Accept json
// @Produce json
// @Param product_id path string true "Enter the product id that you would like to see the details of"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/products/findproductdetails/{product_id} [get]
// @Router /admin/productsDetails/findproductdetails/{product_id} [get]
func (cr *ProductHandler) ListProductDetailsById(c *gin.Context) {
	id := c.Param("product_id")

	productDetails, err := cr.productUseCase.ListProductDetailsById(c.Request.Context(), id)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list product details", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product details", productDetails)
	c.JSON(http.StatusOK, response)
}

// LIST PRODUCT AND PRODUCT_DETAILS
// @Summary API FOR LISTING PRODUCT AND PRODUCT_DETAILS DETAILS BY ID
// @Description LISTING ALL PRODUCT AND PRODUCT_DETAILS FROM ADMINS AND USERS END
// @Tags PRODUCT DETAILS
// @Accept json
// @Produce json
// @Param product_id path string true "Enter the product id that you would like to see the details of"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/products/findproductanddetails/{product_id} [get]
// @Router /admin/productsDetails/findproductanddetails/{product_id} [get]
func (cr *ProductHandler) ListProductAndDetailsById(c *gin.Context) {
	id := c.Param("product_id")

	productAndDetails, err := cr.productUseCase.ListProductAndDetailsById(c.Request.Context(), id)
	if err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to list product details", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Product details", productAndDetails)
	c.JSON(http.StatusOK, response)
}

// EDIT PRODUCT DETAILS BY ID
// @Summary API FOR EDITING PRODUCT DETAILS
// @ID ADMIN-EDIT-PRODUCT-DETAILS
// @Description UPDATING PRODUCT DETAILS FROM ADMINS END
// @Tags PRODUCT DETAILS
// @Accept json
// @Produce json
// @Param productDetails_id path string true "Enter the product details id that you would like to make the change"
// @Param product_details body utils.ProductDetails true "Enter the product details"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/productsDetails/update/{product_id} [patch]
func (cr *ProductHandler) UpdateProductDetailsbyId(c *gin.Context) {
	id := c.Param("productDetails_id")
	var body utils.ProductDetails

	if err := c.BindJSON(&body); err != nil {
		response := utils.ErrorResponse(400, "Error: Failed to bind json body", err.Error(), body)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var product_details domain.ProductDetails
	copier.Copy(&product_details, &body)
	if err := cr.productUseCase.EditProductDetailsById(c.Request.Context(), product_details, id); err != nil {
		response := utils.ErrorResponse(500, "Error: Failed to update the details", err.Error(), body)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.SuccessResponse(200, "Success: Successfully edited the product details", body)
	c.JSON(http.StatusOK, response)
}
