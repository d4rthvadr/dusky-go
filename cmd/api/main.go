package main

import (
	"log"

	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/db"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.InitializeConfig()
	if err != nil {
		log.Fatal("Error initializing config:", err)
	}

	db, err := db.New(config.Db.Addr, config.Db.MaxOpenConns, config.Db.MaxIdleConns, config.Db.MaxIdleTime)
	if err != nil {
		log.Panic("Error connecting to the database:", err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	app := NewApplication(AppConfig{
		addr: config.Server.Host,
	}, store, db)

	mux := app.mount()

	log.Fatal(app.Run(mux))
}
