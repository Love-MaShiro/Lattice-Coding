package config

import (
	"os"

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
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker"`
	Routing        RoutingConfig        `yaml:"routing"`
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
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type HTTPConfig struct {
	Timeout      string `yaml:"timeout"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

type WorkerConfig struct {
	Concurrency int    `yaml:"concurrency"`
	QueueName   string `yaml:"queue_name"`
}

func LoadConfig() *Config {
	path := "configs/config.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return &cfg
}
