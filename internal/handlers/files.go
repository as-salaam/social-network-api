package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type profileDataForUpdate struct {
	Email string `json:"omitempty,email"`
	Type  string `json:"type" binding:"required"`
	Bio   string `json:"bio"`
	Link  string `json:"omitempty,url"`
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

	var profile models.Profile

	if err := h.DB.Where("id = ?", c.Param("profileID")).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a profile:", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Models not found",
			})
			return
		}
		log.Println("getting a profile:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if profile.UserID != claims.UserID {
		log.Println("updating another users profile")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var profileData profileDataForUpdate

	if err := c.ShouldBindJSON(&profileData); err != nil {
		log.Println("binding profile data:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if err := h.DB.Model(&profile).Updates(profileData).Error; err != nil {
		log.Println("updating profile data in DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}
