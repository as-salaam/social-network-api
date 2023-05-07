package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/social-network-api/internal/handlers"
	"log"
)

func main() {
	//DBHost := flag.String("dbhost", "localhost", "Enter the host of the DB server")
	//DBName := flag.String("dbname", "auth_service", "Enter the name of the DB")
	//DBUser := flag.String("dbuser", "postgres", "Enter the name of a DB user")
	//DBPassword := flag.String("dbpassword", "postgres", "Enter the password of user")
	//DBPort := flag.Uint("dbport", 5432, "Enter the port of DB")
	//Timezone := flag.String("dbtimezone", "Asia/Dushanbe", "Enter your timezone to connect to the DB")
	//DBSSLMode := flag.Bool("dbsslmode", false, "Turns on ssl mode while connecting to DB")
	Port := flag.Uint("listenport", 4000, "Which port to listen")
	flag.Parse()

	//db, err := database.DBInit(*DBHost, *DBName, *DBUser, *DBPassword, *DBPort, *Timezone, *DBSSLMode)
	//if err != nil {
	//	log.Fatal("db connection:", err)
	//}

	h := handlers.NewHandler(nil)

	router := gin.Default()

	router.POST("/register", h.Register)
	router.POST("/posts", h.CreatePost)

	log.Fatal("router running:", router.Run(fmt.Sprintf(":%d", *Port)))
}
