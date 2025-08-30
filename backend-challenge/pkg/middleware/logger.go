package middleware

import (
	"github.com/thanhfphan/kart-challenge/pkg/logging"

	"github.com/gin-gonic/gin"
)

func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		log := logging.FromContext(ctx)
		if reqID := RequestIDFromCtx(ctx); reqID != "" {
			log.SetRequestID(reqID)
		}

		ctx = logging.WithLogger(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
