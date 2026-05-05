package app

import (
	"gorm.io/gorm"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/db"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/common/redis"
	"lattice-coding/internal/runtime/llm"
)

type Dependencies struct {
	Config      *config.Config
	Logger      *logger.Logger
	MySQL       *gorm.DB
	Redis       *redis.Client
	VectorDB    *gorm.DB
	LLMExecutor *llm.Executor
}

func NewDependencies(cfg *config.Config, log *logger.Logger) (*Dependencies, error) {
	mysqlDB, err := db.NewMySQL(&cfg.MySQL)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		return nil, err
	}

	vectorDB, err := db.NewPostgres(&cfg.Postgres)
	if err != nil {
		return nil, err
	}

	llmExecutor := newLLMExecutor(cfg)

	return &Dependencies{
		Config:      cfg,
		Logger:      log,
		MySQL:       mysqlDB,
		Redis:       redisClient,
		VectorDB:    vectorDB,
		LLMExecutor: llmExecutor,
	}, nil
}

func (d *Dependencies) Close() error {
	if d.MySQL != nil {
		sqlDB, err := d.MySQL.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	if d.Redis != nil {
		d.Redis.Close()
	}

	if d.VectorDB != nil {
		sqlDB, err := d.VectorDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	if d.Logger != nil {
		d.Logger.Sync()
	}

	return nil
}

func newLLMExecutor(cfg *config.Config) *llm.Executor {
	return llm.NewExecutor(&cfg.LLM)
}
