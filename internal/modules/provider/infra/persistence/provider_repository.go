package persistence

import (
	"lattice-coding/internal/common/db"
	"lattice-coding/internal/modules/provider/domain"

	"gorm.io/gorm"
)

// ProviderPO Provider 持久化对象
type ProviderPO struct {
	db.BasePO
	Name             string `gorm:"column:name;type:varchar(100);not null;index"`
	ProviderType     string `gorm:"column:provider_type;type:varchar(50);not null;index"`
	BaseURL          string `gorm:"column:base_url;type:varchar(255)"`
	APIKeyCiphertext string `gorm:"column:api_key_ciphertext;type:text"`
	IsEnabled        bool   `gorm:"column:is_enabled;default:true"`
}

// TableName 返回表名
func (ProviderPO) TableName() string {
	return "providers"
}

// ToDomain 转换为 Domain 实体
func (po *ProviderPO) ToDomain() *domain.Provider {
	return &domain.Provider{
		ID:               po.ID,
		Name:             po.Name,
		ProviderType:     domain.ProviderType(po.ProviderType),
		BaseURL:          po.BaseURL,
		APIKeyCiphertext: po.APIKeyCiphertext,
		IsEnabled:        po.IsEnabled,
		CreatedAt:        po.CreatedAt,
		UpdatedAt:        po.UpdatedAt,
	}
}

// FromDomain 从 Domain 实体转换
func (po *ProviderPO) FromDomain(d *domain.Provider) {
	po.ID = d.ID
	po.Name = d.Name
	po.ProviderType = string(d.ProviderType)
	po.BaseURL = d.BaseURL
	po.APIKeyCiphertext = d.APIKeyCiphertext
	po.IsEnabled = d.IsEnabled
}

// ProviderRepositoryImpl provider 仓库实现
type ProviderRepositoryImpl struct {
	db *gorm.DB
}

func NewProviderRepositoryImpl(db *gorm.DB) domain.ProviderRepository {
	return &ProviderRepositoryImpl{db: db}
}

func (r *ProviderRepositoryImpl) Create(provider *domain.Provider) error {
	po := &ProviderPO{}
	po.FromDomain(provider)
	return r.db.Create(po).Error
}

func (r *ProviderRepositoryImpl) Update(provider *domain.Provider) error {
	po := &ProviderPO{}
	po.FromDomain(provider)
	return r.db.Model(po).Omit("created_at").Updates(po).Error
}

func (r *ProviderRepositoryImpl) GetByID(id uint64) (*domain.Provider, error) {
	var po ProviderPO
	if err := r.db.First(&po, id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *ProviderRepositoryImpl) List() ([]*domain.Provider, error) {
	var pos []ProviderPO
	if err := r.db.Find(&pos).Error; err != nil {
		return nil, err
	}
	domains := make([]*domain.Provider, len(pos))
	for i, po := range pos {
		domains[i] = po.ToDomain()
	}
	return domains, nil
}

func (r *ProviderRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&ProviderPO{}, id).Error
}

// ModelConfigPO ModelConfig 持久化对象
type ModelConfigPO struct {
	db.BasePO
	ProviderID  uint64   `gorm:"column:provider_id;not null;index"`
	Name        string   `gorm:"column:name;type:varchar(100);not null"`
	ModelName   string   `gorm:"column:model_name;type:varchar(100);not null"`
	MaxTokens   *int     `gorm:"column:max_tokens"`
	Temperature *float64 `gorm:"column:temperature;type:decimal(3,2)"`
	TopP        *float64 `gorm:"column:top_p;type:decimal(3,2)"`
	ExtraConfig string   `gorm:"column:extra_config;type:text"`
	IsEnabled   bool     `gorm:"column:is_enabled;default:true"`
}

// TableName 返回表名
func (ModelConfigPO) TableName() string {
	return "model_configs"
}

// ToDomain 转换为 Domain 实体
func (po *ModelConfigPO) ToDomain() *domain.ModelConfig {
	return &domain.ModelConfig{
		ID:          po.ID,
		ProviderID:  po.ProviderID,
		Name:        po.Name,
		ModelName:   po.ModelName,
		MaxTokens:   po.MaxTokens,
		Temperature: po.Temperature,
		TopP:        po.TopP,
		ExtraConfig: po.ExtraConfig,
		IsEnabled:   po.IsEnabled,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// FromDomain 从 Domain 实体转换
func (po *ModelConfigPO) FromDomain(d *domain.ModelConfig) {
	po.ID = d.ID
	po.ProviderID = d.ProviderID
	po.Name = d.Name
	po.ModelName = d.ModelName
	po.MaxTokens = d.MaxTokens
	po.Temperature = d.Temperature
	po.TopP = d.TopP
	po.ExtraConfig = d.ExtraConfig
	po.IsEnabled = d.IsEnabled
}

// ModelConfigRepositoryImpl model_config 仓库实现
type ModelConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewModelConfigRepositoryImpl(db *gorm.DB) domain.ModelConfigRepository {
	return &ModelConfigRepositoryImpl{db: db}
}

func (r *ModelConfigRepositoryImpl) Create(modelConfig *domain.ModelConfig) error {
	po := &ModelConfigPO{}
	po.FromDomain(modelConfig)
	return r.db.Create(po).Error
}

func (r *ModelConfigRepositoryImpl) GetByID(id uint64) (*domain.ModelConfig, error) {
	var po ModelConfigPO
	if err := r.db.First(&po, id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *ModelConfigRepositoryImpl) List() ([]*domain.ModelConfig, error) {
	var pos []ModelConfigPO
	if err := r.db.Find(&pos).Error; err != nil {
		return nil, err
	}
	domains := make([]*domain.ModelConfig, len(pos))
	for i, po := range pos {
		domains[i] = po.ToDomain()
	}
	return domains, nil
}

func (r *ModelConfigRepositoryImpl) ListByProviderID(providerID uint64) ([]*domain.ModelConfig, error) {
	var pos []ModelConfigPO
	if err := r.db.Where("provider_id = ?", providerID).Find(&pos).Error; err != nil {
		return nil, err
	}
	domains := make([]*domain.ModelConfig, len(pos))
	for i, po := range pos {
		domains[i] = po.ToDomain()
	}
	return domains, nil
}

func (r *ModelConfigRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&ModelConfigPO{}, id).Error
}

// Migrate 自动迁移表结构
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&ProviderPO{}); err != nil {
		return err
	}
	return db.AutoMigrate(&ModelConfigPO{})
}
