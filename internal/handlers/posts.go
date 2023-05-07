package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"log"
	"net/http"
	"path/filepath"
)

type PostData struct {
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

	var postData PostData

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

	if err := h.DB.Where("id = ?", c.Param("postID")).First(&post).Error; err != nil {
		log.Println(err.Error())
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

	if err := h.DB.Where("id = ?", c.Param("postID")).First(&post).Error; err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "InternalServerError",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted post",
	})
}

func (h *Handler) CreatePost(c *gin.Context) {
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
		if extension != ".jpg" || extension != ".jpeg" || extension != ".png" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid file extension",
			})
			return
		}
	}

	var post models.Post

	post.UserID = claims.UserID
	post.Content = content[0]

	if err := h.DB.Create(&post).Error; err != nil {
		log.Println("inserting post data to DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	for _, file := range files {
		extension := filepath.Ext(file.Filename)
		newFileName := "assets/post-files/"
		newFileName += uuid.New().String() + extension

		err := c.SaveUploadedFile(file, newFileName)
		if err != nil {
			log.Println("saving file to filesystem:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}

		var postFile models.File
		postFile.PostID = post.ID
		postFile.Path = "http://127.0.0.1:4000/" + newFileName

		if err = h.DB.Create(&postFile).Error; err != nil {
			log.Println("inserting file data to DB:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}
	}

	if err := h.DB.Preload("Files").First(&post).Error; err != nil {
		log.Println("getting post data from DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, post)
}
