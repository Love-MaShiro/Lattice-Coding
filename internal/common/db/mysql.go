package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lattice-coding/internal/common/config"
)

func NewMySQL(cfg *config.MySQLConfig) (*gorm.DB, error) {
	dsn, err := buildMySQLDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build mysql dsn: %w", err)
	}

	logLevel, err := parseGormLogLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid mysql log_level: %w", err)
	}

	gormLogger := logger.Default.LogMode(logLevel)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if cfg.ConnMaxLifetime != "" {
		lifetime, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err != nil {
			return nil, fmt.Errorf("invalid conn_max_lifetime: %w", err)
		}
		sqlDB.SetConnMaxLifetime(lifetime)
	}

	if cfg.ConnMaxIdleTime != "" {
		idleTime, err := time.ParseDuration(cfg.ConnMaxIdleTime)
		if err != nil {
			return nil, fmt.Errorf("invalid conn_max_idle_time: %w", err)
		}
		sqlDB.SetConnMaxIdleTime(idleTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return db, nil
}

func buildMySQLDSN(cfg *config.MySQLConfig) (string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		url.PathEscape(cfg.Username),
		url.PathEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		url.PathEscape(cfg.Database),
		cfg.Charset,
	)
	return dsn, nil
}

func parseGormLogLevel(level string) (logger.LogLevel, error) {
	switch level {
	case "silent", "Silent", "SILENT", "0":
		return logger.Silent, nil
	case "error", "Error", "ERROR", "1":
		return logger.Error, nil
	case "warn", "Warn", "WARN", "2":
		return logger.Warn, nil
	case "info", "Info", "INFO", "3":
		return logger.Info, nil
	case "":
		return logger.Silent, nil
	default:
		return logger.Silent, fmt.Errorf("unknown log level: %s", level)
	}
}
