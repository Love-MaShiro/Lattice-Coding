package strategy

import (
	"context"
	"testing"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/runtime/llm"
	"lattice-coding/internal/runtime/query"
)

func TestDirectChatStrategy_Execute_ShouldUseHistoryMessages(t *testing.T) {
	llmExecutor := llm.NewExecutor(&config.LLMConfig{
		Pool:   config.PoolConfig{MaxConcurrent: 1},
		Stream: config.PoolConfig{MaxConcurrent: 1},
		Routing: config.RoutingConfig{
			Default: config.RouteConfig{Primary: "fake"},
		},
	})
	client := &recordingLLMClient{response: "done"}
	llmExecutor.RegisterClient("fake", client)

	state := query.NewState(query.QueryRequest{
		Input:         "current question",
		Provider:      "fake",
		SystemPrompt:  "system prompt",
		Summary:       "summary text",
		ModelConfigID: 7,
		Messages: []llm.Message{
			{Role: "user", Content: "old question"},
			{Role: "assistant", Content: "old answer"},
		},
	})

	result, err := NewDirectChatStrategy(llmExecutor).Execute(context.Background(), state)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result.FinalAnswer != "done" {
		t.Fatalf("FinalAnswer = %q, want done", result.FinalAnswer)
	}

	got := client.lastRequest.Messages
	want := []llm.Message{
		{Role: "system", Content: "system prompt"},
		{Role: "system", Content: "Previous conversation summary:\nsummary text"},
		{Role: "user", Content: "old question"},
		{Role: "assistant", Content: "old answer"},
		{Role: "user", Content: "current question"},
	}
	if len(got) != len(want) {
		t.Fatalf("message count = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("message[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

type recordingLLMClient struct {
	response    string
	chunks      []llm.StreamChunk
	lastRequest llm.ChatRequest
}

func (c *recordingLLMClient) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.lastRequest = req
	return &llm.ChatResponse{Content: c.response}, nil
}

func (c *recordingLLMClient) Stream(ctx context.Context, req llm.ChatRequest) (<-chan llm.StreamChunk, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.lastRequest = req
	out := make(chan llm.StreamChunk, len(c.chunks))
	go func() {
		defer close(out)
		for _, chunk := range c.chunks {
			out <- chunk
		}
	}()
	return out, nil
}

func (c *recordingLLMClient) Close() error {
	return nil
}

func TestDirectChatStrategy_Stream_ShouldForwardDelta(t *testing.T) {
	llmExecutor := llm.NewExecutor(&config.LLMConfig{
		Pool:   config.PoolConfig{MaxConcurrent: 1},
		Stream: config.PoolConfig{MaxConcurrent: 1},
		Routing: config.RoutingConfig{
			Default: config.RouteConfig{Primary: "fake"},
		},
	})
	client := &recordingLLMClient{
		chunks: []llm.StreamChunk{
			{Content: "hel"},
			{Content: "lo"},
			{Done: true},
		},
	}
	llmExecutor.RegisterClient("fake", client)

	state := query.NewState(query.QueryRequest{
		RunID:        "run-1",
		Input:        "current question",
		Provider:     "fake",
		SystemPrompt: "system prompt",
		Messages: []llm.Message{
			{Role: "user", Content: "old question"},
		},
	})

	stream, err := NewDirectChatStrategy(llmExecutor).Stream(context.Background(), state)
	if err != nil {
		t.Fatalf("Stream returned error: %v", err)
	}

	var events []query.StreamEvent
	for event := range stream {
		events = append(events, event)
	}

	if len(events) != 5 {
		t.Fatalf("event count = %d, want 5: %+v", len(events), events)
	}
	if events[0].Type != query.StreamEventRunStarted {
		t.Fatalf("first event = %s", events[0].Type)
	}
	if events[1].Type != query.StreamEventLLMDelta || events[1].Content != "hel" {
		t.Fatalf("first delta = %+v", events[1])
	}
	if events[2].Type != query.StreamEventLLMDelta || events[2].Content != "lo" {
		t.Fatalf("second delta = %+v", events[2])
	}
	if events[3].Type != query.StreamEventLLMDone || events[3].Content != "hello" {
		t.Fatalf("llm done = %+v", events[3])
	}
	if events[4].Type != query.StreamEventRunFinished || events[4].Content != "hello" {
		t.Fatalf("run finished = %+v", events[4])
	}
	if got := client.lastRequest.Messages[len(client.lastRequest.Messages)-1]; got.Content != "current question" {
		t.Fatalf("last message = %+v", got)
	}
}
