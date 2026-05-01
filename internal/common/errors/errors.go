package errors

import "fmt"

type ErrorCode string

const (
	// ============ 通用错误 (10000-10099) ============
	Success            ErrorCode = "10000"
	InternalError      ErrorCode = "10001"
	InvalidArgument    ErrorCode = "10002"
	NotFound           ErrorCode = "10003"
	Unauthorized       ErrorCode = "10004"
	Forbidden          ErrorCode = "10005"
	AlreadyExists      ErrorCode = "10006"
	ValidationError    ErrorCode = "10007"
	ServiceUnavailable ErrorCode = "10008"
	Timeout            ErrorCode = "10009"
	DatabaseError      ErrorCode = "10010"
	CacheError         ErrorCode = "10011"
	NetworkError       ErrorCode = "10012"

	// ============ Agent Run 错误 (20000-20099) ============
	RunNotFound        ErrorCode = "20000"
	RunNotStarted      ErrorCode = "20001"
	RunAlreadyStarted  ErrorCode = "20002"
	RunAlreadyFinished ErrorCode = "20003"
	RunCancelled       ErrorCode = "20004"
	RunFailed          ErrorCode = "20005"
	RunTimeout         ErrorCode = "20006"
	RunMaxRetries      ErrorCode = "20007"
	RunInvalidStatus   ErrorCode = "20008"
	RunMissingAgent    ErrorCode = "20009"

	// ============ 模型调用错误 (30000-30099) ============
	LLMError               ErrorCode = "30000"
	LLMServiceError        ErrorCode = "30001"
	LLMAuthenticationError ErrorCode = "30002"
	LLMRateLimit           ErrorCode = "30003"
	LLMQuotaExceeded       ErrorCode = "30004"
	LLMInvalidRequest      ErrorCode = "30005"
	LLMModelNotFound       ErrorCode = "30006"
	LLMContextOverlimit    ErrorCode = "30007"
	LLMGenerationError     ErrorCode = "30008"
	LLMStreamError         ErrorCode = "30009"

	// ============ 工具调用错误 (40000-40099) ============
	ToolError            ErrorCode = "40000"
	ToolNotFound         ErrorCode = "40001"
	ToolNotImplemented   ErrorCode = "40002"
	ToolInvalidParams    ErrorCode = "40003"
	ToolExecutionError   ErrorCode = "40004"
	ToolPermissionDenied ErrorCode = "40005"
	ToolTimeout          ErrorCode = "40006"
	ToolRateLimit        ErrorCode = "40007"
	ToolConnectionError  ErrorCode = "40008"

	// ============ 文件/Git/Shell 错误 (50000-50099) ============
	FileError            ErrorCode = "50000"
	FileNotFound         ErrorCode = "50001"
	FileReadError        ErrorCode = "50002"
	FileWriteError       ErrorCode = "50003"
	FilePermissionDenied ErrorCode = "50004"
	FileTooLarge         ErrorCode = "50005"
	FileFormatError      ErrorCode = "50006"
	GitError             ErrorCode = "50100"
	GitCloneError        ErrorCode = "50101"
	GitPushError         ErrorCode = "50102"
	GitPullError         ErrorCode = "50103"
	GitBranchError       ErrorCode = "50104"
	GitAuthError         ErrorCode = "50105"
	ShellError           ErrorCode = "50200"
	ShellExecutionError  ErrorCode = "50201"
	ShellTimeout         ErrorCode = "50202"
	ShellPermissionError ErrorCode = "50203"

	// ============ 知识库错误 (60000-60099) ============
	KnowledgeError     ErrorCode = "60000"
	KnowledgeNotFound  ErrorCode = "60001"
	KnowledgeNotReady  ErrorCode = "60002"
	EmbeddingError     ErrorCode = "60003"
	VectorStoreError   ErrorCode = "60004"
	IndexError         ErrorCode = "60005"
	DocumentParseError ErrorCode = "60006"
	DocumentTooLong    ErrorCode = "60007"
	ChunkError         ErrorCode = "60008"
	RAGError           ErrorCode = "60009"
)

