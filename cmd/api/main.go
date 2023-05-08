package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/softclub-go-0-0/social-network-api/internal/database"
	"github.com/softclub-go-0-0/social-network-api/internal/handlers"
	"github.com/softclub-go-0-0/social-network-api/internal/middlewares"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file:", err)
	}

	DBHost := flag.String("dbhost", "localhost", "Enter the host of the DB server")
	DBName := flag.String("dbname", "social_network_api", "Enter the name of the DB")
	DBUser := flag.String("dbuser", "postgres", "Enter the name of a DB user")
	DBPassword := flag.String("dbpassword", "postgres", "Enter the password of user")
	DBPort := flag.Uint("dbport", 5432, "Enter the port of DB")
	Timezone := flag.String("dbtimezone", "Asia/Dushanbe", "Enter your timezone to connect to the DB")
	DBSSLMode := flag.Bool("dbsslmode", false, "Turns on ssl mode while connecting to DB")
	Port := flag.Uint("listenport", 4000, "Which port to listen")
	flag.Parse()

	db, err := database.DBInit(*DBHost, *DBName, *DBUser, *DBPassword, *DBPort, *Timezone, *DBSSLMode)
	if err != nil {
		log.Fatal("db connection:", err)
	}

	h := handlers.NewHandler(db)

	router := gin.Default()

	// files route
	router.Static("/assets", "./assets")

	// users routes
	router.GET("/:login", h.GetOneUser) // +

	// auth routes
	router.POST("/register", h.Register) // +
	router.POST("/login", h.Login)       // -

	router.Use(middlewares.AuthMiddleware(db)) // -

	router.POST("/logout", h.Logout) // +

	// profile routes
	router.PUT("/profile", h.UpdateProfile)          // -
	router.POST("/profile/avatar", h.UploadAvatar)   // -
	router.DELETE("/profile/avatar", h.RemoveAvatar) // -

	// posts routes
	router.POST("/posts", h.CreatePost)           // -
	router.GET("/posts/:postID", h.GetPost)       // -
	router.PUT("/posts/:postID", h.UpdatePost)    // +
	router.DELETE("/posts/:postID", h.DeletePost) // -

	log.Fatal("router running:", router.Run(fmt.Sprintf(":%d", *Port)))
}
