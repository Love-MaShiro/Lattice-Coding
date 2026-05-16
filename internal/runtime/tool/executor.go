package tool

import (
	"context"
	"time"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/errors"
)

type ToolExecutor struct {
	registry           *ToolRegistry
	resultProcessor    ResultProcessor
	safetyChecker      SafetyChecker
	auditRecorder      AuditRecorder
	invocationRecorder ToolInvocationRecorder
}

type Executor = ToolExecutor

type ExecutorOption func(*ToolExecutor)

func WithResultProcessor(processor ResultProcessor) ExecutorOption {
	return func(e *ToolExecutor) {
		if processor != nil {
			e.resultProcessor = processor
		}
	}
}

func WithSafetyChecker(checker SafetyChecker) ExecutorOption {
	return func(e *ToolExecutor) {
		if checker != nil {
			e.safetyChecker = checker
		}
	}
}

func WithAuditRecorder(recorder AuditRecorder) ExecutorOption {
	return func(e *ToolExecutor) {
		if recorder != nil {
			e.auditRecorder = recorder
		}
	}
}

func WithToolInvocationRecorder(recorder ToolInvocationRecorder) ExecutorOption {
	return func(e *ToolExecutor) {
		if recorder != nil {
			e.invocationRecorder = recorder
		}
	}
}

func NewToolExecutor(registry *ToolRegistry, opts ...ExecutorOption) *ToolExecutor {
	if registry == nil {
		registry = NewRegistry()
	}
	e := &ToolExecutor{
		registry:           registry,
		resultProcessor:    ResultProcessorFunc(func(ctx context.Context, req ToolRequest, output ToolOutput) (ToolOutput, error) { return output, nil }),
		safetyChecker:      NewRuleBasedSafetyChecker(),
		auditRecorder:      NoopAuditRecorder{},
		invocationRecorder: NoopToolInvocationRecorder{},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}
	return e
}

func NewExecutor(registry *ToolRegistry, opts ...ExecutorOption) *ToolExecutor {
	return NewToolExecutor(registry, opts...)
}

func (e *ToolExecutor) Execute(ctx context.Context, req ToolRequest) (result ToolResult) {
	startedAt := time.Now()
	var descriptor ToolDescriptor
	var invocationID string
	result = ToolResult{
		RequestID: req.ID,
		ToolName:  req.Name,
		StartedAt: startedAt,
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			result.IsError = true
			result.Error = errors.ToolErr("tool runtime panic").Error()
			result.Content = result.Error
		}
		finishedAt := time.Now()
		result.FinishedAt = finishedAt
		result.DurationMs = finishedAt.Sub(startedAt).Milliseconds()
		if e.auditRecorder != nil {
			_ = e.auditRecorder.Record(ctx, AuditEvent{
				EventType:  auditEventType(req, result),
				Request:    req,
				Descriptor: descriptor,
				Result:     result,
				StartedAt:  startedAt,
				FinishedAt: finishedAt,
			})
		}
		if e.invocationRecorder != nil && descriptor.Name != "" {
			status := "success"
			if result.IsError {
				status = "failed"
			}
			_ = e.invocationRecorder.Finish(ctx, ToolInvocationFinished{
				ID:            invocationID,
				RunID:         req.Context.RunID,
				NodeID:        nodeIDFromRequest(req),
				ToolName:      req.Name,
				Result:        result,
				Status:        status,
				CompletedAt:   finishedAt,
				LatencyMs:     result.DurationMs,
				FullResultRef: result.FullResultRef,
			})
		}
	}()

	if req.Name == "" {
		return e.errorResult(result, errors.ToolInvalidParamsErr("tool name is required"))
	}
	if req.Input == nil {
		req.Input = map[string]interface{}{}
	}

	t, ok := e.registry.Get(req.Name)
	if !ok {
		return e.errorResult(result, errors.ToolNotFoundErr("tool not found: "+req.Name))
	}
	descriptor = DescriptorOf(t)
	if e.auditRecorder != nil {
		_ = e.auditRecorder.Record(ctx, AuditEvent{
			EventType:  "tool_called",
			Request:    req,
			Descriptor: descriptor,
			StartedAt:  startedAt,
		})
	}
	if e.invocationRecorder != nil {
		id, err := e.invocationRecorder.Start(ctx, ToolInvocationStarted{
			RunID:     req.Context.RunID,
			NodeID:    nodeIDFromRequest(req),
			ToolName:  req.Name,
			Input:     req.Input,
			StartedAt: startedAt,
		})
		if err != nil {
			return e.errorResult(result, errors.ToolErrWithErr(err, "tool invocation record failed"))
		}
		invocationID = id
	}

	if err := t.Validate(ctx, req.Input); err != nil {
		return e.errorResult(result, errors.ToolInvalidParamsErr(err.Error()))
	}

	if e.safetyChecker != nil {
		safetyResult, err := e.safetyChecker.Check(ctx, req, descriptor)
		if err != nil {
			return e.errorResult(result, errors.ToolPermissionDeniedErr(err.Error()))
		}
		if safetyResult.Decision != SafetyAllow {
			reason := safetyResult.Reason
			if reason == "" {
				reason = "tool safety check denied"
			}
			if safetyResult.Decision == SafetyNeedApproval {
				reason = "tool safety check requires approval: " + reason
			}
			return e.errorResult(result, errors.ToolPermissionDeniedErr(reason))
		}
	}

	decision, reason, err := t.CheckPermission(ctx, req)
	if err != nil {
		return e.errorResult(result, errors.ToolPermissionDeniedErr(err.Error()))
	}
	if decision != PermissionAllow {
		if reason == "" {
			reason = "tool permission denied"
		}
		return e.errorResult(result, errors.ToolPermissionDeniedErr(reason))
	}

	output, err := t.Execute(ctx, req)
	if err != nil {
		return e.errorResultWithOutput(result, output, errors.ToolErrWithErr(err, "tool execution failed"))
	}

	if e.resultProcessor != nil {
		output, err = e.resultProcessor.Process(ctx, req, output)
		if err != nil {
			return e.errorResult(result, errors.ToolErrWithErr(err, "tool result processing failed"))
		}
	}

	result.Content = output.Content
	result.Data = output.Data
	result.Metadata = output.Metadata
	result.FullResultRef = fullResultRefFromMetadata(output.Metadata)
	result.Truncated = output.Truncated
	return result
}