var errorMessages = map[ErrorCode]string{
	// 通用错误
	Success:            "success",
	InternalError:      "系统内部错误",
	InvalidArgument:    "参数无效",
	NotFound:           "资源不存在",
	Unauthorized:       "未授权访问",
	Forbidden:          "禁止访问",
	AlreadyExists:      "资源已存在",
	ValidationError:    "数据校验失败",
	ServiceUnavailable: "服务不可用",
	Timeout:            "请求超时",
	DatabaseError:      "数据库操作失败",
	CacheError:         "缓存操作失败",
	NetworkError:       "网络连接错误",

	// Agent Run 错误
	RunNotFound:        "任务不存在",
	RunNotStarted:      "任务未启动",
	RunAlreadyStarted:  "任务已启动",
	RunAlreadyFinished: "任务已完成",
	RunCancelled:       "任务已取消",
	RunFailed:          "任务执行失败",
	RunTimeout:         "任务执行超时",
	RunMaxRetries:      "任务重试次数超限",
	RunInvalidStatus:   "任务状态无效",
	RunMissingAgent:    "Agent 配置不存在",

	// 模型调用错误
	LLMError:               "LLM 调用失败",
	LLMServiceError:        "LLM 服务错误",
	LLMAuthenticationError: "LLM 认证失败",
	LLMRateLimit:           "LLM 请求受限",
	LLMQuotaExceeded:       "LLM 配额不足",
	LLMInvalidRequest:      "LLM 请求无效",
	LLMModelNotFound:       "LLM 模型不存在",
	LLMContextOverlimit:    "LLM 上下文超限",
	LLMGenerationError:     "LLM 生成失败",
	LLMStreamError:         "LLM 流输出错误",

	// 工具调用错误
	ToolError:            "工具调用失败",
	ToolNotFound:         "工具不存在",
	ToolNotImplemented:   "工具未实现",
	ToolInvalidParams:    "工具参数无效",
	ToolExecutionError:   "工具执行错误",
	ToolPermissionDenied: "工具权限不足",
	ToolTimeout:          "工具执行超时",
	ToolRateLimit:        "工具调用受限",
	ToolConnectionError:  "工具连接错误",

	// 文件/Git/Shell 错误
	FileError:            "文件操作失败",
	FileNotFound:         "文件不存在",
	FileReadError:        "文件读取失败",
	FileWriteError:       "文件写入失败",
	FilePermissionDenied: "文件权限不足",
	FileTooLarge:         "文件过大",
	FileFormatError:      "文件格式错误",
	GitError:             "Git 操作失败",
	GitCloneError:        "Git 克隆失败",
	GitPushError:         "Git 推送失败",
	GitPullError:         "Git 拉取失败",
	GitBranchError:       "Git 分支操作失败",
	GitAuthError:         "Git 认证失败",
	ShellError:           "Shell 执行失败",
	ShellExecutionError:  "Shell 命令执行错误",
	ShellTimeout:         "Shell 执行超时",
	ShellPermissionError: "Shell 权限不足",

	// 知识库错误
	KnowledgeError:     "知识库操作失败",
	KnowledgeNotFound:  "知识库不存在",
	KnowledgeNotReady:  "知识库未就绪",
	EmbeddingError:     "向量嵌入失败",
	VectorStoreError:   "向量存储操作失败",
	IndexError:         "索引操作失败",
	DocumentParseError: "文档解析失败",
	DocumentTooLong:    "文档过长",
	ChunkError:         "文档切分失败",
	RAGError:           "RAG 查询失败",
}

type BizError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *BizError) Unwrap() error {
	return e.Err
}

func NewBizError(code ErrorCode, message ...string) *BizError {
	msg := errorMessages[code]
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &BizError{Code: code, Message: msg}
}

func NewBizErrorWithErr(code ErrorCode, err error, message ...string) *BizError {
	msg := errorMessages[code]
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &BizError{Code: code, Message: msg, Err: err}
}

func WrapError(code ErrorCode, err error, message ...string) *BizError {
	return NewBizErrorWithErr(code, err, message...)
}

// ============ 通用错误构造函数 ============
func Internal(message ...string) *BizError {
	return NewBizError(InternalError, message...)
}

func InternalWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(InternalError, err, message...)
}

func InvalidArg(message ...string) *BizError {
	return NewBizError(InvalidArgument, message...)
}

func NotFoundErr(message ...string) *BizError {
	return NewBizError(NotFound, message...)
}

func UnauthorizedErr(message ...string) *BizError {
	return NewBizError(Unauthorized, message...)
}

func ForbiddenErr(message ...string) *BizError {
	return NewBizError(Forbidden, message...)
}

func AlreadyExistsErr(message ...string) *BizError {
	return NewBizError(AlreadyExists, message...)
}

func Validation(message ...string) *BizError {
	return NewBizError(ValidationError, message...)
}

func ServiceUnavailableErr(message ...string) *BizError {
	return NewBizError(ServiceUnavailable, message...)
}

func TimeoutErr(message ...string) *BizError {
	return NewBizError(Timeout, message...)
}

func DatabaseErr(message ...string) *BizError {
	return NewBizError(DatabaseError, message...)
}

func DatabaseErrWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(DatabaseError, err, message...)
}

func CacheErr(message ...string) *BizError {
	return NewBizError(CacheError, message...)
}

