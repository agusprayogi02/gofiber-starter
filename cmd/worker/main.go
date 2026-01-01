package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starter-gofiber/internal/config"
	"starter-gofiber/internal/worker"
	"starter-gofiber/pkg/logger"

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

	logger.Info("Worker server starting",
		zap.String("env", config.ENV.ENV_TYPE),
	)

	// Load database (required for some handlers)
	config.LoadDB()
	if config.ENV.DB_2_ENABLE {
		config.LoadDB2()
	}

	// Initialize Redis (required for Asynq)
	if !config.ENV.REDIS_ENABLE {
		logger.Fatal("Redis is required for worker server but is not enabled")
	}

	client, err := config.InitRedis()
	if err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer client.Close()

	// Initialize Asynq client and server
	logger.Info("Initializing Asynq worker server")
	asynqClient := config.InitAsynqClient()
	defer asynqClient.Close()

	worker.SetAsynqClient(asynqClient)
	worker.SetRedisConfig(
		config.ENV.REDIS_HOST+":"+config.ENV.REDIS_PORT,
		config.ENV.REDIS_PASSWORD,
		config.ENV.REDIS_DB,
	)

	// Initialize Asynq scheduler
	config.InitAsynqScheduler()
	defer config.CloseAsynq()

	// Register periodic tasks
	if err := worker.RegisterPeriodicTasks(config.AsynqScheduler); err != nil {
		logger.Warn("Failed to register periodic tasks", zap.Error(err))
	}

	// Initialize Asynq server with configurable concurrency
	concurrency := 10 // Can be made configurable via env var
	server := config.InitAsynqServer(concurrency)

	// Create task handler mux
	mux := asynq.NewServeMux()

	// Register legacy task handlers (non-email)
	mux.HandleFunc(worker.TaskSendEmail, worker.HandleSendEmail) // Keep for backward compatibility
	mux.HandleFunc(worker.TaskProcessExport, worker.HandleProcessExport)
	mux.HandleFunc(worker.TaskCleanupOldFiles, worker.HandleCleanupOldFiles)
	mux.HandleFunc(worker.TaskGenerateReport, worker.HandleGenerateReport)
	mux.HandleFunc(worker.TaskSendNotification, worker.HandleSendNotification)

	// Register new email handlers (with templates & SMTP)
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

	logger.Info("Worker server initialized",
		zap.Int("concurrency", concurrency),
	)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start Asynq scheduler in background
	go func() {
		logger.Info("Starting Asynq scheduler")
		if err := config.AsynqScheduler.Run(); err != nil {
			logger.Error("Asynq scheduler error", zap.Error(err))
		}
	}()

	// Start worker server in background
	go func() {
		logger.Info("Starting Asynq worker server")
		if err := server.Run(mux); err != nil {
			logger.Fatal("Failed to start Asynq worker server", zap.Error(err))
		}
	}()

	// Set default shutdown timeout if not configured
	shutdownTimeout := config.ENV.SHUTDOWN_TIMEOUT
	if shutdownTimeout == 0 {
		shutdownTimeout = 10 // Default 10 seconds
	}

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutting down worker server...",
		zap.Int("timeout_seconds", shutdownTimeout),
	)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	// Channel to track shutdown completion
	shutdownDone := make(chan error, 1)

	// Start graceful shutdown in goroutine
	go func() {
		var shutdownErrors []error

		// Shutdown Asynq scheduler first
		if config.AsynqScheduler != nil {
			logger.Info("Shutting down Asynq scheduler...")
			config.AsynqScheduler.Shutdown()
			logger.Info("Asynq scheduler shut down successfully")
		}

		// Shutdown Asynq server
		if config.AsynqServer != nil {
			logger.Info("Shutting down Asynq worker server...")
			config.AsynqServer.Shutdown()
			logger.Info("Asynq worker server shut down successfully")
		}

		// Close database connections
		if config.DB != nil {
			logger.Info("Closing database connections...")
			sqlDB, err := config.DB.DB()
			if err == nil {
				if err := sqlDB.Close(); err != nil {
					logger.Error("Error closing database", zap.Error(err))
					shutdownErrors = append(shutdownErrors, err)
				} else {
					logger.Info("Database connections closed successfully")
				}
			}
		}

		// Close Redis connections
		if config.ENV.REDIS_ENABLE {
			logger.Info("Closing Redis connections...")
			// Redis client will be closed by defer in main
			logger.Info("Redis connections closed successfully")
		}

		// Flush logger
		logger.SyncLogger()

		if len(shutdownErrors) > 0 {
			shutdownDone <- fmt.Errorf("shutdown completed with %d errors", len(shutdownErrors))
		} else {
			shutdownDone <- nil
		}
	}()

	// Wait for shutdown completion or timeout
	select {
	case err := <-shutdownDone:
		if err != nil {
			logger.Error("Shutdown completed with errors", zap.Error(err))
		} else {
			logger.Info("Worker server shut down gracefully")
		}
	case <-ctx.Done():
		logger.Warn("Shutdown timeout exceeded, forcing exit",
			zap.Int("timeout_seconds", shutdownTimeout),
		)
	}

	logger.Info("Worker server exited")
}
