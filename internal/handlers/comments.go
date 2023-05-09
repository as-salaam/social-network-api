package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/social-network-api/internal/models"
	"log"
	"net/http"
)

type createCommentData struct {
	Text string `json:"text" binding:"required"`
}

func (h *Handler) CreateComment(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		log.Println("claims doesn't exist")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims := claimsData.(*models.Claims)

	var commentData createCommentData
	if err := c.ShouldBindJSON(&commentData); err != nil {
		log.Println("invalid comment data")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	var post models.Post
	err := h.DB.Where("id = ?", c.Param("postID")).First(&post).Error
	if err != nil {
		log.Println("getting post:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Post not found",
		})
		return
	}

	var comment models.Comment

	comment.UserID = claims.UserID
	comment.PostID = post.ID
	comment.Text = commentData.Text

	if err := h.DB.Create(&comment).Error; err != nil {
		log.Println("inserting post data to DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}
