package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/auth"
)

// The middleware verifies the presence and validity of a token stored in a cookie and sets the user's email in the Gin context if the authorization is successful.
func AuthorizationMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(role + "-token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Needs to login",
			})
			return
		}
		claims, err1 := auth.ValidateToken(tokenString)
		if err1 != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err1.Error(),
			})
			return
		}
		c.Set(role+"-email", claims.Email)
		c.Set(role+"-id", claims.ID)
		c.Next()
	}
}
