package persistence

import "lattice-coding/internal/common/db"

type KnowledgePackPO struct {
	db.BasePO
	PackKey       string `gorm:"column:pack_key;type:varchar(120);not null;uniqueIndex"`
	Query         string `gorm:"column:query;type:text"`
	Intent        string `gorm:"column:intent;type:varchar(80);not null;default:'';index"`
	Route         string `gorm:"column:route;type:varchar(80);not null;default:'';index"`
	Status        string `gorm:"column:status;type:varchar(20);not null;default:'draft';index"`
	TokenEstimate int    `gorm:"column:token_estimate;not null;default:0"`
	MaxTokens     int    `gorm:"column:max_tokens;not null;default:0"`
	PromptContext string `gorm:"column:prompt_context;type:longtext"`
	Warnings      string `gorm:"column:warnings;type:json"`
	Options       string `gorm:"column:options;type:json"`
	Meta          string `gorm:"column:meta;type:json"`
}

func (KnowledgePackPO) TableName() string {
	return "knowledge_pack"
}

type KnowledgeItemPO struct {
	db.BasePO
	PackID        uint64  `gorm:"column:pack_id;not null;index"`
	ItemKey       string  `gorm:"column:item_key;type:varchar(120);not null;index"`
	SourceKind    string  `gorm:"column:source_kind;type:varchar(50);not null;index"`
	SourceID      string  `gorm:"column:source_id;type:varchar(200);not null;default:'';index"`
	SourceType    string  `gorm:"column:source_type;type:varchar(80);not null;default:'';index"`
	Title         string  `gorm:"column:title;type:varchar(500);not null;default:''"`
	Content       string  `gorm:"column:content;type:longtext"`
	Location      string  `gorm:"column:location;type:varchar(500);not null;default:''"`
	Score         float64 `gorm:"column:score;not null;default:0"`
	TokenEstimate int     `gorm:"column:token_estimate;not null;default:0"`
	CitationKey   string  `gorm:"column:citation_key;type:varchar(120);not null;default:'';index"`
	Metadata      string  `gorm:"column:metadata;type:json"`
	SortOrder     int     `gorm:"column:sort_order;not null;default:0;index"`
}

func (KnowledgeItemPO) TableName() string {
	return "knowledge_item"
}

type KnowledgeCitationPO struct {
	db.BasePO
	PackID      uint64  `gorm:"column:pack_id;not null;index"`
	CitationKey string  `gorm:"column:citation_key;type:varchar(120);not null;index"`
	SourceKind  string  `gorm:"column:source_kind;type:varchar(50);not null;index"`
	SourceID    string  `gorm:"column:source_id;type:varchar(200);not null;default:'';index"`
	Title       string  `gorm:"column:title;type:varchar(500);not null;default:''"`
	Location    string  `gorm:"column:location;type:varchar(500);not null;default:''"`
	URI         string  `gorm:"column:uri;type:varchar(1000);not null;default:''"`
	Score       float64 `gorm:"column:score;not null;default:0"`
	Metadata    string  `gorm:"column:metadata;type:json"`
	SortOrder   int     `gorm:"column:sort_order;not null;default:0;index"`
}

func (KnowledgeCitationPO) TableName() string {
	return "knowledge_citation"
}
