package persistence

import "lattice-coding/internal/common/db"

type WorkflowPO struct {
	db.BasePO
	Name        string `gorm:"column:name;type:varchar(200);not null;index"`
	Description string `gorm:"column:description;type:text"`
	Status      string `gorm:"column:status;type:varchar(20);not null;default:'draft';index"`
	Version     int    `gorm:"column:version;not null;default:1"`
	Meta        string `gorm:"column:meta;type:json"`
}

func (WorkflowPO) TableName() string {
	return "workflow"
}

type WorkflowNodePO struct {
	db.BasePO
	WorkflowID uint64 `gorm:"column:workflow_id;not null;index"`
	NodeKey    string `gorm:"column:node_key;type:varchar(120);not null;index"`
	Name       string `gorm:"column:name;type:varchar(200);not null;default:''"`
	Type       string `gorm:"column:type;type:varchar(50);not null;index"`
	Config     string `gorm:"column:config;type:json"`
	Position   string `gorm:"column:position;type:json"`
	SortOrder  int    `gorm:"column:sort_order;not null;default:0;index"`
	Meta       string `gorm:"column:meta;type:json"`
}

func (WorkflowNodePO) TableName() string {
	return "workflow_node"
}

type WorkflowEdgePO struct {
	db.BasePO
	WorkflowID uint64 `gorm:"column:workflow_id;not null;index"`
	EdgeKey    string `gorm:"column:edge_key;type:varchar(160);not null;index"`
	SourceKey  string `gorm:"column:source_key;type:varchar(120);not null;index"`
	TargetKey  string `gorm:"column:target_key;type:varchar(120);not null;index"`
	Condition  string `gorm:"column:condition_expr;type:text"`
	SortOrder  int    `gorm:"column:sort_order;not null;default:0;index"`
	Meta       string `gorm:"column:meta;type:json"`
}

func (WorkflowEdgePO) TableName() string {
	return "workflow_edge"
}
