package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rganes5/go-gin-clean-arch/pkg/api/handler"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	signup := api.Group("/admin")
	{
		signup.POST("/signup", adminHandler.AdminSignUp)
	}
}
