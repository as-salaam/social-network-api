package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/softclub-go-0-0/social-network-api/internal/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("X-Auth-Token")
		if tokenString == "" {
			log.Print("empty token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		claims := &models.Claims{}
		jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return "key", nil
		})
		if err != nil || !jwtToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		var token models.Token
		if err = db.Where("id = ?", claims.TokenID).First(&token).Error; err != nil {
			log.Print("getting token:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		if !token.IsValid {
			log.Print("token is invalid:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		c.Set("authClaims", claims)
		c.Next()
	}
}
