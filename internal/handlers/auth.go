package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/softclub-go-0-0/instagram-api-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
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

	var newProfile models.Profile
	newProfile.UserID = user.ID

	if result := h.DB.Create(&newProfile); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, user)

}

func (h *Handler) Login(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var user models.User
	if result := h.DB.Where("login = ?", credentials.Login).First(&user); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var token models.Token
	token.ExpiresAt = time.Now().Add(10 * time.Minute)

	if result := h.DB.Create(&token); result.Error != nil {
		log.Println("inserting token data to DB:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	claims := &models.Claims{
		UserID:  user.ID,
		TokenID: token.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(token.ExpiresAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	claimsData, exist := c.Get("authClaims")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "",
		})
		return
	}

	claims := claimsData.(*models.Claims)

	var token models.Token

	if err := h.DB.Where("id = ?", c.Param("TokenID")).First(&token).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
	}

	token.IsValid = false

	if result := h.DB.Save(&token); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}
