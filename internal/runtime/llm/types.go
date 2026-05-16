package llm

import "time"

type ChatRequest struct {
	Provider      string
	Model         string
	ModelConfigID uint64
	Messages      []Message
	Temperature   *float64
	TopP          *float64
	MaxTokens     int
}

type Message struct {
	Role    string
	Content string
}

type ChatResponse struct {
	Content   string
	Usage     *Usage
	LatencyMs int64
}

type StreamChunk struct {
	Content   string
	Done      bool
	Usage     *Usage
	LatencyMs int64
	Err       error
}

type CallResult struct {
	Provider  string
	Model     string
	LatencyMs int64
	Tokens    int
	Error     error
	Success   bool
	Fallback  bool
}

type CallOptions struct {
	Timeout     time.Duration
	Retry       bool
	MaxAttempts int
}
