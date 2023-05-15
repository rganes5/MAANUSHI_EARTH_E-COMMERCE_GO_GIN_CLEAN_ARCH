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
			users.PATCH("/:userid/make", adminHandler.AccessHandler)
		}
		category := home.Group("/category")
		{
			category.POST("/add", adminHandler.AddCategory)
			category.POST("/delete/:categoryid", adminHandler.DeleteCategory)
			category.GET("/listall", adminHandler.ListCategories)
		}
		products := home.Group("/products")
		{
			products.POST("/add", productHandler.AddProduct)
			products.POST("/delete", productHandler.DeleteProduct)
		}

	}
}
