package main

import (
	"log"
	"os"

	"github.com/sankalp-12/clip-url/database"
	"github.com/sankalp-12/clip-url/routes"
)

func main() {
	db, err := database.CreateClient()
	if err != nil {
		log.Fatalln("Internal server error: Unable to connect to the DB")
	}

	r := routes.SetupRouter(db)
	r.Run(":" + os.Getenv("PORT"))
}
