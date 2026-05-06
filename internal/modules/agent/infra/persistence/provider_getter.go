package persistence

import (
	"context"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/agent/application"

	"gorm.io/gorm"
)

type ProviderGetter struct {
	db *gorm.DB
}

func NewProviderGetter(db *gorm.DB) application.ModelConfigGetter {
	return &ProviderGetter{db: db}
}

func (g *ProviderGetter) GetModelConfig(ctx context.Context, id uint64) (uint64, error) {
	var po struct {
		ProviderID uint64 `gorm:"column:provider_id"`
	}
	if err := g.db.WithContext(ctx).
		Model(&providerModelConfigPO{}).
		Select("provider_id").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.NotFoundErr("ModelConfig 不存在")
		}
		return 0, errors.DatabaseErrWithErr(err, "查询 ModelConfig 失败")
	}
	return po.ProviderID, nil
}

func (g *ProviderGetter) GetProvider(ctx context.Context, id uint64) error {
	var count int64
	if err := g.db.WithContext(ctx).
		Model(&providerPO{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error; err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if count == 0 {
		return errors.NotFoundErr("Provider 不存在")
	}
	return nil
}

type providerPO struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (providerPO) TableName() string {
	return "providers"
}

type providerModelConfigPO struct {
	ProviderID uint64 `gorm:"column:provider_id"`
}

func (providerModelConfigPO) TableName() string {
	return "model_configs"
}
