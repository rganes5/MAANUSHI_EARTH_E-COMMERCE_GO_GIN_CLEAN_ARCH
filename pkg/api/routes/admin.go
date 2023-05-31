package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/handler"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler, productHandler *handler.ProductHandler) {
	signup := api.Group("/admin")
	{
		signup.POST("/signup", adminHandler.AdminSignUp)
	}
	login := api.Group("/admin")
	{
		login.POST("/login", adminHandler.AdminLogin)
	}
	home := login.Group("/")
	{
		home.Use(middleware.AuthorizationMiddleware("admin"))
		home.GET("/home", adminHandler.HomeHandler)
		home.POST("/logout", adminHandler.Logout)
		users := home.Group("/user")
		{
			users.GET("/", adminHandler.ListUsers)
			users.PATCH("/:user_id/make", adminHandler.AccessHandler)
		}
		category := home.Group("/category")
		{
			category.POST("/add", productHandler.AddCategory)
			category.PATCH("/update/:category_id", productHandler.UpdateCategory)
			category.POST("/delete/:category_id", productHandler.DeleteCategory)
			category.GET("/listall", productHandler.ListCategories)
		}
		products := home.Group("/products")
		{
			products.POST("/add", productHandler.AddProduct)
			products.POST("/delete/:product_id", productHandler.DeleteProduct)
			products.PATCH("/update/:product_id", productHandler.EditProduct)
			products.GET("/listall", productHandler.ListProducts)
			products.GET("/listbasedoncategory/:category_id", productHandler.ListProductsBasedOnCategory)

		}
		productDetails := home.Group("/productsDetails")
		{
			productDetails.POST("/add", productHandler.AddProductDetails)
			productDetails.PATCH("/update/:productDetails_id", productHandler.UpdateProductDetailsbyId)
			productDetails.GET("/findproductdetails/:product_id", productHandler.ListProductDetailsById)
			productDetails.GET("/findproductanddetails/:product_id", productHandler.ListProductAndDetailsById)

		}

	}
}
