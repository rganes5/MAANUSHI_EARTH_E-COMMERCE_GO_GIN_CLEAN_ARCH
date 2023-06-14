package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/auth"
)

// The middleware verifies the presence and validity of a token stored in a cookie and sets the user's email in the Gin context if the authorization is successful.
//
//	func AuthorizationMiddleware(role string) gin.HandlerFunc {
//		return func(c *gin.Context) {
//			tokenString, err := c.Cookie(role + "-token")
//			if err != nil || tokenString == "" {
//				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//					"error": "Needs to login",
//				})
//				return
//			}
//			claims, err1 := auth.ValidateToken(tokenString)
//			if err1 != nil {
//				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//					"error": err1.Error(),
//				})
//				return
//			}
//			c.Set(role+"-email", claims.Email)
//			c.Set(role+"-id", claims.ID)
//			c.Next()
//		}
//	}
func AuthorizationMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(role + "-token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Please login first",
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
		// Check if the token is about to expire or has already expired
		if time.Now().After(claims.ExpiresAt.Time) {
			// Token has expired, remove the cookie
			c.SetCookie(role+"-token", "", -1, "/", "", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired, please login again",
			})
			return
		}
		// Refresh the token expiry time on each API call
		expiryTime := time.Now().Add(10 * time.Minute)
		claims.ExpiresAt = jwt.NewNumericDate(expiryTime)
		tokenString, err = auth.GenerateJWT(claims.Email, claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}
		c.SetCookie(role+"-token", tokenString, int(10*time.Minute.Seconds()), "/", "", false, true)
		c.Set(role+"-email", claims.Email)
		c.Set(role+"-id", claims.ID)
		c.Next()
	}
}
