package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"log"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("X-Auth-Token")
		if tokenString == "" {
			log.Print()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "User unauthorized",
			})
		}

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return "key", nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("authClaims", claims)
		c.Next()
	}
}
