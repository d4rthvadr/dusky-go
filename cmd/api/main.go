package main

import (
	"context"
	"expvar"
	"log"
	"runtime"

	"github.com/d4rthvadr/dusky-go/internal/auth"
	"github.com/d4rthvadr/dusky-go/internal/cache"
	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/db"
	"github.com/d4rthvadr/dusky-go/internal/mailer"
	ratelimiter "github.com/d4rthvadr/dusky-go/internal/ratelmiter"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils/logger"
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
	logger := logger.NewLogger()
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

	var cacheStorage cache.CacheStorage
	var rdb *cache.RedisClient
	if config.CacheConfig.Enabled {
		rdb = cache.NewRedisClient(&cache.RedisOptions{
			Addr:     config.CacheConfig.Addr,
			Password: config.CacheConfig.Password,
			DB:       config.CacheConfig.DB,
		})

		// test Redis connection
		if err := rdb.Ping(context.Background()); err != nil {
			logger.Panic("Error connecting to Redis:", err)
		}
		cacheStorage = cache.NewCache(rdb)
	}

	if rdb != nil {
		defer rdb.Close()
	}

	store := store.NewStorage(db)

	// Initialize the rate limiter
	var rateLimiter ratelimiter.Limiter
	if config.RateLimiter.Enabled {
		rateLimiter = ratelimiter.NewFixedWindowRateLimiter(config.RateLimiter.RequestsPerTimeFrame, config.RateLimiter.TimeFrame)
	}

	appConfig := AppConfig{
		addr:   config.Server.Host,
		apiUrl: config.ApiUrl,
	}

	isProdEnv := config.Environment == "production"

	mailConfig := config.Mail

	var maxEmailRetries int

	mailer, err := mailer.NewSendGridMailer(mailConfig.SendGrid.APIKey, mailConfig.FromEmail, maxEmailRetries)
	if err != nil {
		logger.Fatal("Error initializing mailer:", err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(config.JWT.SecretKey, config.JWT.Audience, config.JWT.Issuer, int64(config.JWT.Expiry))

	app := NewApplication(appOptions{
		config:           appConfig,
		store:            store,
		db:               db,
		cache:            cacheStorage,
		logger:           logger,
		mailConfig:       mailConfig,
		mailer:           mailer,
		jwtAuthenticator: jwtAuthenticator,
		rateLimiter:      rateLimiter,
		isProdEnv:        isProdEnv,
	})

	// Metrics collection
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {

		return db.Stats()
	}))

	expvar.Publish("go_routines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	if err := app.Run(mux); err != nil {
		log.Fatal(err)
	}
}
