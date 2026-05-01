package app

import (
	"gorm.io/gorm"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/db"
	"lattice-coding/internal/common/redis"
)

type Bootstrap struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewBootstrap(cfg *config.Config) (*Bootstrap, error) {
	db, err := db.NewMySQL(&cfg.MySQL)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Bootstrap{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func (b *Bootstrap) Close() error {
	if b.DB != nil {
		sqlDB, err := b.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if b.Redis != nil {
		b.Redis.Close()
	}

	return nil
}
