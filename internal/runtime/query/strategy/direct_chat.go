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
	resp, call := s.llm.Call(ctx, llm.ChatRequest{
		Provider:      req.Provider,
		Model:         req.Model,
		ModelConfigID: req.ModelConfigID,
		Messages:      buildDirectChatMessages(req),
		Temperature:   req.Temperature,
		TopP:          req.TopP,
		MaxTokens:     req.MaxTokens,
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
		RunID:       state.RunID,
		Mode:        query.ExecutionModeDirectChat,
		Content:     resp.Content,
		FinalAnswer: resp.Content,
		Messages:    directChatQueryMessages(append(buildDirectChatMessages(req), llm.Message{Role: "assistant", Content: resp.Content})),
		Steps:       state.Steps,
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

func (s *DirectChatStrategy) Stream(ctx context.Context, state *query.QueryState) (query.QueryStream, error) {
	if s.llm == nil {
		return nil, errors.New("direct chat strategy requires llm executor")
	}
	if err := state.Budget.UseStep(); err != nil {
		return nil, err
	}

	req := state.Request
	llmStream, call := s.llm.Stream(ctx, llm.ChatRequest{
		Provider:      req.Provider,
		Model:         req.Model,
		ModelConfigID: req.ModelConfigID,
		Messages:      buildDirectChatMessages(req),
		Temperature:   req.Temperature,
		TopP:          req.TopP,
		MaxTokens:     req.MaxTokens,
	})
	if !call.Success {
		if call.Error != nil {
			return nil, call.Error
		}
		return nil, errors.New("llm stream failed")
	}

	out := make(chan query.StreamEvent, 16)
	go func() {
		defer close(out)
		defer state.Finish()

		out <- query.StreamEvent{
			Type:  query.StreamEventRunStarted,
			RunID: state.RunID,
			Metadata: map[string]interface{}{
				"mode":       query.ExecutionModeDirectChat.String(),
				"started_at": state.StartedAt,
			},
		}

		var content string
		for chunk := range llmStream {
			if chunk.Err != nil {
				out <- query.StreamEvent{
					Type:    query.StreamEventRunError,
					RunID:   state.RunID,
					Err:     chunk.Err,
					Done:    true,
					Content: content,
				}
				return
			}
			if chunk.Usage != nil {
				_ = state.Budget.UseTokens(chunk.Usage.TotalTokens)
			}
			if chunk.Content != "" {
				content += chunk.Content
				out <- query.StreamEvent{
					Type:    query.StreamEventLLMDelta,
					RunID:   state.RunID,
					Content: chunk.Content,
				}
			}
			if chunk.Done {
				out <- query.StreamEvent{
					Type:    query.StreamEventLLMDone,
					RunID:   state.RunID,
					Content: content,
					Done:    true,
				}
			}
		}

		state.Steps = append(state.Steps, query.StepResult{
			Iteration:   1,
			Name:        query.ExecutionModeDirectChat.String(),
			Content:     content,
			StartedAt:   state.StartedAt,
			CompletedAt: time.Now(),
		})
		out <- query.StreamEvent{
			Type:    query.StreamEventRunFinished,
			RunID:   state.RunID,
			Content: content,
			Done:    true,
			Metadata: map[string]interface{}{
				"mode":         query.ExecutionModeDirectChat.String(),
				"completed_at": time.Now(),
				"budget":       state.Budget.Snapshot(),
			},
		}
	}()

	return out, nil
}

func buildDirectChatMessages(req query.QueryRequest) []llm.Message {
	messages := make([]llm.Message, 0, len(req.Messages)+3)
	if req.SystemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: req.SystemPrompt})
	}
	if req.Summary != "" {
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: "Previous conversation summary:\n" + req.Summary,
		})
	}
	messages = append(messages, req.Messages...)
	if req.Input != "" {
		messages = append(messages, llm.Message{Role: "user", Content: req.Input})
	}
	return messages
}

func directChatQueryMessages(messages []llm.Message) []query.Message {
	result := make([]query.Message, 0, len(messages))
	for _, message := range messages {
		result = append(result, query.Message{Role: message.Role, Content: message.Content})
	}
	return result
}