func (e *ToolExecutor) Register(t Tool) error {
	return e.registry.Register(t)
}

func (e *ToolExecutor) Registry() *ToolRegistry {
	return e.registry
}

func (e *ToolExecutor) SetToolInvocationRecorder(recorder ToolInvocationRecorder) {
	if recorder != nil {
		e.invocationRecorder = recorder
	}
}

func (e *ToolExecutor) SetAuditRecorder(recorder AuditRecorder) {
	if recorder != nil {
		e.auditRecorder = recorder
	}
}

func (e *ToolExecutor) ListDescriptors() []ToolDescriptor {
	return e.registry.ListDescriptors()
}

func (e *ToolExecutor) errorResult(result ToolResult, err error) ToolResult {
	result.IsError = true
	result.Error = ErrorContent(err)
	result.Content = result.Error
	return result
}

func (e *ToolExecutor) errorResultWithOutput(result ToolResult, output ToolOutput, err error) ToolResult {
	result = e.errorResult(result, err)
	if output.Content != "" {
		result.Content = output.Content
	}
	result.Data = output.Data
	result.Metadata = output.Metadata
	result.FullResultRef = fullResultRefFromMetadata(output.Metadata)
	result.Truncated = output.Truncated
	return result
}

func auditEventType(req ToolRequest, result ToolResult) string {
	if result.IsError && result.Data == nil {
		return "tool_denied"
	}
	switch req.Name {
	case "file.edit":
		return "file_edit"
	case "shell.run":
		return "shell_run"
	default:
		return "tool_finished"
	}
}

func nodeIDFromRequest(req ToolRequest) string {
	if req.Context.Metadata == nil {
		return ""
	}
	if value, ok := req.Context.Metadata["node_id"].(string); ok {
		return value
	}
	return ""
}

func fullResultRefFromMetadata(metadata map[string]interface{}) string {
	if metadata == nil {
		return ""
	}
	if value, ok := metadata["full_result_ref"].(string); ok {
		return value
	}
	return ""
}

var defaultExecutor *ToolExecutor

func Init(cfg *config.Config) {
	_ = cfg
	defaultExecutor = NewExecutor(NewRegistry())
}

func Default() *Executor {
	if defaultExecutor == nil {
		Init(nil)
	}
	return defaultExecutor
}

func Register(t Tool) error {
	return Default().Register(t)
}

func Execute(ctx context.Context, req ToolRequest) ToolResult {
	return Default().Execute(ctx, req)
}

func List() []ToolDescriptor {
	return Default().ListDescriptors()
}
