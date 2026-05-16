package application

import (
	"context"
	"errors"
	"testing"
	"time"

	chatdomain "lattice-coding/internal/modules/chat/domain"
	"lattice-coding/internal/runtime/query"
)

func TestChatCompletion_ShouldCallQueryEngine(t *testing.T) {
	strategy := &fakeQueryStrategy{answer: "assistant answer"}
	service, messages := newCommandServiceForQueryTest(strategy)

	result, err := service.Complete(context.Background(), &CompletionCommand{
		AgentID: 1,
		Message: "hello",
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if result.Content != "assistant answer" {
		t.Fatalf("content = %q", result.Content)
	}
	if strategy.calls != 1 {
		t.Fatalf("query strategy calls = %d, want 1", strategy.calls)
	}
	if strategy.lastRequest.Input != "hello" {
		t.Fatalf("query input = %q", strategy.lastRequest.Input)
	}
	if strategy.lastRequest.ModelConfigID != 99 {
		t.Fatalf("model_config_id = %d, want 99", strategy.lastRequest.ModelConfigID)
	}
	if len(strategy.lastRequest.AllowedTools) == 0 {
		t.Fatal("expected static allowed tools")
	}
	if strategy.lastRequest.WorkingDir == "" {
		t.Fatal("expected working dir")
	}
	if got := messages.countRole(chatdomain.MessageRoleAssistant); got != 1 {
		t.Fatalf("assistant messages = %d, want 1", got)
	}
}

func TestChatCompletion_QueryEngineError_ShouldNotSaveAssistantMessage(t *testing.T) {
	strategy := &fakeQueryStrategy{err: errors.New("query failed")}
	service, messages := newCommandServiceForQueryTest(strategy)

	_, err := service.Complete(context.Background(), &CompletionCommand{
		AgentID: 1,
		Message: "hello",
	})
	if err == nil {
		t.Fatal("Complete returned nil error")
	}
	if strategy.calls != 1 {
		t.Fatalf("query strategy calls = %d, want 1", strategy.calls)
	}
	if got := messages.countRole(chatdomain.MessageRoleAssistant); got != 0 {
		t.Fatalf("assistant messages = %d, want 0", got)
	}
	if got := messages.countRole(chatdomain.MessageRoleUser); got != 1 {
		t.Fatalf("user messages = %d, want 1", got)
	}
}

func TestChatStream_ShouldCallQueryEngineStream(t *testing.T) {
	strategy := &fakeQueryStrategy{
		streamEvents: []query.StreamEvent{
			{Type: query.StreamEventRunStarted, RunID: "run-1"},
			{Type: query.StreamEventLLMDelta, RunID: "run-1", Content: "hel"},
			{Type: query.StreamEventLLMDelta, RunID: "run-1", Content: "lo"},
			{Type: query.StreamEventLLMDone, RunID: "run-1", Content: "hello", Done: true},
			{Type: query.StreamEventRunFinished, RunID: "run-1", Content: "hello", Done: true},
		},
	}
	service, messages := newCommandServiceForQueryTest(strategy)

	var events []query.StreamEvent
	result, err := service.StreamComplete(context.Background(), &CompletionCommand{
		AgentID: 1,
		Message: "hello",
	}, func(event query.StreamEvent) error {
		events = append(events, event)
		return nil
	})
	if err != nil {
		t.Fatalf("StreamComplete returned error: %v", err)
	}
	if result.Content != "hello" {
		t.Fatalf("content = %q, want hello", result.Content)
	}
	if strategy.streamCalls != 1 {
		t.Fatalf("query stream calls = %d, want 1", strategy.streamCalls)
	}
	if !strategy.lastRequest.Stream {
		t.Fatal("query request Stream = false, want true")
	}
	if got := messages.countRole(chatdomain.MessageRoleAssistant); got != 1 {
		t.Fatalf("assistant messages = %d, want 1", got)
	}
	if len(events) != 5 {
		t.Fatalf("events = %d, want 5", len(events))
	}
}

func TestChatStream_WhenError_ShouldReturnSSEError(t *testing.T) {
	streamErr := errors.New("stream failed")
	strategy := &fakeQueryStrategy{
		streamEvents: []query.StreamEvent{
			{Type: query.StreamEventRunStarted, RunID: "run-1"},
			{Type: query.StreamEventRunError, RunID: "run-1", Err: streamErr, Done: true},
		},
	}
	service, messages := newCommandServiceForQueryTest(strategy)

	var errorEvent *query.StreamEvent
	_, err := service.StreamComplete(context.Background(), &CompletionCommand{
		AgentID: 1,
		Message: "hello",
	}, func(event query.StreamEvent) error {
		if event.Type == query.StreamEventRunError {
			copied := event
			errorEvent = &copied
		}
		return nil
	})
	if err == nil {
		t.Fatal("StreamComplete returned nil error")
	}
	if errorEvent == nil || errorEvent.Err == nil {
		t.Fatalf("missing SSE error event: %+v", errorEvent)
	}
	if got := messages.countRole(chatdomain.MessageRoleAssistant); got != 0 {
		t.Fatalf("assistant messages = %d, want 0", got)
	}
}

func newCommandServiceForQueryTest(strategy *fakeQueryStrategy) (*CommandService, *fakeMessageRepo) {
	sessionRepo := &fakeSessionRepo{}
	messageRepo := &fakeMessageRepo{}
	agentGetter := fakeAgentGetter{agent: &AgentRuntimeDTO{
		ID:            1,
		Name:          "agent",
		ModelConfigID: 99,
		SystemPrompt:  "system",
		Temperature:   0.2,
		TopP:          0.8,
		MaxTokens:     128,
		Enabled:       true,
	}}
	queryEngine := query.NewEngine(query.WithStrategy(strategy))
	return NewCommandService(
		sessionRepo,
		messageRepo,
		agentGetter,
		queryEngine,
		nil,
		nil,
		MemoryConfig{CompressionThreshold: 100, RetainAfterCompress: 20, CacheTTL: time.Hour},
	), messageRepo
}

type fakeQueryStrategy struct {
	answer       string
	err          error
	calls        int
	streamCalls  int
	streamEvents []query.StreamEvent
	streamErr    error
	lastRequest  query.QueryRequest
}

func (s *fakeQueryStrategy) Mode() query.ExecutionMode {
	return query.ExecutionModeDirectChat
}

func (s *fakeQueryStrategy) Execute(_ context.Context, state *query.QueryState) (*query.QueryResult, error) {
	s.calls++
	s.lastRequest = state.Request
	if s.err != nil {
		return nil, s.err
	}
	return &query.QueryResult{
		Content:     s.answer,
		FinalAnswer: s.answer,
	}, nil
}

func (s *fakeQueryStrategy) Stream(_ context.Context, state *query.QueryState) (query.QueryStream, error) {
	s.streamCalls++
	s.lastRequest = state.Request
	if s.streamErr != nil {
		return nil, s.streamErr
	}
	out := make(chan query.StreamEvent, len(s.streamEvents))
	go func() {
		defer close(out)
		for _, event := range s.streamEvents {
			out <- event
		}
	}()
	return out, nil
}

type fakeAgentGetter struct {
	agent *AgentRuntimeDTO
}

func (g fakeAgentGetter) GetAgentForChat(context.Context, uint64) (*AgentRuntimeDTO, error) {
	return g.agent, nil
}

type fakeSessionRepo struct {
	nextID   uint64
	sessions map[uint64]*chatdomain.ChatSession
}

func (r *fakeSessionRepo) ensure() {
	if r.sessions == nil {
		r.sessions = map[uint64]*chatdomain.ChatSession{}
	}
	if r.nextID == 0 {
		r.nextID = 1
	}
}

func (r *fakeSessionRepo) Create(_ context.Context, session *chatdomain.ChatSession) error {
	r.ensure()
	session.ID = r.nextID
	r.nextID++
	r.sessions[session.ID] = session
	return nil
}

func (r *fakeSessionRepo) Update(context.Context, *chatdomain.ChatSession) error {
	return nil
}

func (r *fakeSessionRepo) FindByID(_ context.Context, id uint64) (*chatdomain.ChatSession, error) {
	r.ensure()
	return r.sessions[id], nil
}

func (r *fakeSessionRepo) FindPage(context.Context, *chatdomain.PageRequest) (*chatdomain.PageResult[*chatdomain.ChatSession], error) {
	return nil, nil
}

func (r *fakeSessionRepo) DeleteByID(context.Context, uint64) error {
	return nil
}

func (r *fakeSessionRepo) UpdateSummary(_ context.Context, id uint64, summary string, summarizedUntilMessageID uint64) error {
	r.ensure()
	if session := r.sessions[id]; session != nil {
		session.Summary = summary
		session.SummarizedUntilMessageID = summarizedUntilMessageID
	}
	return nil
}

type fakeMessageRepo struct {
	nextID   uint64
	messages []*chatdomain.ChatMessage
}

func (r *fakeMessageRepo) Create(_ context.Context, message *chatdomain.ChatMessage) error {
	if r.nextID == 0 {
		r.nextID = 1
	}
	message.ID = r.nextID
	r.nextID++
	r.messages = append(r.messages, message)
	return nil
}

func (r *fakeMessageRepo) FindBySessionID(ctx context.Context, sessionID uint64, limit int) ([]*chatdomain.ChatMessage, error) {
	return r.find(sessionID, 0, 0, limit), nil
}

func (r *fakeMessageRepo) FindBySessionIDAfterID(ctx context.Context, sessionID uint64, afterID uint64, limit int) ([]*chatdomain.ChatMessage, error) {
	return r.find(sessionID, afterID, 0, limit), nil
}

func (r *fakeMessageRepo) FindBySessionIDBeforeID(ctx context.Context, sessionID uint64, beforeID uint64, limit int) ([]*chatdomain.ChatMessage, error) {
	return r.find(sessionID, 0, beforeID, limit), nil
}

func (r *fakeMessageRepo) CountBySessionID(_ context.Context, sessionID uint64) (int64, error) {
	return int64(len(r.find(sessionID, 0, 0, 0))), nil
}

func (r *fakeMessageRepo) find(sessionID uint64, afterID uint64, beforeID uint64, limit int) []*chatdomain.ChatMessage {
	result := make([]*chatdomain.ChatMessage, 0)
	for _, message := range r.messages {
		if message.SessionID != sessionID {
			continue
		}
		if afterID > 0 && message.ID <= afterID {
			continue
		}
		if beforeID > 0 && message.ID >= beforeID {
			continue
		}
		result = append(result, message)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

func (r *fakeMessageRepo) countRole(role chatdomain.MessageRole) int {
	count := 0
	for _, message := range r.messages {
		if message.Role == role {
			count++
		}
	}
	return count
}
