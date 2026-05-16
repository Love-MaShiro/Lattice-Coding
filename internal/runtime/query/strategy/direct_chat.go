package strategy

import (
	"context"
	"errors"
	"time"

	"lattice-coding/internal/runtime/llm"
	"lattice-coding/internal/runtime/query"
)

type DirectChatStrategy struct {
	llm *llm.Executor
}

func NewDirectChatStrategy(llmExecutor *llm.Executor) *DirectChatStrategy {
	return &DirectChatStrategy{llm: llmExecutor}
}

func (s *DirectChatStrategy) Mode() query.ExecutionMode {
	return query.ExecutionModeDirectChat
}

func (s *DirectChatStrategy) Execute(ctx context.Context, state *query.QueryState) (*query.QueryResult, error) {
	if s.llm == nil {
		return nil, errors.New("direct chat strategy requires llm executor")
	}
	if err := state.Budget.UseStep(); err != nil {
		return nil, err
	}

	startedAt := time.Now()
	req := state.Request
	resp, call := s.llm.Chat(ctx, llm.ChatRequest{
		Provider:      req.Provider,
		Model:         req.Model,
		ModelConfigID: req.ModelConfigID,
		Messages: []llm.Message{
			{Role: "user", Content: req.Input},
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
	if resp.Usage != nil {
		if err := state.Budget.UseTokens(resp.Usage.TotalTokens); err != nil {
			return nil, err
		}
	}

	step := query.StepResult{
		Iteration:   1,
		Name:        query.ExecutionModeDirectChat.String(),
		Content:     resp.Content,
		StartedAt:   startedAt,
		CompletedAt: time.Now(),
	}
	state.Steps = append(state.Steps, step)

	result := &query.QueryResult{
		RunID:   state.RunID,
		Mode:    query.ExecutionModeDirectChat,
		Content: resp.Content,
		Messages: []query.Message{
			{Role: "user", Content: req.Input},
			{Role: "assistant", Content: resp.Content},
		},
		Steps: state.Steps,
	}
	if resp.Usage != nil {
		result.Usage = query.Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}
	return result, nil
}
