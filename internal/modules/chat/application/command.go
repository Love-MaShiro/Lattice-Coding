package application

type CreateSessionCommand struct {
	Title   string `json:"title"`
	AgentID uint64 `json:"agent_id"`
}

type UpdateSessionSummaryCommand struct {
	SessionID                uint64 `json:"session_id"`
	Summary                  string `json:"summary"`
	SummarizedUntilMessageID uint64 `json:"summarized_until_message_id"`
}

type CreateMessageCommand struct {
	SessionID uint64 `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Meta      string `json:"meta"`
}

type CompletionCommand struct {
	AgentID   uint64 `json:"agent_id"`
	SessionID uint64 `json:"session_id"`
	Message   string `json:"message"`
}
