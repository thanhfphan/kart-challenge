package logging

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ctxKey string

const loggerKey = ctxKey("logger_key")

type RequestInfo struct {
	RequestID string
}

type ContextLogger struct {
	reqInfo RequestInfo
}

func init() {
	handler := newHandler()
	slog.SetDefault(slog.New(handler))
}

func newHandler() slog.Handler {
	if os.Getenv("ENVIRONMENT") == "local" {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevel(os.Getenv("LOG_LEVEL")),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "request_id" {
					return slog.Attr{}
				} else if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339))
				}

				return a
			},
		})
	}

	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(os.Getenv("LOG_LEVEL")),
	})
}

func getLogLevel(logLevel string) slog.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func (c *ContextLogger) SetRequestID(reqID string) {
	c.reqInfo.RequestID = reqID
}

func (c *ContextLogger) Infof(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.withBasicFields().Info(msg)
}

func (c *ContextLogger) Info(msg string, args ...interface{}) {
	c.withBasicFields().Info(msg, args...)
}

func (c *ContextLogger) Debugf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.withBasicFields().Debug(msg)
}

func (c *ContextLogger) Debug(msg string, args ...interface{}) {
	c.withBasicFields().Debug(msg, args...)
}

func (c *ContextLogger) Warn(msg string, args ...interface{}) {
	c.withBasicFields().Warn(msg, args...)
}

func (c *ContextLogger) Warnf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.withBasicFields().Warn(msg)
}

func (c *ContextLogger) Error(msg string, args ...interface{}) {
	c.withBasicFields().Error(msg, args...)
}

func (c *ContextLogger) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c.withBasicFields().Error(msg)
}

func (c *ContextLogger) Fatalf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c.withBasicFields().Error(msg)
}

func (c *ContextLogger) withBasicFields() *slog.Logger {
	logger := slog.With(slog.String("caller", getCaller()))
	if c.reqInfo.RequestID != "" {
		logger = logger.With("request_id", c.reqInfo.RequestID)
	}

	return logger
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}

	lastSlashIndex := strings.LastIndex(file, "/")
	secondLastSlashIndex := strings.LastIndex(file[:lastSlashIndex], "/")
	if secondLastSlashIndex == -1 {
		return file[lastSlashIndex+1:] + ":" + strconv.Itoa(line) // Just the file name if no parent directory
	}

	parentDir := file[secondLastSlashIndex+1 : lastSlashIndex]
	shortFile := file[lastSlashIndex+1:]

	return parentDir + "/" + shortFile + ":" + strconv.Itoa(line)
}
