package main

import (
	"log"

	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/db"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

//	@title			Dusky API
//	@description	This is a sample server for a social media application called Dusky. It provides endpoints for managing users, posts, comments, and feeds.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Provide your API key in the Authorization header as follows: "Bearer {your_api_key}"
func main() {

	err := godotenv.Load()
	logger := utils.NewLogger()
	defer logger.Sync()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	config, err := config.InitializeConfig()
	if err != nil {
		logger.Fatal("Error initializing config:", err)
	}

	db, err := db.New(config.Db.Addr, config.Db.MaxOpenConns, config.Db.MaxIdleConns, config.Db.MaxIdleTime)
	if err != nil {
		logger.Panic("Error connecting to the database:", err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	appConfig := AppConfig{
		addr:   config.Server.Host,
		apiUrl: config.ApiUrl,
	}

	mailConfig := config.Mail

	app := NewApplication(appConfig, store, db, logger, mailConfig)

	mux := app.mount()

	log.Fatal(app.Run(mux))
}
