package agent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lattice-coding/internal/runtime/llm"
	"lattice-coding/internal/runtime/tool"
)

const (
	defaultMaxIterations = 8
	defaultMaxToolCalls  = 6
	defaultTimeout       = 2 * time.Minute
)

type reactStrategy struct {
	llm        *llm.Executor
	tools      *tool.Executor
	scratchpad Scratchpad
}

func NewReActStrategy(llmExecutor *llm.Executor, toolExecutor *tool.Executor, scratchpad Scratchpad) ExecutionStrategy {
	return &reactStrategy{
		llm:        llmExecutor,
		tools:      toolExecutor,
		scratchpad: scratchpad,
	}
}

func (s *reactStrategy) Name() string {
	return StrategyReAct
}

func (s *reactStrategy) Execute(ctx context.Context, req Request) (*Result, error) {
	if s.llm == nil {
		return nil, errors.New("react strategy requires llm executor")
	}
	if s.tools == nil {
		return nil, errors.New("react strategy requires tool executor")
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	maxIterations := req.MaxIterations
	if maxIterations <= 0 {
		maxIterations = defaultMaxIterations
	}
	maxToolCalls := req.MaxToolCalls
	if maxToolCalls <= 0 {
		maxToolCalls = defaultMaxToolCalls
	}

	state := &state{request: req}
	allowed := allowedTools(req.AllowedTools)
	scratchpad := s.scratchpad
	if scratchpad == nil {
		scratchpad = NewMemoryScratchpad()
	}

	for iteration := 1; iteration <= maxIterations; iteration++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		descriptors := filterToolDescriptors(s.tools.ListDescriptors(), allowed)
		messages := BuildReActMessages(req, descriptors, state.steps)
		resp, call := s.llm.Chat(ctx, llm.ChatRequest{
			Provider:      req.Provider,
			Model:         req.Model,
			ModelConfigID: req.ModelConfigID,
			Messages: []llm.Message{
				{Role: "system", Content: messages[0]},
				{Role: "user", Content: messages[1]},
			},
		})
		if !call.Success {
			if call.Error != nil {
				return nil, call.Error
			}
			return nil, errors.New("llm call failed")
		}
		if resp == nil {
			return nil, errors.New("llm returned empty response")
		}

		action, err := ParseReActAction(resp.Content)
		if err != nil {
			step := errorStep(iteration, "parse_react_action", resp.Content, err)
			state.appendStep(step)
			_ = scratchpad.Append(ctx, step)
			continue
		}

		if action.Type == ReActActionFinal {
			step := finalStep(iteration, action.Answer)
			state.appendStep(step)
			if err := scratchpad.Append(ctx, step); err != nil {
				return nil, err
			}
			return &Result{
				RunID:   req.RunID,
				Content: action.Answer,
				Messages: []Message{
					{Role: "user", Content: req.Input},
					{Role: "assistant", Content: action.Answer},
				},
				Metadata: map[string]interface{}{
					"strategy":    StrategyReAct,
					"react_steps": state.steps,
				},
			}, nil
		}

		if state.toolCalls >= maxToolCalls {
			step := errorStep(iteration, "react", "max tool calls reached", errors.New("react max tool calls reached"))
			state.appendStep(step)
			_ = scratchpad.Append(ctx, step)
			return reactFailureResult(req, "未能完成请求：ReAct 已达到工具调用上限。请查看工具调用轨迹了解失败原因。", state.steps), nil
		}

		step := s.executeToolCall(ctx, req, iteration, action, allowed)
		state.appendStep(step)
		if err := scratchpad.Append(ctx, step); err != nil {
			return nil, err
		}
	}

	step := errorStep(maxIterations, "react", "max iterations reached", errors.New("react max iterations reached"))
	state.appendStep(step)
	_ = scratchpad.Append(ctx, step)
	return reactFailureResult(req, "未能完成请求：ReAct 已达到最大迭代次数。请查看工具调用轨迹了解失败原因。", state.steps), nil
}

func reactFailureResult(req Request, answer string, steps []ReActStep) *Result {
	return &Result{
		RunID:   req.RunID,
		Content: answer,
		Messages: []Message{
			{Role: "user", Content: req.Input},
			{Role: "assistant", Content: answer},
		},
		Metadata: map[string]interface{}{
			"strategy":    StrategyReAct,
			"react_steps": steps,
			"is_error":    true,
		},
	}
}

func (s *reactStrategy) executeToolCall(ctx context.Context, req Request, iteration int, action *ReActAction, allowed map[string]struct{}) ReActStep {
	startedAt := time.Now()
	step := ReActStep{
		Iteration:   iteration,
		Reason:      action.Reason,
		Action:      action.Tool,
		ActionInput: action.Args,
		StartedAt:   startedAt,
	}

	if _, ok := allowed[action.Tool]; !ok {
		step.IsError = true
		step.Observation = "tool is not allowed: " + action.Tool
		step.CompletedAt = time.Now()
		return step
	}

	result := s.tools.Execute(ctx, tool.ToolRequest{
		Name:  action.Tool,
		Input: action.Args,
		Context: tool.ToolContext{
			RunID:      req.RunID,
			AgentID:    req.AgentID,
			UserID:     req.UserID,
			ProjectID:  req.ProjectID,
			SessionID:  req.SessionID,
			WorkingDir: req.WorkingDir,
			Metadata: map[string]interface{}{
				"node_id": req.NodeID,
			},
		},
	})
	step.IsError = result.IsError
	step.Observation = result.Content
	if result.IsError && result.Error != "" {
		step.Observation = result.Error
	}
	step.CompletedAt = time.Now()
	return step
}

func allowedTools(names []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(names))
	for _, name := range names {
		if name != "" {
			allowed[name] = struct{}{}
		}
	}
	return allowed
}

func filterToolDescriptors(descriptors []tool.ToolDescriptor, allowed map[string]struct{}) []tool.ToolDescriptor {
	if len(allowed) == 0 {
		return nil
	}
	filtered := make([]tool.ToolDescriptor, 0, len(descriptors))
	for _, descriptor := range descriptors {
		if _, ok := allowed[descriptor.Name]; ok {
			filtered = append(filtered, descriptor)
		}
	}
	return filtered
}

func errorStep(iteration int, action string, observation string, err error) ReActStep {
	now := time.Now()
	return ReActStep{
		Iteration:   iteration,
		Action:      action,
		Observation: fmt.Sprintf("%s: %s", err.Error(), observation),
		IsError:     true,
		StartedAt:   now,
		CompletedAt: now,
	}
}

func finalStep(iteration int, answer string) ReActStep {
	now := time.Now()
	return ReActStep{
		Iteration:   iteration,
		Action:      ReActActionFinal,
		Observation: answer,
		StartedAt:   now,
		CompletedAt: now,
	}
}
