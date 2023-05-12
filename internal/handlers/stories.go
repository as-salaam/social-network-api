package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/social-network-api/internal/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) CreateStory(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims := claimsData.(*models.Claims)

	form, _ := c.MultipartForm()

	file := form.File["file"][0]
	extension := filepath.Ext(file.Filename)
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".mp4" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid file extension",
		})
		return
	}

	newFileName := "assets/story-files/"
	newFileName += uuid.New().String() + extension

	err := c.SaveUploadedFile(file, newFileName)
	if err != nil {
		log.Println("saving file to filesystem:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	var story models.Story
	story.UserID = claims.UserID
	story.FilePath = "/" + newFileName

	if err = h.DB.Create(&story).Error; err != nil {
		log.Println("inserting story data to DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, story)
}

func (h *Handler) GetStory(c *gin.Context) {
	var story models.Story
	if err := h.DB.Where("id = ?", c.Param("storyID")).Preload("Files").First(&story).Error; err != nil {
		log.Println("getting story from DB:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	var user models.User
	if err := h.DB.Where("id = ?", story.UserID).First(&user).Error; err != nil {
		log.Println("getting user from DB:", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	claimsData, exist := c.Get("authClaims")
	if !exist {
		log.Println("claims doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	if story.UserID == claims.UserID {
		c.JSON(http.StatusOK, story)
		return
	} else {
		var profile models.Profile
		if result := h.DB.Where("user_id = ?", user.ID).First(&profile); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Not Found",
			})
			return
		}

		if profile.Type == "public" {
			c.JSON(http.StatusOK, story)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "This account is private",
			})
			return
		}
	}
}

type StoryData struct {
	Content string `json:"content" binding:"required"`
}

func (h *Handler) DeleteStory(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		log.Println("claims doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	var story models.Post
	if err := h.DB.Where("id = ?", c.Param("postID")).Preload("Files").First(&story).Error; err != nil {
		log.Println("getting story from db", err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	if story.UserID != claims.UserID {
		log.Println("deleting another user's story")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	for _, file := range story.Files {
		pathSlice := strings.Split(file.Path, "/")
		fileName := pathSlice[len(pathSlice)-1]
		err := os.RemoveAll("assets/story-files/" + fileName)
		if err != nil {
			log.Println("removing file from filesystem:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}
	}

	if err := h.DB.Delete(&story).Error; err != nil {
		log.Println("deleting post from DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted story",
	})
}
