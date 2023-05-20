package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/handler"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler) {
	//sets up a route group for the "/user" endpoint
	signup := api.Group("/user")
	{
		signup.POST("/signup", userHandler.UserSignUp)
		signup.POST("/signup/otp/verify", userHandler.SignupOtpverify)
	}
	login := api.Group("/user")
	{
		login.POST("/login", userHandler.LoginHandler)

	}
	home := login.Group("/")
	{
		//AuthorizationMiddleware as middleware to perform authorization checks for users accessing the "/user" endpoint.
		home.Use(middleware.AuthorizationMiddleware("user"))
		home.GET("/home", userHandler.Homehandler)
		home.POST("/logout", userHandler.LogoutHandler)
		product := home.Group("/products")
		{
			product.GET("/listall", userHandler.ListProducts)
		}
		userprofile := home.Group("/profile")
		{
			userprofile.POST("/add/address", userHandler.AddAddress)
		}
	}

}
