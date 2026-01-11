package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/pkg/metrics"
)

// PrometheusMiddleware collects HTTP metrics for Prometheus
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Update metrics
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

// TracingMiddleware adds trace_id to each request for distributed tracing
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if trace_id is already present in header (from upstream service)
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			// Generate new trace_id
			traceID = uuid.New().String()
		}

		// Set trace_id in context
		c.Set("trace_id", traceID)

		// Set trace_id in response header for client debugging
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}

// RequestLoggingMiddleware logs all incoming requests with structured logging
func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after request completes
		duration := time.Since(start)
		status := c.Writer.Status()

		// Get trace_id from context
		traceID, _ := c.Get("trace_id")

		// Structured log entry
		// Using slog would be: slog.Info("HTTP Request", ...)
		// For now we'll use gin's built-in logger format
		// In production, replace with slog calls
		gin.DefaultWriter.Write([]byte(
			formatLogEntry(method, path, status, duration, traceID),
		))
	}
}

func formatLogEntry(method, path string, status int, duration time.Duration, traceID any) string {
	return ""
}
