package app

import (
	"gorm.io/gorm"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/db"
	"lattice-coding/internal/common/redis"
)

type Dependencies struct {
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := db.NewMySQL(&cfg.MySQL)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,
	}, nil
}

func (d *Dependencies) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if d.Redis != nil {
		d.Redis.Close()
	}

	return nil
}
