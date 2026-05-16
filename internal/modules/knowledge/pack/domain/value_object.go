package domain

type PackStatus string

const (
	PackStatusDraft PackStatus = "draft"
	PackStatusReady PackStatus = "ready"
	PackStatusStale PackStatus = "stale"
)

type PackSourceKind string

const (
	PackSourceRAG                PackSourceKind = "rag"
	PackSourceLocalFile          PackSourceKind = "local_file"
	PackSourceCodeSymbol         PackSourceKind = "code_symbol"
	PackSourceToolResult         PackSourceKind = "tool_result"
	PackSourceWebPage            PackSourceKind = "web_page"
	PackSourceConversationMemory PackSourceKind = "conversation_memory"
	PackSourceProjectSnapshot    PackSourceKind = "project_snapshot"
)
