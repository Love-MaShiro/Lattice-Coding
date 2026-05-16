package persistence

import "time"

type RunPO struct {
	ID          string     `gorm:"column:id;type:varchar(64);primaryKey"`
	AgentID     string     `gorm:"column:agent_id;type:varchar(64);index"`
	SessionID   string     `gorm:"column:session_id;type:varchar(64);index"`
	WorkflowID  string     `gorm:"column:workflow_id;type:varchar(64);index"`
	Status      string     `gorm:"column:status;type:varchar(20);not null;index"`
	Input       string     `gorm:"column:input;type:longtext"`
	Output      string     `gorm:"column:output;type:longtext"`
	Error       string     `gorm:"column:error;type:longtext"`
	StartedAt   time.Time  `gorm:"column:started_at;index"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (RunPO) TableName() string {
	return "runs"
}

type ToolInvocationPO struct {
	ID            string     `gorm:"column:id;type:varchar(64);primaryKey"`
	RunID         string     `gorm:"column:run_id;type:varchar(64);index"`
	NodeID        string     `gorm:"column:node_id;type:varchar(120);index"`
	ToolName      string     `gorm:"column:tool_name;type:varchar(120);not null;index"`
	InputJSON     string     `gorm:"column:input_json;type:longtext"`
	OutputJSON    string     `gorm:"column:output_json;type:longtext"`
	IsError       bool       `gorm:"column:is_error;not null;default:false;index"`
	LatencyMs     int64      `gorm:"column:latency_ms;not null;default:0"`
	Status        string     `gorm:"column:status;type:varchar(20);not null;index"`
	FullResultRef string     `gorm:"column:full_result_ref;type:varchar(500)"`
	StartedAt     time.Time  `gorm:"column:started_at;index"`
	CompletedAt   *time.Time `gorm:"column:completed_at"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (ToolInvocationPO) TableName() string {
	return "tool_invocations"
}
