package llm

import "time"

type ChatRequest struct {
	Provider string
	Model    string
	Messages []Message
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

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
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

var (
	ErrPoolFull        = NewBizError("LLM_001", "系统繁忙，请稍后重试")
	ErrTimeout         = NewBizError("LLM_002", "模型响应超时")
	ErrNoProvider      = NewBizError("LLM_003", "无可用的模型提供商")
	ErrAllFallbackFail = NewBizError("LLM_004", "所有模型提供商均不可用")
)

type BizError struct {
	Code    string
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

func NewBizError(code, message string) *BizError {
	return &BizError{Code: code, Message: message}
}

func (e *BizError) WithMessage(msg string) *BizError {
	return &BizError{Code: e.Code, Message: msg}
}
