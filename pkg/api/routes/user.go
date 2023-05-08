package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/go-gin-clean-arch/pkg/api/handler"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler) {
	// signup := api.Group("/user")
	{
		// signup.POST("/signup", userHandler.SignUp)
		// signup.POST("/signup/otp/verify", userHandler.SignupOtpverify)

	}
}
