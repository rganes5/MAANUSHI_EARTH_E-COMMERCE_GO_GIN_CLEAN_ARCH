package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rganes5/go-gin-clean-arch/pkg/auth"
)

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
				"error": err1,
			})
			return
		}
		c.Set(role+"-email", claims.Email)
		c.Next()
	}
}
