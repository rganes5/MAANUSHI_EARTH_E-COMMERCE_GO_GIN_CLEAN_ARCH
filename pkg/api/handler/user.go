package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domain "github.com/rganes5/go-gin-clean-arch/pkg/domain"
	"github.com/rganes5/go-gin-clean-arch/pkg/support"
	services "github.com/rganes5/go-gin-clean-arch/pkg/usecase/interface"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

// type Response struct {
// 	ID   uint   `copier:"must"`
// 	Name string `copier:"must"`
// }

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// Variable declared contataining type as users which is already initialiazed in domanin folder.
var signUp_user domain.Users

// USER SIGN-UP
// BINDING
func (cr *UserHandler) UserSignUp(c *gin.Context) {
	if err := c.BindJSON(&signUp_user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if ok := support.Email_validator(signUp_user.Email); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Email Format is incorrect",
		})
		return
	}

	// if ok := support.MobileNum_validator(signUp_user.PhoneNum); !ok {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"error": "Enter a valid number",
	// 	})
	// 	return
	// }

	if err := support.MobileNum_validator(signUp_user.PhoneNum); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _, err := cr.userUseCase.FindByEmail(c.Request.Context(), signUp_user.Email); err == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User with the email already exits!",
		})
		return
	}

	signUp_user.Password, _ = support.HashPassword(signUp_user.Password)
	err := cr.userUseCase.SignUpUser(c.Request.Context(), signUp_user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"User registration": "Success",
	})

}

// FindAll godoc
// @summary Get all users
// @description Get all users
// @tags users
// @security ApiKeyAuth
// @id FindAll
// @produce json
// @Router /api/users [get]
// @response 200 {object} []Response "OK"
// func (cr *UserHandler) FindAll(c *gin.Context) {
// 	users, err := cr.userUseCase.FindAll(c.Request.Context())

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusNotFound)
// 	} else {
// 		response := []Response{}
// 		copier.Copy(&response, &users)

// 		c.JSON(http.StatusOK, response)
// 	}
// }

// func (cr *UserHandler) FindByID(c *gin.Context) {
// 	paramsId := c.Param("id")
// 	id, err := strconv.Atoi(paramsId)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "cannot parse id",
// 		})
// 		return
// 	}

// 	user, err := cr.userUseCase.FindByID(c.Request.Context(), uint(id))

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusNotFound)
// 	} else {
// 		response := Response{}
// 		copier.Copy(&response, &user)

// 		c.JSON(http.StatusOK, response)
// 	}
// }

// func (cr *UserHandler) Save(c *gin.Context) {
// 	var user domain.Users

// 	if err := c.BindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 		return
// 	}

// 	user, err := cr.userUseCase.Save(c.Request.Context(), user)

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusNotFound)
// 	} else {
// 		response := Response{}
// 		copier.Copy(&response, &user)

// 		c.JSON(http.StatusOK, response)
// 	}
// }

// func (cr *UserHandler) Delete(c *gin.Context) {
// 	paramsId := c.Param("id")
// 	id, err := strconv.Atoi(paramsId)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Cannot parse id",
// 		})
// 		return
// 	}

// 	ctx := c.Request.Context()
// 	user, err := cr.userUseCase.FindByID(ctx, uint(id))

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusNotFound)
// 	}

// 	if user == (domain.Users{}) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "User is not booking yet",
// 		})
// 		return
// 	}

// 	cr.userUseCase.Delete(ctx, user)

// 	c.JSON(http.StatusOK, gin.H{"message": "User is deleted successfully"})
// }
