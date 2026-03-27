package middleware

import (
	"log/slog"
	"time"

	"bsky-schwarz/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()[:8]
		}
		c.Set("request_id", requestID)

		log := logger.With("request_id", requestID)
		c.Set("logger", log)

		start := time.Now()

		log.Info("request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"remote_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logLevel := "info"
		msg := "request completed"
		if status >= 400 && status < 500 {
			logLevel = "warn"
			msg = "request failed"
		} else if status >= 500 {
			logLevel = "error"
			msg = "request failed"
		}

		switch logLevel {
		case "warn":
			log.Warn(msg,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"response_size", c.Writer.Size(),
			)
		case "error":
			log.Error(msg,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"response_size", c.Writer.Size(),
			)
		default:
			log.Info(msg,
				"status", status,
				"duration_ms", duration.Milliseconds(),
				"response_size", c.Writer.Size(),
			)
		}
	}
}

func GetLogger(c *gin.Context) *slog.Logger {
	if log, exists := c.Get("logger"); exists {
		return log.(*slog.Logger)
	}
	return logger.Log
}

func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return "unknown"
}
