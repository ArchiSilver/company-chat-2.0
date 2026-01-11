package http

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db          *sqlx.DB
	redisClient *redis.Client
	startTime   time.Time
}

func NewHealthHandler(db *sqlx.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		startTime:   time.Now(),
	}
}

func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
	// Health endpoints (no auth required)
	router.GET("/health", h.HealthCheck)
	router.GET("/health/ready", h.ReadinessCheck)
	router.GET("/health/live", h.LivenessCheck)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// HealthCheck returns overall health status
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := "healthy"
	httpStatus := http.StatusOK

	// Check database
	dbStatus := "up"
	if err := h.db.Ping(); err != nil {
		dbStatus = "down"
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	// Check Redis
	redisStatus := "up"
	if _, err := h.redisClient.Ping(c.Request.Context()).Result(); err != nil {
		redisStatus = "down"
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	// Memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
		"version":   "1.0.0",
		"checks": gin.H{
			"database": gin.H{
				"status": dbStatus,
			},
			"redis": gin.H{
				"status": redisStatus,
			},
		},
		"system": gin.H{
			"goroutines":    runtime.NumGoroutine(),
			"memory_alloc":  memStats.Alloc / 1024 / 1024,      // MB
			"memory_sys":    memStats.Sys / 1024 / 1024,        // MB
			"gc_cycles":     memStats.NumGC,
		},
	})
}

// ReadinessCheck checks if the service is ready to accept traffic
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check database connection
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database connection failed",
		})
		return
	}

	// Check Redis connection
	if _, err := h.redisClient.Ping(c.Request.Context()).Result(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "redis connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
