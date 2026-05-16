package domain

type ContextSourceKind string

const (
	ContextSourceRAG             ContextSourceKind = "rag"
	ContextSourceProjectFile     ContextSourceKind = "project_file"
	ContextSourceCodeSymbol      ContextSourceKind = "code_symbol"
	ContextSourceToolResult      ContextSourceKind = "tool_result"
	ContextSourceWebPage         ContextSourceKind = "web_page"
	ContextSourceProjectSnapshot ContextSourceKind = "project_snapshot"
)

type ContextCandidateStatus string

const (
	ContextCandidatePending  ContextCandidateStatus = "pending"
	ContextCandidateSelected ContextCandidateStatus = "selected"
	ContextCandidateRejected ContextCandidateStatus = "rejected"
)

type ContextSignalKind string

const (
	ContextSignalExplicitMention ContextSignalKind = "explicit_mention"
	ContextSignalPathMatch       ContextSignalKind = "path_match"
	ContextSignalSymbolMatch     ContextSignalKind = "symbol_match"
	ContextSignalRecentEdit      ContextSignalKind = "recent_edit"
	ContextSignalRetrievalScore  ContextSignalKind = "retrieval_score"
	ContextSignalToolOutput      ContextSignalKind = "tool_output"
)
