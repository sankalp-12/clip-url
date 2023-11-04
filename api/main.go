package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sankalp-12/clip-url/controllers"
	"github.com/sankalp-12/clip-url/database"
	"github.com/sankalp-12/clip-url/routes"
	"golang.design/x/clipboard"
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

	err = clipboard.Init()
	if err != nil {
		log.Fatalln("Internal server error: Unable to initialize clipboard watcher")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go controllers.SetupWatcher(ctx)

	r := routes.SetupRouter(db)
	r.Run(":" + os.Getenv("PORT"))
}
