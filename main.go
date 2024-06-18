package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
}

func main() {
	r := gin.Default()

	// Trust the Nginx reverse proxy running on localhost
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Use(validateSecret)

	fmt.Println("Starting Deployinator server on port 4444")

	r.GET("/*projectName", handleProjectDeploy)
	r.Run(":4444")
}
