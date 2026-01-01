package main

import (
	"os"
	"os/signal"
	"syscall"

	"starter-gofiber/config"
	"starter-gofiber/helper"
	"starter-gofiber/jobs"
	"starter-gofiber/middleware"
	"starter-gofiber/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
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

			// Initialize Asynq for background jobs (Redis required)
			helper.Info("Initializing Asynq job queue")
			asynqClient := config.InitAsynqClient()
			helper.SetAsynqClient(asynqClient)
			helper.SetRedisConfig(
				config.ENV.REDIS_HOST+":"+config.ENV.REDIS_PORT,
				config.ENV.REDIS_PASSWORD,
				config.ENV.REDIS_DB,
			)

			config.InitAsynqScheduler()
			defer config.CloseAsynq()

			// Register periodic tasks
			if err := jobs.RegisterPeriodicTasks(config.AsynqScheduler); err != nil {
				helper.Warn("Failed to register periodic tasks", zap.Error(err))
			}

			// Start Asynq worker server in background
			go startWorkerServer()

			// Start Asynq scheduler
			go func() {
				if err := config.AsynqScheduler.Run(); err != nil {
					helper.Error("Asynq scheduler error", zap.Error(err))
				}
			}()
		}
	}

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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + config.ENV.PORT); err != nil {
			helper.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-quit
	helper.Info("Shutting down server...")

	// Shutdown Asynq server
	if config.AsynqServer != nil {
		helper.Info("Shutting down Asynq worker server...")
		config.AsynqServer.Shutdown()
	}

	// Shutdown Fiber app
	if err := app.Shutdown(); err != nil {
		helper.Error("Error shutting down server", zap.Error(err))
	}

	helper.Info("Server exited")
}

// startWorkerServer starts Asynq worker server to process background jobs
func startWorkerServer() {
	helper.Info("Starting Asynq worker server")

	// Initialize server with 10 concurrent workers
	server := config.InitAsynqServer(10)

	// Create task handler mux
	mux := asynq.NewServeMux()

	// Register task handlers
	mux.HandleFunc(helper.TaskSendEmail, jobs.HandleSendEmail)
	mux.HandleFunc(helper.TaskSendVerificationCode, jobs.HandleSendVerificationEmail)
	mux.HandleFunc(helper.TaskSendPasswordReset, jobs.HandleSendPasswordReset)
	mux.HandleFunc(helper.TaskProcessExport, jobs.HandleProcessExport)
	mux.HandleFunc(helper.TaskCleanupOldFiles, jobs.HandleCleanupOldFiles)
	mux.HandleFunc(helper.TaskGenerateReport, jobs.HandleGenerateReport)
	mux.HandleFunc(helper.TaskSendNotification, jobs.HandleSendNotification)

	// Register periodic task handlers
	mux.HandleFunc("system:health_check", jobs.HandleHealthCheck)
	mux.HandleFunc("cleanup:expired_tokens", jobs.HandleCleanupExpiredTokens)
	mux.HandleFunc("archive:monthly", jobs.HandleMonthlyArchive)
	mux.HandleFunc("backup:database", jobs.HandleDatabaseBackup)
	mux.HandleFunc("email:daily_digest", jobs.HandleDailyEmailDigest)
	mux.HandleFunc("metrics:collect", jobs.HandleMetricsCollection)

	// Start server
	if err := server.Run(mux); err != nil {
		helper.Fatal("Failed to start Asynq worker server", zap.Error(err))
	}
}
