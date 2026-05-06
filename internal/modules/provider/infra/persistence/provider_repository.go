package persistence

import (
	"context"
	"time"

	"lattice-coding/internal/common/db"
	"lattice-coding/internal/modules/provider/domain"

	"gorm.io/gorm"
)

type ProviderPO struct {
	db.BasePO
	Name                 string     `gorm:"column:name;type:varchar(100);not null;index"`
	ProviderType         string     `gorm:"column:provider_type;type:varchar(50);not null;index"`
	BaseURL              string     `gorm:"column:base_url;type:varchar(255)"`
	AuthType             string     `gorm:"column:auth_type;type:varchar(50);not null;default:'api_key'"`
	APIKeyCiphertext     string     `gorm:"column:api_key_ciphertext;type:text"`
	AuthConfigCiphertext string     `gorm:"column:auth_config_ciphertext;type:text"`
	Config               string     `gorm:"column:config;type:text"`
	Enabled              bool       `gorm:"column:enabled;default:true"`
	HealthStatus         string     `gorm:"column:health_status;type:varchar(20);not null;default:'unknown'"`
	LastCheckedAt        *time.Time `gorm:"column:last_checked_at"`
	LastError            string     `gorm:"column:last_error;type:text"`
}

func (ProviderPO) TableName() string {
	return "providers"
}

type ModelConfigPO struct {
	db.BasePO
	ProviderID   uint64 `gorm:"column:provider_id;not null;index"`
	Name         string `gorm:"column:name;type:varchar(100);not null"`
	Model        string `gorm:"column:model;type:varchar(100);not null"`
	ModelType    string `gorm:"column:model_type;type:varchar(50);not null;default:'chat'"`
	Params       string `gorm:"column:params;type:text"`
	Capabilities string `gorm:"column:capabilities;type:text"`
	IsDefault    bool   `gorm:"column:is_default;default:false"`
	Enabled      bool   `gorm:"column:enabled;default:true"`
}

func (ModelConfigPO) TableName() string {
	return "model_configs"
}

