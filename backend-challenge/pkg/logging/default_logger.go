package logging

import (
	"context"

	"github.com/google/uuid"
)

func FromContext(ctx context.Context) *ContextLogger {
	if logger, ok := ctx.Value(loggerKey).(*ContextLogger); ok {
		return logger
	}

	return createLoggerWithRandomRequestID()
}

func WithLogger(ctx context.Context, logger *ContextLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func NewContextLogger() *ContextLogger {
	return createLoggerWithRandomRequestID()
}

func createLoggerWithRandomRequestID() *ContextLogger {
	return &ContextLogger{
		reqInfo: RequestInfo{
			RequestID: uuid.NewString(),
		},
	}
}
