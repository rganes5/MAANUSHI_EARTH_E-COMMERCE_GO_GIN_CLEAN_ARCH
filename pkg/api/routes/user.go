package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/go-gin-clean-arch/pkg/api/handler"
	"github.com/rganes5/go-gin-clean-arch/pkg/api/middleware"
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
		login.POST("/login", userHandler.UserLoginHandler)

	}
	home := api.Group("/user")
	{
		//AuthorizationMiddleware as middleware to perform authorization checks for users accessing the "/user" endpoint.
		home.Use(middleware.AuthorizationMiddleware("user"))
		home.GET("/home", userHandler.UserHomehandler)
		home.POST("/logout", userHandler.UserLogoutHandler)
	}
}
