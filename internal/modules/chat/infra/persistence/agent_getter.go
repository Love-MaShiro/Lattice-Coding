package persistence

import (
	"context"
	"errors"

	"lattice-coding/internal/modules/chat/application"

	"gorm.io/gorm"
)

type AgentRuntimePO struct {
	ID              uint64  `gorm:"column:id"`
	Name            string  `gorm:"column:name"`
	ModelConfigID   uint64  `gorm:"column:model_config_id"`
	SystemPrompt    string  `gorm:"column:system_prompt"`
	Temperature     float64 `gorm:"column:temperature"`
	TopP            float64 `gorm:"column:top_p"`
	MaxTokens       int     `gorm:"column:max_tokens"`
	MaxContextTurns int     `gorm:"column:max_context_turns"`
	Enabled         bool    `gorm:"column:enabled"`
}

func (AgentRuntimePO) TableName() string {
	return "agent"
}

type AgentGetter struct {
	db *gorm.DB
}

func NewAgentGetter(db *gorm.DB) application.AgentGetter {
	return &AgentGetter{db: db}
}

func (g *AgentGetter) GetAgentForChat(ctx context.Context, id uint64) (*application.AgentRuntimeDTO, error) {
	var po AgentRuntimePO
	if err := g.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &application.AgentRuntimeDTO{
		ID:              po.ID,
		Name:            po.Name,
		ModelConfigID:   po.ModelConfigID,
		SystemPrompt:    po.SystemPrompt,
		Temperature:     po.Temperature,
		TopP:            po.TopP,
		MaxTokens:       po.MaxTokens,
		MaxContextTurns: po.MaxContextTurns,
		Enabled:         po.Enabled,
	}, nil
}
