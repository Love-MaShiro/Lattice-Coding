package strategy

import (
	"context"
	"errors"
	"time"

	"lattice-coding/internal/runtime/agent"
	"lattice-coding/internal/runtime/query"
)

type PureReActStrategy struct {
	runtime agent.Runtime
}

func NewPureReActStrategy(runtime agent.Runtime) *PureReActStrategy {
	return &PureReActStrategy{runtime: runtime}
}

func (s *PureReActStrategy) Mode() query.ExecutionMode {
	return query.ExecutionModePureReAct
}

func (s *PureReActStrategy) Execute(ctx context.Context, state *query.QueryState) (*query.QueryResult, error) {
	if s.runtime == nil {
		return nil, errors.New("pure react strategy requires agent runtime")
	}
	if err := state.Budget.UseStep(); err != nil {
		return nil, err
	}

	req := state.Request
	result, err := s.runtime.Run(ctx, agent.Request{
		RunID:         state.RunID,
		AgentID:       req.AgentID,
		UserID:        req.UserID,
		SessionID:     req.SessionID,
		ProjectID:     req.ProjectID,
		NodeID:        req.NodeID,
		Input:         req.Input,
		Strategy:      agent.StrategyReAct,
		Provider:      req.Provider,
		Model:         req.Model,
		ModelConfigID: req.ModelConfigID,
		AllowedTools:  req.AllowedTools,
		WorkingDir:    req.WorkingDir,
		MaxIterations: req.Budget.MaxSteps,
		MaxToolCalls:  req.Budget.MaxToolCalls,
		Timeout:       req.Timeout,
		Metadata:      req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("pure react agent returned nil result")
	}

	steps := reactStepsFromMetadata(result.Metadata)
	if len(steps) > 0 {
		for range steps {
			if err := state.Budget.UseStep(); err != nil {
				return nil, err
			}
		}
		state.Steps = append(state.Steps, steps...)
	}

	return &query.QueryResult{
		RunID:       state.RunID,
		Mode:        query.ExecutionModePureReAct,
		Content:     result.Content,
		FinalAnswer: result.Content,
		Messages:    queryMessages(result.Messages),
		Steps:       state.Steps,
		Metadata:    result.Metadata,
	}, nil
}

func queryMessages(messages []agent.Message) []query.Message {
	if len(messages) == 0 {
		return nil
	}
	out := make([]query.Message, 0, len(messages))
	for _, message := range messages {
		out = append(out, query.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return out
}

func reactStepsFromMetadata(metadata map[string]interface{}) []query.StepResult {
	value, ok := metadata["react_steps"]
	if !ok {
		return nil
	}
	reactSteps, ok := value.([]agent.ReActStep)
	if !ok {
		return nil
	}
	steps := make([]query.StepResult, 0, len(reactSteps))
	for _, step := range reactSteps {
		name := step.Action
		if name == "" {
			name = "react"
		}
		steps = append(steps, query.StepResult{
			Iteration:   step.Iteration,
			Name:        name,
			Content:     step.Observation,
			IsError:     step.IsError,
			StartedAt:   step.StartedAt,
			CompletedAt: nonZeroTime(step.CompletedAt),
			Metadata: map[string]interface{}{
				"reason":       step.Reason,
				"action":       step.Action,
				"action_input": step.ActionInput,
			},
		})
	}
	return steps
}

func nonZeroTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now()
	}
	return value
}
