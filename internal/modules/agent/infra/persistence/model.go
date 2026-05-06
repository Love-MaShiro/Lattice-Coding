package persistence

import (
	"lattice-coding/internal/common/db"
)

type AgentPO struct {
	db.BasePO
	Name          string `gorm:"column:name;type:varchar(100);not null;index"`
	Description   string `gorm:"column:description;type:text"`
	ProviderID    uint64 `gorm:"column:provider_id;not null;index"`
	ModelConfigID uint64 `gorm:"column:model_config_id;not null;index"`
	SystemPrompt  string `gorm:"column:system_prompt;type:text"`
	Tools         string `gorm:"column:tools;type:json"`
	MaxSteps      int    `gorm:"column:max_steps;not null;default:20"`
	Timeout       int    `gorm:"column:timeout;not null;default:1200"`
	Enabled       bool   `gorm:"column:enabled;not null;default:true;index"`
}

func (AgentPO) TableName() string {
	return "agents"
}
