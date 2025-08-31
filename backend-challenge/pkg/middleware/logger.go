package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
)

// Logger adds a logger with unique request ID to the request context
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		logger := logging.NewContextLogger()
		logger.SetRequestID(requestID)

		ctx := logging.WithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", requestID)

		logger.Infof("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		logger.Infof("Request completed with status: %d", c.Writer.Status())
	}
}
