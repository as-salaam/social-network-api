package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"log"
	"net/http"
	"os"
	"strings"
)

type PostUpdateData struct {
	Content string `json:"content" binding:"required"`
}

func (h *Handler) UpdatePost(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims := claimsData.(*models.Claims)

	var post models.Post
	err := h.DB.Where("id = ?", c.Param("postID")).First(&post).Error
	if err != nil {
		log.Println("getting post:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	if post.UserID != claims.UserID {
		log.Println("updating another users post")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var postData PostUpdateData

	if err := c.ShouldBindJSON(&postData); err != nil {
		log.Println("binding post data:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
			"errors":  err.Error(),
		})
		return
	}

	if err := h.DB.Model(&post).Updates(postData).Error; err != nil {
		log.Println("updating post data in DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) DeletePost(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	var post models.Post
	if err := h.DB.Where("id = ?", c.Param("postID")).Preload("Files").First(&post).Error; err != nil {
		log.Println("getting post from db", err)
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

	for _, file := range post.Files {
		pathSlice := strings.Split(file.Path, "/")
		fileName := pathSlice[len(pathSlice)-1]
		err := os.RemoveAll("assets/post-files/" + fileName)
		if err != nil {
			log.Println("removing file from filesystem:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}
	}

	if err := h.DB.Model(&post).Association("Files").Clear(); err != nil {
		log.Println("deleting files from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "InternalServerError",
		})
		return
	}

	if err := h.DB.Delete(&post).Error; err != nil {
		log.Println("deleting post from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "InternalServerError",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted post",
	})
}
