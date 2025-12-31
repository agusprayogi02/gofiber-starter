package helper

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	// Auth metrics
	authAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"result"},
	)

	authTokensCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_tokens_created_total",
			Help: "Total number of authentication tokens created",
		},
	)

	// DB metrics
	dbConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_in_use",
			Help: "Number of database connections currently in use",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)
)

// RecordAuthAttempt records authentication attempts to Prometheus
func RecordAuthAttempt(success bool) {
	result := "failure"
	if success {
		result = "success"
		authTokensCreated.Inc()
	}
	authAttemptsTotal.WithLabelValues(result).Inc()
}

// UpdateDBMetrics updates database connection pool metrics
func UpdateDBMetrics(inUse, idle int) {
	dbConnectionsInUse.Set(float64(inUse))
	dbConnectionsIdle.Set(float64(idle))
}

// StartDBMetricsUpdater starts a goroutine to periodically update database metrics
// Pass *sql.DB from config.DB.DB() to avoid import cycle
func StartDBMetricsUpdater(getDB func() (*sql.DB, error)) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			db, err := getDB()
			if err != nil {
				Error("Failed to get database connection", zap.Error(err))
				continue
			}

			stats := db.Stats()
			UpdateDBMetrics(stats.InUse, stats.Idle)
		}
	}()

	Info("Database metrics updater started")
}
