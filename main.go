package main

import (
	"log"

	"portalnesia.com/api/models"
	"portalnesia.com/api/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	models.SetupDB()

	r := routes.SetupRouters()

	r.Listen(":3000")
}
