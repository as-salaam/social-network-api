package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type profileDataForUpdate struct {
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	Bio    string `json:"bio"`
	Link   string `json:"link"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var profile models.Profile

	if err := h.DB.Where("id = ?", c.Param("profileID")).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("getting a profile:", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		}
		log.Println("getting a profile:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var profileData profileDataForUpdate

	if err := c.ShouldBindJSON(&profileData); err != nil {
		log.Println("binding profile data:", err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	if err := h.DB.Model(&profile).Updates(profileData).Error; err != nil {
		log.Println("updating profile data in DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, profile)
}
