package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/common/response"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		log.WithCtx(c.Request.Context()).Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
			zap.Int("length", c.Writer.Size()),
		)
	}
}

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithCtx(c.Request.Context()).Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
				)
				response.FailWithStatus(c, http.StatusInternalServerError, errors.Internal("系统内部错误"))
			}
		}()
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err.Err)
			return
		}

		if gin.IsDebugging() {
			return
		}

		if c.Writer.Status() == http.StatusInternalServerError {
			response.Fail(c, errors.Internal("系统内部错误"))
		}
	}
}

func handleError(c *gin.Context, err error) {
	if bizErr, ok := err.(*errors.BizError); ok {
		switch bizErr.Code {
		case errors.Unauthorized:
			response.FailWithStatus(c, http.StatusUnauthorized, bizErr)
		case errors.Forbidden:
			response.FailWithStatus(c, http.StatusForbidden, bizErr)
		case errors.NotFound:
			response.FailWithStatus(c, http.StatusNotFound, bizErr)
		default:
			response.Fail(c, bizErr)
		}
		return
	}

	response.Fail(c, errors.InternalWithErr(err, "系统内部错误"))
}

func CORS() gin.HandlerFunc {
	allowOrigins := getEnvOrDefault("CORS_ALLOW_ORIGINS", "*")
	allowMethods := getEnvOrDefault("CORS_ALLOW_METHODS", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	allowHeaders := getEnvOrDefault("CORS_ALLOW_HEADERS", "Content-Type, Authorization, X-Trace-ID, X-Run-ID, X-Session-ID, X-Request-ID")
	exposeHeaders := getEnvOrDefault("CORS_EXPOSE_HEADERS", "X-Trace-ID, X-Run-ID, X-Session-ID")
	allowCredentials := getEnvOrDefault("CORS_ALLOW_CREDENTIALS", "true")
	maxAge := getEnvOrDefault("CORS_MAX_AGE", "86400")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" && allowOrigins != "*" {
			if strings.Contains(allowOrigins, origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else if strings.Contains(allowOrigins, "http://localhost") || strings.Contains(allowOrigins, "http://127.0.0.1") {
				if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		c.Writer.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", allowCredentials)
		c.Writer.Header().Set("Access-Control-Max-Age", maxAge)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.FailWithStatus(c, http.StatusUnauthorized, errors.UnauthorizedErr("未授权访问"))
			return
		}
		c.Next()
	}
}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}

		runID := c.GetHeader("X-Run-ID")
		sessionID := c.GetHeader("X-Session-ID")

		logCtx := &logger.LogContext{
			TraceID:   traceID,
			RunID:     runID,
			SessionID: sessionID,
		}
		ctx := logger.NewContext(c.Request.Context(), logCtx)
		c.Request = c.Request.WithContext(ctx)

		c.Set("traceID", traceID)
		c.Set("runID", runID)
		c.Set("sessionID", sessionID)
		c.Writer.Header().Set("X-Trace-ID", traceID)

		if runID != "" {
			c.Writer.Header().Set("X-Run-ID", runID)
		}
		if sessionID != "" {
			c.Writer.Header().Set("X-Session-ID", sessionID)
		}

		c.Next()
	}
}

func GetTraceID(c *gin.Context) string {
	if v, exists := c.Get("traceID"); exists {
		return v.(string)
	}
	return ""
}

func GetRunID(c *gin.Context) string {
	if v, exists := c.Get("runID"); exists {
		return v.(string)
	}
	return ""
}

func GetSessionID(c *gin.Context) string {
	if v, exists := c.Get("sessionID"); exists {
		return v.(string)
	}
	return ""
}

func generateTraceID() string {
	return "trace-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(result)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}
