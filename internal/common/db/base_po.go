package db

import (
	"time"

	"gorm.io/gorm"
)

type BasePO struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
