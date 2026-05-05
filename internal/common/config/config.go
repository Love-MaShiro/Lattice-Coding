package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	MySQL    MySQLConfig    `yaml:"mysql"`
	Redis    RedisConfig    `yaml:"redis"`
	Postgres PostgresConfig `yaml:"postgres"`
	LLM      LLMConfig      `yaml:"llm"`
	Logging  LoggingConfig  `yaml:"logging"`
	HTTP     HTTPConfig     `yaml:"http"`
	Worker   WorkerConfig   `yaml:"worker"`
}

type AppConfig struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Port  string `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime string `yaml:"conn_max_idle_time"`
	LogLevel        string `yaml:"log_level"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	DialTimeout  string `yaml:"dial_timeout"`
	PoolTimeout  string `yaml:"pool_timeout"`
	MaxRetries   int    `yaml:"max_retries"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

type LLMConfig struct {
	Pool           PoolConfig           `yaml:"pool"`
	Stream         PoolConfig           `yaml:"stream"`
	Timeout        TimeoutConfig        `yaml:"timeout"`
	Retry          RetryConfig          `yaml:"retry"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker"`
	Routing        RoutingConfig        `yaml:"routing"`
}

type TimeoutConfig struct {
	SyncCall    string `yaml:"sync_call"`
	StreamIdle  string `yaml:"stream_idle"`
	HealthCheck string `yaml:"health_check"`
}

type RetryConfig struct {
	MaxAttempts     int    `yaml:"max_attempts"`
	InitialInterval string `yaml:"initial_interval"`
	MaxInterval     string `yaml:"max_interval"`
}

type PoolConfig struct {
	MaxConcurrent  int    `yaml:"max_concurrent"`
	AcquireTimeout string `yaml:"acquire_timeout"`
}

type CircuitBreakerConfig struct {
	Enabled               bool    `yaml:"enabled"`
	Window                string  `yaml:"window"`
	MinRequests           int     `yaml:"min_requests"`
	FailureRateThreshold  float64 `yaml:"failure_rate_threshold"`
	SlowCallThreshold     string  `yaml:"slow_call_threshold"`
	SlowCallRateThreshold float64 `yaml:"slow_call_rate_threshold"`
	OpenDuration          string  `yaml:"open_duration"`
	HalfOpenMaxRequests   int     `yaml:"half_open_max_requests"`
}

type RoutingConfig struct {
	Default RouteConfig `yaml:"default"`
}

