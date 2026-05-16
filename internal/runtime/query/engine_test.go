package query

import (
	"context"
	"testing"
)

func TestQueryEngine_Stream_ShouldSelectDirectChatStrategy(t *testing.T) {
	strategy := &fakeStreamStrategy{
		events: []StreamEvent{
			{Type: StreamEventRunStarted, RunID: "run-1"},
			{Type: StreamEventLLMDelta, RunID: "run-1", Content: "hello"},
			{Type: StreamEventRunFinished, RunID: "run-1", Content: "hello", Done: true},
		},
	}
	engine := NewEngine(WithStrategy(strategy))

	stream, err := engine.Stream(context.Background(), QueryRequest{
		RunID: "run-1",
		Mode:  ExecutionModeDirectChat,
		Input: "hello",
	})
	if err != nil {
		t.Fatalf("Stream returned error: %v", err)
	}

	var events []StreamEvent
	for event := range stream {
		events = append(events, event)
	}

	if strategy.streamCalls != 1 {
		t.Fatalf("stream calls = %d, want 1", strategy.streamCalls)
	}
	if !strategy.lastRequest.Stream {
		t.Fatal("request Stream flag = false, want true")
	}
	if len(events) != len(strategy.events) {
		t.Fatalf("event count = %d, want %d", len(events), len(strategy.events))
	}
	if events[1].Content != "hello" {
		t.Fatalf("delta content = %q", events[1].Content)
	}
}

type fakeStreamStrategy struct {
	events      []StreamEvent
	streamCalls int
	lastRequest QueryRequest
}

func (s *fakeStreamStrategy) Mode() ExecutionMode {
	return ExecutionModeDirectChat
}

func (s *fakeStreamStrategy) Execute(context.Context, *QueryState) (*QueryResult, error) {
	return nil, nil
}

func (s *fakeStreamStrategy) Stream(_ context.Context, state *QueryState) (QueryStream, error) {
	s.streamCalls++
	s.lastRequest = state.Request
	out := make(chan StreamEvent, len(s.events))
	go func() {
		defer close(out)
		for _, event := range s.events {
			out <- event
		}
	}()
	return out, nil
}
