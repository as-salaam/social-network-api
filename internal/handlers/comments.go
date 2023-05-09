package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/social-network-api/internal/models"
	"log"
	"net/http"
	"path/filepath"
)

func (h *Handler) CreateComment(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims := claimsData.(*models.Claims)

	form, _ := c.MultipartForm()

	content, exists := form.Value["content"]
	if !exists || len(content) != 1 || content[0] == "" {
		log.Println("invalid post data")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation error",
			"errors":  "Content field is required and it should be text",
		})
		return
	}

	files := form.File["photos[]"]
	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid file extension",
			})
			return
		}
	}

	var comment models.Comment

	comment.UserID = claims.UserID
	comment.Content = content[0]

	if err := h.DB.Create(&comment).Error; err != nil {
		log.Println("inserting post data to DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
}
