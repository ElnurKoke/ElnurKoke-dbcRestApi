package main

import (
	"dbc/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	router := gin.Default()

	routes.InitRoutes(router)

	log.Fatal(router.Run(":8080"))
}
