package main

import (
	"log"

	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.InitializeConfig()
	if err != nil {
		log.Fatal("Error initializing config:", err)
	}

	app := NewApplication(AppConfig{
		addr: config.Server.Host,
	})

	mux := app.mount()

	log.Fatal(app.Run(mux))
}