type RouteConfig struct {
	Primary  string   `yaml:"primary"`
	Fallback []string `yaml:"fallback"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

type HTTPConfig struct {
	Port         string `yaml:"port"`
	Timeout      string `yaml:"timeout"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

type WorkerConfig struct {
	Concurrency int    `yaml:"concurrency"`
	QueueName   string `yaml:"queue_name"`
}

func LoadConfig() (*Config, error) {
	return LoadConfigWithPath("configs/config.yaml")
}

func LoadConfigWithPath(paths ...string) (*Config, error) {
	var cfg Config

	for _, path := range paths {
		if err := loadConfigFile(path, &cfg); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
	}

	if err := applyEnvOverrides(&cfg); err != nil {
		return nil, err
	}

	if err := setDefaults(&cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadConfigFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

func applyEnvOverrides(cfg *Config) error {
	if env := getEnv("APP_NAME"); env != "" {
		cfg.App.Name = env
	}
	if env := getEnv("APP_ENV"); env != "" {
		cfg.App.Env = env
	}
	if env := getEnv("APP_PORT"); env != "" {
		cfg.App.Port = env
	}

	if env := getEnv("MYSQL_HOST"); env != "" {
		cfg.MySQL.Host = env
	}
	if env := getEnv("MYSQL_PORT"); env != "" {
		cfg.MySQL.Port = parseInt(env)
	}
	if env := getEnv("MYSQL_USERNAME"); env != "" {
		cfg.MySQL.Username = env
	}
	if env := getEnv("MYSQL_PASSWORD"); env != "" {
		cfg.MySQL.Password = env
	}
	if env := getEnv("MYSQL_DATABASE"); env != "" {
		cfg.MySQL.Database = env
	}

	if env := getEnv("REDIS_HOST"); env != "" {
		cfg.Redis.Host = env
	}
	if env := getEnv("REDIS_PORT"); env != "" {
		cfg.Redis.Port = parseInt(env)
	}
	if env := getEnv("REDIS_PASSWORD"); env != "" {
		cfg.Redis.Password = env
	}

	if env := getEnv("POSTGRES_HOST"); env != "" {
		cfg.Postgres.Host = env
	}
	if env := getEnv("POSTGRES_PORT"); env != "" {
		cfg.Postgres.Port = parseInt(env)
	}
	if env := getEnv("POSTGRES_USERNAME"); env != "" {
		cfg.Postgres.Username = env
	}
	if env := getEnv("POSTGRES_PASSWORD"); env != "" {
		cfg.Postgres.Password = env
	}
	if env := getEnv("POSTGRES_DATABASE"); env != "" {
		cfg.Postgres.Database = env
	}

	if env := getEnv("HTTP_PORT"); env != "" {
		cfg.HTTP.Port = env
	}

	if env := getEnv("LOGGING_LEVEL"); env != "" {
		cfg.Logging.Level = env
	}

	return nil
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func parseInt(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}

func setDefaults(cfg *Config) error {
	if cfg.App.Name == "" {
		cfg.App.Name = "lattice-coding"
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "development"
	}
	if cfg.App.Port == "" {
		cfg.App.Port = "8080"
	}

	if cfg.MySQL.Host == "" {
		cfg.MySQL.Host = "localhost"
	}
	if cfg.MySQL.Port == 0 {
		cfg.MySQL.Port = 3306
	}
	if cfg.MySQL.Charset == "" {
		cfg.MySQL.Charset = "utf8mb4"
	}
	if cfg.MySQL.MaxOpenConns == 0 {
		cfg.MySQL.MaxOpenConns = 20
	}
	if cfg.MySQL.MaxIdleConns == 0 {
		cfg.MySQL.MaxIdleConns = 10
	}
	if cfg.MySQL.ConnMaxLifetime == "" {
		cfg.MySQL.ConnMaxLifetime = "300s"
	}
	if cfg.MySQL.ConnMaxIdleTime == "" {
		cfg.MySQL.ConnMaxIdleTime = "600s"
	}

	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "localhost"
	}
	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 20
	}
	if cfg.Redis.MinIdleConns == 0 {
		cfg.Redis.MinIdleConns = 5
	}
	if cfg.Redis.DialTimeout == "" {
		cfg.Redis.DialTimeout = "5s"
	}
	if cfg.Redis.PoolTimeout == "" {
		cfg.Redis.PoolTimeout = "4s"
	}

	if cfg.Postgres.Host == "" {
		cfg.Postgres.Host = "localhost"
	}
	if cfg.Postgres.Port == 0 {
		cfg.Postgres.Port = 5432
	}
	if cfg.Postgres.SSLMode == "" {
		cfg.Postgres.SSLMode = "disable"
	}

	if cfg.HTTP.Port == "" {
		cfg.HTTP.Port = cfg.App.Port
	}
	if cfg.HTTP.Timeout == "" {
		cfg.HTTP.Timeout = "30s"
	}
	if cfg.HTTP.ReadTimeout == "" {
		cfg.HTTP.ReadTimeout = "30s"
	}
	if cfg.HTTP.WriteTimeout == "" {
		cfg.HTTP.WriteTimeout = "60s"
	}
	if cfg.HTTP.IdleTimeout == "" {
		cfg.HTTP.IdleTimeout = "120s"
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "json"
	}
	if cfg.Logging.Output == "" {
		cfg.Logging.Output = "stdout"
	}
	if cfg.Logging.MaxSize == 0 {
		cfg.Logging.MaxSize = 100
	}
	if cfg.Logging.MaxBackups == 0 {
		cfg.Logging.MaxBackups = 3
	}
	if cfg.Logging.MaxAge == 0 {
		cfg.Logging.MaxAge = 7
	}

	if cfg.Worker.Concurrency == 0 {
		cfg.Worker.Concurrency = 10
	}
	if cfg.Worker.QueueName == "" {
		cfg.Worker.QueueName = "agent_run"
	}

	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.MySQL.Username == "" {
		return fmt.Errorf("mysql.username is required")
	}
	if cfg.MySQL.Database == "" {
		return fmt.Errorf("mysql.database is required")
	}

	if cfg.Postgres.Username == "" {
		return fmt.Errorf("postgres.username is required")
	}
	if cfg.Postgres.Database == "" {
		return fmt.Errorf("postgres.database is required")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[strings.ToLower(cfg.Logging.Level)] {
		return fmt.Errorf("invalid logging.level: %s", cfg.Logging.Level)
	}

	validEnvs := map[string]bool{
		"development": true, "test": true, "production": true,
	}
	if !validEnvs[strings.ToLower(cfg.App.Env)] {
		return fmt.Errorf("invalid app.env: %s", cfg.App.Env)
	}

	return nil
}
