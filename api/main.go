package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sankalp-12/clip-url/database"
	"github.com/sankalp-12/clip-url/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Internal server error: Unable to load the env file")
	}

	db, err := database.CreateClient()
	if err != nil {
		log.Fatalln("Internal server error: Unable to connect to BadgerDB")
	}

	r := routes.SetupRouter(db)
	r.Run(":" + os.Getenv("PORT"))
}
