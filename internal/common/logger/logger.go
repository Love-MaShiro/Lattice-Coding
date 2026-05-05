package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lattice-coding/internal/common/config"
)

type Logger struct {
	*zap.Logger
}

type LogContext struct {
	TraceID   string
	RunID     string
	SessionID string
	Provider  string
	Model     string
}

type contextKey string

const (
	logContextKey contextKey = "log_context"
)

func NewLogger(cfg *config.Config) *Logger {
	var config zap.Config

	if cfg.App.Env == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level = zap.NewAtomicLevelAt(getLevel(cfg.Logging.Level))
	config.Encoding = cfg.Logging.Format
	config.OutputPaths = []string{cfg.Logging.Output}

	logger, err := config.Build()
	if err != nil {
		zap.L().Fatal("failed to initialize logger", zap.Error(err))
		os.Exit(1)
	}

	return &Logger{logger}
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}

func (l *Logger) WithTrace(traceID string) *Logger {
	return &Logger{l.With(zap.String("trace_id", traceID))}
}

func (l *Logger) WithRun(runID string) *Logger {
	return &Logger{l.With(zap.String("run_id", runID))}
}

func (l *Logger) WithSession(sessionID string) *Logger {
	return &Logger{l.With(zap.String("session_id", sessionID))}
}

func (l *Logger) WithProvider(provider string) *Logger {
	return &Logger{l.With(zap.String("provider", provider))}
}

func (l *Logger) WithModel(model string) *Logger {
	return &Logger{l.With(zap.String("model", model))}
}

func (l *Logger) WithCtx(ctx context.Context) *Logger {
	if logCtx := GetContext(ctx); logCtx != nil {
		return &Logger{l.With(
			zap.String("trace_id", logCtx.TraceID),
			zap.String("run_id", logCtx.RunID),
			zap.String("session_id", logCtx.SessionID),
			zap.String("provider", logCtx.Provider),
			zap.String("model", logCtx.Model),
		)}
	}
	return l
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{l.With(fields...)}
}

func (l *Logger) ErrorWithStack(err error, msg string, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Error(msg, fields...)
}

func (l *Logger) ErrorWithCtx(ctx context.Context, err error, msg string, fields ...zap.Field) {
	logger := l.WithCtx(ctx)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	logger.Error(msg, fields...)
}

func (l *Logger) InfoWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithCtx(ctx).Info(msg, fields...)
}

func (l *Logger) DebugWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithCtx(ctx).Debug(msg, fields...)
}

func (l *Logger) WarnWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithCtx(ctx).Warn(msg, fields...)
}

func (l *Logger) LLMCall(ctx context.Context, provider, model string, latencyMs int64, tokens int, err error) {
	logger := l.WithCtx(ctx).WithProvider(provider).WithModel(model)
	fields := []zap.Field{
		zap.Int64("latency_ms", latencyMs),
		zap.Int("tokens", tokens),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("LLM call failed", fields...)
	} else {
		logger.Info("LLM call success", fields...)
	}
}

func (l *Logger) ToolCall(ctx context.Context, toolName string, latencyMs int64, err error) {
	logger := l.WithCtx(ctx)
	fields := []zap.Field{
		zap.String("tool", toolName),
		zap.Int64("latency_ms", latencyMs),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("Tool call failed", fields...)
	} else {
		logger.Info("Tool call success", fields...)
	}
}

func (l *Logger) RunEvent(ctx context.Context, runID, eventType, status string, details map[string]interface{}) {
	logger := l.WithCtx(ctx).WithRun(runID)
	fields := []zap.Field{
		zap.String("event_type", eventType),
		zap.String("status", status),
	}
	if details != nil {
		for k, v := range details {
			fields = append(fields, zap.Any(k, v))
		}
	}
	logger.Info("Run event", fields...)
}

func NewContext(ctx context.Context, logCtx *LogContext) context.Context {
	return context.WithValue(ctx, logContextKey, logCtx)
}

func GetContext(ctx context.Context) *LogContext {
	if v := ctx.Value(logContextKey); v != nil {
		return v.(*LogContext)
	}
	return nil
}

func GetTraceID(ctx context.Context) string {
	if logCtx := GetContext(ctx); logCtx != nil {
		return logCtx.TraceID
	}
	return ""
}

func GetRunID(ctx context.Context) string {
	if logCtx := GetContext(ctx); logCtx != nil {
		return logCtx.RunID
	}
	return ""
}

func GetSessionID(ctx context.Context) string {
	if logCtx := GetContext(ctx); logCtx != nil {
		return logCtx.SessionID
	}
	return ""
}
