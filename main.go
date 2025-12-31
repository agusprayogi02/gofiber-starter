package main

import (
	"database/sql"
	"errors"

	"starter-gofiber/config"
	"starter-gofiber/helper"
	"starter-gofiber/middleware"
	"starter-gofiber/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig() // required first, because it will load .env file

	// Initialize structured logging
	if err := helper.InitLogger(config.ENV.ENV_TYPE); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer helper.SyncLogger()

	helper.Info("Application starting",
		zap.String("env", config.ENV.ENV_TYPE),
		zap.String("port", config.ENV.PORT),
	)

	// Initialize RSA private key
	if err := helper.InitPrivateKey(config.ENV.LOCATION_CERT); err != nil {
		helper.Fatal("Failed to initialize private key", zap.Error(err))
	}

	// Initialize encryption for sensitive data
	if err := helper.InitEncryption(config.ENV.ENCRYPTION_KEY); err != nil {
		helper.Fatal("Failed to initialize encryption", zap.Error(err))
	}

	// Initialize Sentry for error tracking
	if err := helper.InitSentry(config.ENV.SENTRY_DSN, config.ENV.ENV_TYPE); err != nil {
		helper.Warn("Failed to initialize Sentry", zap.Error(err))
	}
	defer helper.FlushSentry()

	config.LoadTimezone()
	config.LoadPermissions()
	config.LoadStorage()
	config.LoadDB()
	if config.ENV.DB_2_ENABLE {
		config.LoadDB2()
	}

	// Initialize Redis cache
	if config.ENV.REDIS_ENABLE {
		client, err := config.InitRedis()
		if err != nil {
			helper.Warn("Failed to initialize Redis", zap.Error(err))
		} else {
			helper.InitRedisClient(client)
			defer helper.CloseRedis()
		}
	}

	// Start database metrics updater
	helper.StartDBMetricsUpdater(func() (*sql.DB, error) {
		if config.DB == nil {
			return nil, errors.New("database not initialized")
		}
		return config.DB.DB()
	})

	// Initialize API Key middleware
	middleware.InitAPIKeyMiddleware(config.DB)

	conf := fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: helper.ErrorHelper,
	}
	if config.ENV.ENV_TYPE == "prod" {
		conf.Prefork = true
	}

	app := fiber.New(conf)
	config.App(app)
	router.AppRouter(app)

	helper.Info("Server starting", zap.String("port", config.ENV.PORT))

	err := app.Listen(":" + config.ENV.PORT)
	if err != nil {
		helper.Fatal("Failed to start server", zap.Error(err))
	}
}