func NetworkErr(message ...string) *BizError {
	return NewBizError(NetworkError, message...)
}

// ============ Agent Run 错误构造函数 ============
func RunNotFoundErr(message ...string) *BizError {
	return NewBizError(RunNotFound, message...)
}

func RunNotStartedErr(message ...string) *BizError {
	return NewBizError(RunNotStarted, message...)
}

func RunAlreadyStartedErr(message ...string) *BizError {
	return NewBizError(RunAlreadyStarted, message...)
}

func RunAlreadyFinishedErr(message ...string) *BizError {
	return NewBizError(RunAlreadyFinished, message...)
}

func RunCancelledErr(message ...string) *BizError {
	return NewBizError(RunCancelled, message...)
}

func RunFailedErr(message ...string) *BizError {
	return NewBizError(RunFailed, message...)
}

func RunFailedWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(RunFailed, err, message...)
}

func RunTimeoutErr(message ...string) *BizError {
	return NewBizError(RunTimeout, message...)
}

func RunMaxRetriesErr(message ...string) *BizError {
	return NewBizError(RunMaxRetries, message...)
}

func RunInvalidStatusErr(message ...string) *BizError {
	return NewBizError(RunInvalidStatus, message...)
}

func RunMissingAgentErr(message ...string) *BizError {
	return NewBizError(RunMissingAgent, message...)
}

// ============ 模型调用错误构造函数 ============
func LLMErr(message ...string) *BizError {
	return NewBizError(LLMError, message...)
}

func LLMErrWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(LLMError, err, message...)
}

func LLMAuthErr(message ...string) *BizError {
	return NewBizError(LLMAuthenticationError, message...)
}

func LLMRateLimitErr(message ...string) *BizError {
	return NewBizError(LLMRateLimit, message...)
}

func LLMQuotaExceededErr(message ...string) *BizError {
	return NewBizError(LLMQuotaExceeded, message...)
}

func LLMModelNotFoundErr(message ...string) *BizError {
	return NewBizError(LLMModelNotFound, message...)
}

func LLMContextOverlimitErr(message ...string) *BizError {
	return NewBizError(LLMContextOverlimit, message...)
}

// ============ 工具调用错误构造函数 ============
func ToolErr(message ...string) *BizError {
	return NewBizError(ToolError, message...)
}

func ToolErrWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(ToolError, err, message...)
}

func ToolNotFoundErr(message ...string) *BizError {
	return NewBizError(ToolNotFound, message...)
}

func ToolInvalidParamsErr(message ...string) *BizError {
	return NewBizError(ToolInvalidParams, message...)
}

func ToolPermissionDeniedErr(message ...string) *BizError {
	return NewBizError(ToolPermissionDenied, message...)
}

func ToolTimeoutErr(message ...string) *BizError {
	return NewBizError(ToolTimeout, message...)
}

// ============ 文件/Git/Shell 错误构造函数 ============
func FileErr(message ...string) *BizError {
	return NewBizError(FileError, message...)
}

func FileNotFoundErr(message ...string) *BizError {
	return NewBizError(FileNotFound, message...)
}

func FileReadErr(message ...string) *BizError {
	return NewBizError(FileReadError, message...)
}

func FileWriteErr(message ...string) *BizError {
	return NewBizError(FileWriteError, message...)
}

func FilePermissionDeniedErr(message ...string) *BizError {
	return NewBizError(FilePermissionDenied, message...)
}

func GitErr(message ...string) *BizError {
	return NewBizError(GitError, message...)
}

func GitCloneErr(message ...string) *BizError {
	return NewBizError(GitCloneError, message...)
}

func GitAuthErr(message ...string) *BizError {
	return NewBizError(GitAuthError, message...)
}

func ShellErr(message ...string) *BizError {
	return NewBizError(ShellError, message...)
}

func ShellErrWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(ShellError, err, message...)
}

// ============ 知识库错误构造函数 ============
func KnowledgeErr(message ...string) *BizError {
	return NewBizError(KnowledgeError, message...)
}

func KnowledgeNotFoundErr(message ...string) *BizError {
	return NewBizError(KnowledgeNotFound, message...)
}

func EmbeddingErr(message ...string) *BizError {
	return NewBizError(EmbeddingError, message...)
}

func EmbeddingErrWithErr(err error, message ...string) *BizError {
	return NewBizErrorWithErr(EmbeddingError, err, message...)
}

func VectorStoreErr(message ...string) *BizError {
	return NewBizError(VectorStoreError, message...)
}

func DocumentParseErr(message ...string) *BizError {
	return NewBizError(DocumentParseError, message...)
}

func RAGErr(message ...string) *BizError {
	return NewBizError(RAGError, message...)
}
