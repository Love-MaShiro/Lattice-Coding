package persistence

import "lattice-coding/internal/common/db"

type ContextSourcePO struct {
	db.BasePO
	SourceKey string `gorm:"column:source_key;type:varchar(120);not null;uniqueIndex"`
	Kind      string `gorm:"column:kind;type:varchar(50);not null;index"`
	Name      string `gorm:"column:name;type:varchar(200);not null;default:'';index"`
	URI       string `gorm:"column:uri;type:varchar(1000);not null;default:''"`
	Scope     string `gorm:"column:scope;type:varchar(200);not null;default:'';index"`
	Metadata  string `gorm:"column:metadata;type:json"`
}

func (ContextSourcePO) TableName() string {
	return "context_source"
}

type ContextCandidatePO struct {
	db.BasePO
	CandidateKey  string  `gorm:"column:candidate_key;type:varchar(120);not null;uniqueIndex"`
	SourceKey     string  `gorm:"column:source_key;type:varchar(120);not null;index"`
	SourceKind    string  `gorm:"column:source_kind;type:varchar(50);not null;index"`
	Title         string  `gorm:"column:title;type:varchar(500);not null;default:''"`
	Content       string  `gorm:"column:content;type:longtext"`
	Location      string  `gorm:"column:location;type:varchar(500);not null;default:''"`
	Score         float64 `gorm:"column:score;not null;default:0;index"`
	TokenEstimate int     `gorm:"column:token_estimate;not null;default:0"`
	Status        string  `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	Metadata      string  `gorm:"column:metadata;type:json"`
}

func (ContextCandidatePO) TableName() string {
	return "context_candidate"
}

type ContextSignalPO struct {
	db.BasePO
	CandidateID uint64  `gorm:"column:candidate_id;not null;index"`
	SignalKey   string  `gorm:"column:signal_key;type:varchar(120);not null;index"`
	Kind        string  `gorm:"column:kind;type:varchar(50);not null;index"`
	Weight      float64 `gorm:"column:weight;not null;default:0"`
	Reason      string  `gorm:"column:reason;type:text"`
	Metadata    string  `gorm:"column:metadata;type:json"`
}

func (ContextSignalPO) TableName() string {
	return "context_signal"
}

type ContextPolicyPO struct {
	db.BasePO
	PolicyKey   string `gorm:"column:policy_key;type:varchar(120);not null;uniqueIndex"`
	Name        string `gorm:"column:name;type:varchar(200);not null;default:''"`
	Description string `gorm:"column:description;type:text"`
	MaxTokens   int    `gorm:"column:max_tokens;not null;default:0"`
	MaxItems    int    `gorm:"column:max_items;not null;default:0"`
	Rules       string `gorm:"column:rules;type:json"`
	Metadata    string `gorm:"column:metadata;type:json"`
}

func (ContextPolicyPO) TableName() string {
	return "context_policy"
}
