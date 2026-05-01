package middleware

import (
	"net/http"
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
		log.Info(path,
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
				log.Error("panic recovered",
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
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID")

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
			traceID = "trace-" + time.Now().Format("20060102150405") + "-" + randomString(8)
		}
		c.Set("traceID", traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[int(time.Now().UnixNano())%len(letters)]
	}
	return string(result)
}
