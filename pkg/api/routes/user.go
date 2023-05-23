package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/handler"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler) {
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
	forgotPassword := api.Group("/user/forgot/password")
	{
		forgotPassword.POST("/", userHandler.ForgotPassword)
		forgotPassword.PATCH("/otp/verify", userHandler.ForgotPasswordOtpVerify)
	}
	home := login.Group("/")
	{
		//AuthorizationMiddleware as middleware to perform authorization checks for users accessing the "/user" endpoint.
		home.Use(middleware.AuthorizationMiddleware("user"))
		home.POST("/logout", userHandler.LogoutHandler)
		product := home.Group("/products")
		{
			product.GET("/listall", productHandler.ListProducts)
			product.GET("/findproductdetails/:productid", productHandler.ListProductDetailsById)
			product.GET("/findproductanddetails/:productid", productHandler.ListProductAndDetailsById)
		}
		userprofile := home.Group("/profile")
		{
			userprofile.GET("/home", userHandler.HomeHandler)
			userprofile.PATCH("/edit/profile", userHandler.UpdateProfile)
			userprofile.POST("/add/address", userHandler.AddAddress)
			userprofile.GET("/list/address", userHandler.ListAddress)
			userprofile.PATCH("/edit/address/:addressid", userHandler.UpdateAddress)
			userprofile.POST("/delete/address/:addressid", userHandler.DeleteAddress)
		}
	}

}
