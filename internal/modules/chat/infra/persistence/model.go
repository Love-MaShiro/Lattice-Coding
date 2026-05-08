package persistence

import "lattice-coding/internal/common/db"

type ChatSessionPO struct {
	db.BasePO
	Title                    string `gorm:"column:title;type:varchar(200);not null;default:''"`
	AgentID                  uint64 `gorm:"column:agent_id;not null;index"`
	ModelConfigID            uint64 `gorm:"column:model_config_id;not null;index"`
	Status                   string `gorm:"column:status;type:varchar(20);not null;default:'active';index"`
	Summary                  string `gorm:"column:summary;type:longtext"`
	SummarizedUntilMessageID uint64 `gorm:"column:summarized_until_message_id;not null;default:0;index"`
	Meta                     string `gorm:"column:meta;type:json"`
}

func (ChatSessionPO) TableName() string {
	return "chat_session"
}

type ChatMessagePO struct {
	db.BasePO
	SessionID  uint64 `gorm:"column:session_id;not null;index"`
	Role       string `gorm:"column:role;type:varchar(20);not null;index"`
	Content    string `gorm:"column:content;type:longtext;not null"`
	TokenCount int    `gorm:"column:token_count;not null;default:0"`
	Meta       string `gorm:"column:meta;type:json"`
}

func (ChatMessagePO) TableName() string {
	return "chat_message"
}
