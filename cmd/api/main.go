package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"starter-gofiber/internal/config"
	"starter-gofiber/internal/handler/middleware"
	"starter-gofiber/internal/infrastructure/cache"
	"starter-gofiber/internal/infrastructure/email"
	"starter-gofiber/internal/worker"
	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/pkg/logger"
	"starter-gofiber/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig() // required first, because it will load .env file

	// Initialize structured logging
	if err := logger.InitLogger(config.ENV.ENV_TYPE); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.SyncLogger()

	logger.Info("Application starting",
		zap.String("env", config.ENV.ENV_TYPE),
		zap.String("port", config.ENV.PORT),
	)

	// Initialize RSA private key
	// Get project root path and join with certificate location
	projectRoot, err := os.Getwd()
	if err != nil {
		logger.Fatal("Failed to get working directory", zap.Error(err))
	}
	certPath := filepath.Join(projectRoot, config.ENV.LOCATION_CERT)
	if err := crypto.InitPrivateKey(certPath); err != nil {
		logger.Fatal("Failed to initialize private key", zap.Error(err))
	}

	// Initialize encryption for sensitive data
	if err := crypto.InitEncryption(config.ENV.ENCRYPTION_KEY); err != nil {
		logger.Fatal("Failed to initialize encryption", zap.Error(err))
	}

	// Initialize Sentry for error tracking
	if err := logger.InitSentry(config.ENV.SENTRY_DSN, config.ENV.ENV_TYPE); err != nil {
		logger.Warn("Failed to initialize Sentry", zap.Error(err))
	}
	defer logger.FlushSentry()

	// Initialize email configuration
	if err := email.InitEmail(); err != nil {
		logger.Warn("Failed to initialize email config", zap.Error(err))
	}

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
			logger.Warn("Failed to initialize Redis", zap.Error(err))
		} else {
			cache.InitRedisClient(client)
			defer cache.CloseRedis()

			// Initialize Asynq for background jobs (Redis required)
			logger.Info("Initializing Asynq job queue")
			asynqClient := config.InitAsynqClient()
			worker.SetAsynqClient(asynqClient)
			worker.SetRedisConfig(
				config.ENV.REDIS_HOST+":"+config.ENV.REDIS_PORT,
				config.ENV.REDIS_PASSWORD,
				config.ENV.REDIS_DB,
			)

			config.InitAsynqScheduler()
			defer config.CloseAsynq()

			// Register periodic tasks
			if err := worker.RegisterPeriodicTasks(config.AsynqScheduler); err != nil {
				logger.Warn("Failed to register periodic tasks", zap.Error(err))
			}

			// Start Asynq worker server in background
			go startWorkerServer()

			// Start Asynq scheduler
			go func() {
				if err := config.AsynqScheduler.Run(); err != nil {
					logger.Error("Asynq scheduler error", zap.Error(err))
				}
			}()
		}
	}

	// Initialize API Key middleware
	middleware.InitAPIKeyMiddleware(config.DB)

	conf := fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: apierror.ErrorHelper,
	}
	if config.ENV.ENV_TYPE == "prod" {
		conf.Prefork = true
	}

	app := fiber.New(conf)
	config.App(app)
	router.AppRouter(app)

	logger.Info("Server starting", zap.String("port", config.ENV.PORT))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + config.ENV.PORT); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutting down server...")

	// Shutdown Asynq server
	if config.AsynqServer != nil {
		logger.Info("Shutting down Asynq worker server...")
		config.AsynqServer.Shutdown()
	}

	// Shutdown Fiber app
	if err := app.Shutdown(); err != nil {
		logger.Error("Error shutting down server", zap.Error(err))
	}

	logger.Info("Server exited")
}

// startWorkerServer starts Asynq worker server to process background jobs
func startWorkerServer() {
	logger.Info("Starting Asynq worker server")

	// Initialize server with 10 concurrent workers
	server := config.InitAsynqServer(10)

	// Create task handler mux
	mux := asynq.NewServeMux()

	// Register legacy task handlers (non-email)
	mux.HandleFunc(worker.TaskSendEmail, worker.HandleSendEmail) // Keep for backward compatibility
	mux.HandleFunc(worker.TaskProcessExport, worker.HandleProcessExport)
	mux.HandleFunc(worker.TaskCleanupOldFiles, worker.HandleCleanupOldFiles)
	mux.HandleFunc(worker.TaskGenerateReport, worker.HandleGenerateReport)
	mux.HandleFunc(worker.TaskSendNotification, worker.HandleSendNotification)

	// Register new email handlers (with templates & SMTP)
	// Note: These replace legacy email handlers to avoid duplicate registration
	mux.HandleFunc(worker.TypeEmailWelcome, worker.HandleEmailWelcome)
	mux.HandleFunc(worker.TypeEmailPasswordReset, worker.HandleEmailPasswordReset)
	mux.HandleFunc(worker.TypeEmailVerification, worker.HandleEmailVerification)
	mux.HandleFunc(worker.TypeEmailCustom, worker.HandleEmailCustom)

	// Register periodic task handlers
	mux.HandleFunc("system:health_check", worker.HandleHealthCheck)
	mux.HandleFunc("cleanup:expired_tokens", worker.HandleCleanupExpiredTokens)
	mux.HandleFunc("archive:monthly", worker.HandleMonthlyArchive)
	mux.HandleFunc("backup:database", worker.HandleDatabaseBackup)
	mux.HandleFunc("email:daily_digest", worker.HandleDailyEmailDigest)
	mux.HandleFunc("metrics:collect", worker.HandleMetricsCollection)

	// Start server
	if err := server.Run(mux); err != nil {
		logger.Fatal("Failed to start Asynq worker server", zap.Error(err))
	}
}
