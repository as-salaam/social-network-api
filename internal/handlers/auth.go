package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserRegistrationData struct {
	Login        string `json:"login" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Confirmation string `json:"confirmation" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var userData UserRegistrationData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	var user models.User
	if result := h.DB.Where("login = ?", userData.Login).First(&user); result.Error == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "There is a user registered with such login already",
		})
		return
	}

	if userData.Password != userData.Confirmation {
		log.Println("passwords doesn't match")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Passwords doesn't match",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 12)
	if err != nil {
		log.Print("generating password hash:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	user.Login = userData.Login
	user.Password = string(hashedPassword)

	if result := h.DB.Create(&user); result.Error != nil {
		log.Println("inserting user data to DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
