package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"log"
	"net/http"
)

func (h *Handler) DeletePost(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	postID := c.Param("postID")

	var post models.Post
	result := h.DB.First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "NotFound",
		})
		return
	}

	if post.UserID != claims.UserID {
		log.Println("deleting another user's post")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	result = h.DB.Delete(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "InternalServerError",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted post",
	})
}