type ProviderHealthPO struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ProviderID    uint64    `gorm:"column:provider_id;not null;index"`
	ModelConfigID uint64    `gorm:"column:model_config_id;not null;index"`
	Status        string    `gorm:"column:status;type:varchar(20);not null;default:'unknown'"`
	LatencyMs     int64     `gorm:"column:latency_ms;default:0"`
	ErrorCode     string    `gorm:"column:error_code;type:varchar(50)"`
	ErrorMessage  string    `gorm:"column:error_message;type:text"`
	CheckedAt     time.Time `gorm:"column:checked_at;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (ProviderHealthPO) TableName() string {
	return "provider_healths"
}

type ProviderRepositoryImpl struct {
	db *gorm.DB
}

func NewProviderRepositoryImpl(db *gorm.DB) domain.ProviderRepository {
	return &ProviderRepositoryImpl{db: db}
}

func (r *ProviderRepositoryImpl) Create(ctx context.Context, provider *domain.Provider) error {
	po := &ProviderPO{}
	ConvertProviderToPO(provider, po)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *ProviderRepositoryImpl) Update(ctx context.Context, provider *domain.Provider) error {
	po := &ProviderPO{}
	ConvertProviderToPO(provider, po)
	return r.db.WithContext(ctx).Model(po).Omit("created_at").Updates(po).Error
}

func (r *ProviderRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.Provider, error) {
	var po ProviderPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, err
	}
	return ConvertPOToProvider(&po), nil
}

func (r *ProviderRepositoryImpl) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.Provider], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ProviderPO{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	var pos []ProviderPO
	if err := r.db.WithContext(ctx).Offset(offset).Limit(req.PageSize).Find(&pos).Error; err != nil {
		return nil, err
	}

	items := make([]*domain.Provider, len(pos))
	for i := range pos {
		items[i] = ConvertPOToProvider(&pos[i])
	}

	return &domain.PageResult[*domain.Provider]{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (r *ProviderRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&ProviderPO{}, id).Error
}

func (r *ProviderRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ProviderPO{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProviderRepositoryImpl) UpdateEnabled(ctx context.Context, id uint64, enabled bool) error {
	return r.db.WithContext(ctx).Model(&ProviderPO{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func (r *ProviderRepositoryImpl) UpdateHealthStatus(ctx context.Context, id uint64, status domain.HealthStatus, lastError string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&ProviderPO{}).Where("id = ?", id).Updates(map[string]interface{}{
		"health_status":   status,
		"last_error":      lastError,
		"last_checked_at": &now,
	}).Error
}

type ModelConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewModelConfigRepositoryImpl(db *gorm.DB) domain.ModelConfigRepository {
	return &ModelConfigRepositoryImpl{db: db}
}

func (r *ModelConfigRepositoryImpl) Create(ctx context.Context, modelConfig *domain.ModelConfig) error {
	po := &ModelConfigPO{}
	ConvertModelConfigToPO(modelConfig, po)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *ModelConfigRepositoryImpl) Update(ctx context.Context, modelConfig *domain.ModelConfig) error {
	po := &ModelConfigPO{}
	ConvertModelConfigToPO(modelConfig, po)
	return r.db.WithContext(ctx).Model(po).Omit("created_at").Updates(po).Error
}

func (r *ModelConfigRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.ModelConfig, error) {
	var po ModelConfigPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, err
	}
	return ConvertPOToModelConfig(&po), nil
}

func (r *ModelConfigRepositoryImpl) FindByProviderID(ctx context.Context, providerID uint64) ([]*domain.ModelConfig, error) {
	var pos []ModelConfigPO
	if err := r.db.WithContext(ctx).Where("provider_id = ?", providerID).Find(&pos).Error; err != nil {
		return nil, err
	}
	items := make([]*domain.ModelConfig, len(pos))
	for i := range pos {
		items[i] = ConvertPOToModelConfig(&pos[i])
	}
	return items, nil
}

func (r *ModelConfigRepositoryImpl) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.ModelConfig], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ModelConfigPO{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	var pos []ModelConfigPO
	if err := r.db.WithContext(ctx).Offset(offset).Limit(req.PageSize).Find(&pos).Error; err != nil {
		return nil, err
	}

	items := make([]*domain.ModelConfig, len(pos))
	for i := range pos {
		items[i] = ConvertPOToModelConfig(&pos[i])
	}

	return &domain.PageResult[*domain.ModelConfig]{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (r *ModelConfigRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&ModelConfigPO{}, id).Error
}

func (r *ModelConfigRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ModelConfigPO{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ModelConfigRepositoryImpl) UpdateEnabled(ctx context.Context, id uint64, enabled bool) error {
	return r.db.WithContext(ctx).Model(&ModelConfigPO{}).Where("id = ?", id).Update("enabled", enabled).Error
}

type ProviderHealthRepositoryImpl struct {
	db *gorm.DB
}

func NewProviderHealthRepositoryImpl(db *gorm.DB) domain.ProviderHealthRepository {
	return &ProviderHealthRepositoryImpl{db: db}
}

func (r *ProviderHealthRepositoryImpl) Create(ctx context.Context, health *domain.ProviderHealth) error {
	po := &ProviderHealthPO{}
	ConvertProviderHealthToPO(health, po)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *ProviderHealthRepositoryImpl) FindLatestByProviderID(ctx context.Context, providerID uint64) (*domain.ProviderHealth, error) {
	var po ProviderHealthPO
	if err := r.db.WithContext(ctx).Where("provider_id = ?", providerID).Order("checked_at DESC").First(&po).Error; err != nil {
		return nil, err
	}
	return ConvertPOToProviderHealth(&po), nil
}

func (r *ProviderHealthRepositoryImpl) FindLatestByModelConfigID(ctx context.Context, modelConfigID uint64) (*domain.ProviderHealth, error) {
	var po ProviderHealthPO
	if err := r.db.WithContext(ctx).Where("model_config_id = ?", modelConfigID).Order("checked_at DESC").First(&po).Error; err != nil {
		return nil, err
	}
	return ConvertPOToProviderHealth(&po), nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&ProviderPO{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ModelConfigPO{}); err != nil {
		return err
	}
	return db.AutoMigrate(&ProviderHealthPO{})
}
