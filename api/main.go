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

	org := "Gofers"
	bucket := "Clip-URL"
	influx_client := database.CreateInfluxClient()
	defer influx_client.Close()
	write_api := influx_client.WriteAPIBlocking(org, bucket)

	err = clipboard.Init()
	if err != nil {
		log.Fatalln("Internal server error: Unable to initialize clipboard watcher")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go controllers.SetupWatcher(ctx)

	r := routes.SetupRouter(db, write_api)
	r.Run(":" + os.Getenv("PORT"))
}
