package persistence

import (
	"lattice-coding/internal/common/db"
)

type AgentPO struct {
	db.BasePO
	Name            string  `gorm:"column:name;type:varchar(100);not null;index"`
	Description     string  `gorm:"column:description;type:text"`
	AgentType       string  `gorm:"column:agent_type;type:varchar(50);not null;default:'customer_service';index"`
	ModelConfigID   uint64  `gorm:"column:model_config_id;not null;index"`
	SystemPrompt    string  `gorm:"column:system_prompt;type:text"`
	Temperature     float64 `gorm:"column:temperature;type:decimal(3,2);not null;default:0.70"`
	TopP            float64 `gorm:"column:top_p;type:decimal(3,2);not null;default:1.00"`
	MaxTokens       int     `gorm:"column:max_tokens;not null;default:4096"`
	MaxContextTurns int     `gorm:"column:max_context_turns;not null;default:10"`
	MaxSteps        int     `gorm:"column:max_steps;not null;default:20"`
	Enabled         bool    `gorm:"column:enabled;not null;default:true;index"`
}

func (AgentPO) TableName() string {
	return "agent"
}

type AgentToolPO struct {
	db.BasePO
	AgentID  uint64 `gorm:"column:agent_id;not null;index"`
	ToolID   uint64 `gorm:"column:tool_id;not null;index"`
	ToolType string `gorm:"column:tool_type;type:varchar(50);not null;default:''"`
}

func (AgentToolPO) TableName() string {
	return "agent_tool"
}
