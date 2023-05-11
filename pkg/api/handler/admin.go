package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rganes5/go-gin-clean-arch/pkg/domain"
	services "github.com/rganes5/go-gin-clean-arch/pkg/usecase/interface"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

// type Response struct {
// 	ID   uint   `copier:"must"`
// 	Name string `copier:"must"`
// }

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

// Variable declared contataining type as Admin which is already initialiazed in domain folder.
var signUp_admin domain.Admin

func (cr *AdminHandler) AdminSignUp(c *gin.Context) {
	//Binding
	if err := c.BindJSON(&signUp_admin); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

}
