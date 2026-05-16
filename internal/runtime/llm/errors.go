package llm

import "errors"

var (
	ErrPoolFull        = NewBizError("LLM_001", "llm worker pool is full")
	ErrTimeout         = NewBizError("LLM_002", "llm request timeout")
	ErrNoProvider      = NewBizError("LLM_003", "no llm provider available")
	ErrAllFallbackFail = NewBizError("LLM_004", "all llm providers failed")
)

var (
	ErrUnsupportedProviderType = errors.New("unsupported provider type")
	ErrProviderDisabled        = errors.New("provider is disabled")
	ErrModelConfigDisabled     = errors.New("model config is disabled")
	ErrNoModelConfigFound      = errors.New("no model config found")
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
