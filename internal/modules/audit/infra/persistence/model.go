package persistence

import "lattice-coding/internal/common/db"

type AuditLogPO struct {
	db.BasePO
	RunID        string `gorm:"column:run_id;type:varchar(64);index"`
	TraceID      string `gorm:"column:trace_id;type:varchar(80);index"`
	EventType    string `gorm:"column:event_type;type:varchar(80);not null;index"`
	ToolName     string `gorm:"column:tool_name;type:varchar(120);index"`
	ResourceType string `gorm:"column:resource_type;type:varchar(80);index"`
	ResourceID   string `gorm:"column:resource_id;type:varchar(120);index"`
	Message      string `gorm:"column:message;type:varchar(1000)"`
	PayloadJSON  string `gorm:"column:payload_json;type:longtext"`
}

func (AuditLogPO) TableName() string {
	return "audit_logs"
}
