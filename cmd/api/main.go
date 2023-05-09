package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/softclub-go-0-0/social-network-api/internal/database"
	"github.com/softclub-go-0-0/social-network-api/internal/handlers"
	"github.com/softclub-go-0-0/social-network-api/internal/middlewares"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file:", err)
	}

	dbport, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalln("Error parsing DB_PORT:", err)
	}

	sslmode, err := strconv.ParseBool(os.Getenv("DB_SSL_MODE"))
	if err != nil {
		log.Fatalln("Error parsing DB_SSL_MODE:", err)
	}

	db, err := database.DBInit(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		uint(dbport),
		os.Getenv("TIMEZONE"),
		sslmode,
	)
	if err != nil {
		log.Fatal("db connection error:", err)
	}

	h := handlers.NewHandler(db)

	router := gin.Default()

	// files route
	router.Static("/assets", "./assets")

	// users routes
	router.GET("/:login", h.GetOneUser)

	// auth routes
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	router.Use(middlewares.AuthMiddleware(db))

	router.POST("/logout", h.Logout)

	// profile routes
	router.PUT("/profile", h.UpdateProfile)
	router.POST("/profile/avatar", h.UploadAvatar)
	router.DELETE("/profile/avatar", h.RemoveAvatar)

	// posts routes
	router.POST("/posts", h.CreatePost)
	router.GET("/posts/:postID", h.GetPost)
	router.PUT("/posts/:postID", h.UpdatePost)
	router.DELETE("/posts/:postID", h.DeletePost)
	router.POST("/posts/:postID/comments", h.CreateComment)

	log.Fatal("router running:", router.Run(fmt.Sprintf(":%d", os.Getenv("APP_PORT"))))
}
