package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/social-network-api/internal/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type profileDataForUpdate struct {
	Email string `json:"email" binding:"omitempty,email"`
	Type  string `json:"type" binding:"required"`
	Bio   string `json:"bio"`
	Link  string `json:"link" binding:"omitempty,url"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	var user models.User

	if err := h.DB.Where("id = ?", claims.UserID).Preload("Profile").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a user's profile:", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}
		log.Println("getting a user's profile:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	var profileData profileDataForUpdate

	if err := c.ShouldBindJSON(&profileData); err != nil {
		log.Println("binding profile data:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	user.Profile.Email = profileData.Email
	user.Profile.Type = profileData.Type
	user.Profile.Bio = profileData.Bio
	user.Profile.Link = profileData.Link

	// todo: find out why using .Updates(profileData) doesn't put updated_at timestamp
	if err := h.DB.Save(&user.Profile).Error; err != nil {
		log.Println("updating profile data in DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, user.Profile)
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		log.Println("claims doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	var user models.User
	if err := h.DB.Where("id = ?", claims.UserID).Preload("Profile").First(&user).Error; err != nil {
		log.Println("getting user's profile:", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	uploadedPhoto, err := c.FormFile("photo")
	if err != nil {
		log.Println("getting file from request:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
		})
		return
	}

	extension := filepath.Ext(uploadedPhoto.Filename)

	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
		log.Println("invalid file extension:", extension)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid file extension. Upload an image please.",
		})
		return
	}

	newFileName := "assets/avatars/"
	newFileName += uuid.New().String() + extension

	err = c.SaveUploadedFile(uploadedPhoto, newFileName)
	if err != nil {
		log.Println("saving photo to filesystem", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	user.Profile.Avatar = "/" + newFileName

	if err = h.DB.Save(&user.Profile).Error; err != nil {
		log.Println("updating profile", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, user.Profile)
}

func (h *Handler) RemoveAvatar(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		log.Println("claims doesn't exist")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	claims := claimsData.(*models.Claims)

	var user models.User
	if err := h.DB.Where("id = ?", claims.UserID).Preload("Profile").First(&user).Error; err != nil {
		log.Println("getting user's profile:", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	pathSlice := strings.Split(user.Profile.Avatar, "/")
	fileName := pathSlice[len(pathSlice)-1]
	err := os.RemoveAll("assets/avatars/" + fileName)
	if err != nil {
		log.Println("removing avatar from filesystem", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	user.Profile.Avatar = ""

	if err = h.DB.Save(&user.Profile).Error; err != nil {
		log.Println("updating profile", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, user.Profile)
}
