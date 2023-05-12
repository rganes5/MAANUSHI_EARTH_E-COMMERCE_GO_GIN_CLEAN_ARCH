package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/handler"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	signup := api.Group("/admin")
	{
		signup.POST("/signup", adminHandler.AdminSignUp)
	}
	login := api.Group("/admin")
	{
		login.POST("/login", adminHandler.AdminLogin)
	}
	home := api.Group("/admin")
	{
		home.Use(middleware.AuthorizationMiddleware("admin"))
		home.GET("/home", adminHandler.HomeHandler)
		home.POST("/logout", adminHandler.Logout)

	}
}
