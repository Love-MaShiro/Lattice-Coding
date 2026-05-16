package query

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithMessage(message string) *Error {
	return &Error{Code: e.Code, Message: message}
}

var (
	ErrInvalidRequest     = &Error{Code: "QUERY_INVALID_REQUEST", Message: "invalid query request"}
	ErrModeNotSupported   = &Error{Code: "QUERY_MODE_NOT_SUPPORTED", Message: "query execution mode not supported"}
	ErrStrategyNotFound   = &Error{Code: "QUERY_STRATEGY_NOT_FOUND", Message: "query strategy not found"}
	ErrBudgetExceeded     = &Error{Code: "QUERY_BUDGET_EXCEEDED", Message: "query budget exceeded"}
	ErrQueryTimeout       = &Error{Code: "QUERY_TIMEOUT", Message: "query timeout"}
	ErrQueryInterrupted   = &Error{Code: "QUERY_INTERRUPTED", Message: "query interrupted"}
	ErrQueryFailed        = &Error{Code: "QUERY_FAILED", Message: "query failed"}
	ErrRunBindFailed      = &Error{Code: "QUERY_RUN_BIND_FAILED", Message: "query run bind failed"}
	ErrContextBuildFailed = &Error{Code: "QUERY_CONTEXT_BUILD_FAILED", Message: "query context build failed"}
)
