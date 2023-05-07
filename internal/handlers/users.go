package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"log"
	"net/http"
)

func (h *Handler) GetOneUser(c *gin.Context) {
	var user models.User
	if err := h.DB.Where("login = ?", c.Param("login")).Preload("Profile").Preload("Posts.Files").First(&user).Error; err != nil {
		log.Println("getting a user:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "no such user",
		})
		return
	}

	if user.Profile.Type != "public" {
		c.JSON(http.StatusOK, gin.H{
			"message": "this profile is private",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
