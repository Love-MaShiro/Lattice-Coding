package query

type QueryResult struct {
	RunID    string
	Mode     ExecutionMode
	Content  string
	Messages []Message
	Steps    []StepResult
	Usage    Usage
	Metadata map[string]interface{}
}

type Result = QueryResult

type Message struct {
	Role    string
	Content string
}

type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}
