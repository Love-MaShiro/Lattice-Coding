package domain

type WorkflowStatus string

const (
	WorkflowStatusDraft    WorkflowStatus = "draft"
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusArchived WorkflowStatus = "archived"
)

type NodeType string

const (
	NodeTypeStart             NodeType = "start"
	NodeTypeEnd               NodeType = "end"
	NodeTypeLLM               NodeType = "llm"
	NodeTypeTool              NodeType = "tool"
	NodeTypeCondition         NodeType = "condition"
	NodeTypeParallel          NodeType = "parallel"
	NodeTypeKnowledgeRoute    NodeType = "knowledge_route"
	NodeTypeKnowledgeRetrieve NodeType = "knowledge_retrieve"
	NodeTypeContextBuild      NodeType = "context_build"
	NodeTypeContextCompress   NodeType = "context_compress"
	NodeTypeWebSearch         NodeType = "web_search"
	NodeTypeMCPCall           NodeType = "mcp_call"
	NodeTypeCodeSearch        NodeType = "code_search"
	NodeTypeFileRead          NodeType = "file_read"
	NodeTypeShellCommand      NodeType = "shell_command"
)

type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusRunning   RunStatus = "running"
	RunStatusSucceeded RunStatus = "succeeded"
	RunStatusFailed    RunStatus = "failed"
	RunStatusCancelled RunStatus = "cancelled"
)
