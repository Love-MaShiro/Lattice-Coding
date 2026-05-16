package strategy

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/runtime/agent"
	"lattice-coding/internal/runtime/llm"
	"lattice-coding/internal/runtime/query"
	"lattice-coding/internal/runtime/tool"
)

func TestQueryEnginePureReActToolCallChain(t *testing.T) {
	ctx := context.Background()

	llmExecutor := llm.NewExecutor(&config.LLMConfig{
		Pool:   config.PoolConfig{MaxConcurrent: 1},
		Stream: config.PoolConfig{MaxConcurrent: 1},
		Routing: config.RoutingConfig{
			Default: config.RouteConfig{Primary: "fake"},
		},
	})
	fakeLLM := &reactChainLLMClient{
		responses: []string{
			`{"type":"tool_call","reason":"echo input","tool":"demo.echo","args":{"text":"hello"}}`,
			`{"type":"final","answer":"tool returned echo: hello"}`,
		},
	}
	llmExecutor.RegisterClient("fake", fakeLLM)

	toolRegistry := tool.NewRegistry()
	if err := toolRegistry.Register(echoTool{}); err != nil {
		t.Fatalf("register echo tool: %v", err)
	}
	toolExecutor := tool.NewExecutor(toolRegistry, tool.WithSafetyChecker(tool.NoopSafetyChecker{}))
	agentRuntime := agent.NewAgentRuntime(llmExecutor, toolExecutor)

	engine := query.NewEngine(
		query.WithStrategy(NewPureReActStrategy(agentRuntime)),
	)

	result, err := engine.Run(ctx, query.QueryRequest{
		RunID:        "run-1",
		UserID:       "user-1",
		Input:        "call echo tool with hello, then answer",
		Mode:         query.ExecutionModePureReAct,
		Provider:     "fake",
		AllowedTools: []string{"demo.echo"},
		Budget: query.QueryBudget{
			MaxSteps:     4,
			MaxToolCalls: 2,
		},
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("query run returned error: %v", err)
	}
	if result.Content != "tool returned echo: hello" {
		t.Fatalf("content = %q", result.Content)
	}
	if result.Mode != query.ExecutionModePureReAct {
		t.Fatalf("mode = %q", result.Mode)
	}
	if len(result.Steps) != 2 {
		t.Fatalf("steps length = %d, want 2", len(result.Steps))
	}
	if result.Steps[0].Name != "demo.echo" {
		t.Fatalf("first step name = %q, want demo.echo", result.Steps[0].Name)
	}
	if result.Steps[0].Content != "echo: hello" {
		t.Fatalf("first step content = %q, want echo: hello", result.Steps[0].Content)
	}
	if result.Steps[0].IsError {
		t.Fatal("first step should not be an error")
	}
	if result.Steps[1].Name != agent.ReActActionFinal {
		t.Fatalf("second step name = %q, want final", result.Steps[1].Name)
	}
	if len(result.Messages) != 2 {
		t.Fatalf("messages length = %d, want 2", len(result.Messages))
	}
	if got := fakeLLM.callCount(); got != 2 {
		t.Fatalf("llm calls = %d, want 2", got)
	}
	secondPrompt := fakeLLM.userPrompt(1)
	if !strings.Contains(secondPrompt, "Observation: echo: hello") {
		t.Fatalf("second prompt did not include tool observation:\n%s", secondPrompt)
	}
}

func TestQueryEnginePureReActRejectsDisallowedTool(t *testing.T) {
	ctx := context.Background()

	llmExecutor := llm.NewExecutor(&config.LLMConfig{
		Pool:   config.PoolConfig{MaxConcurrent: 1},
		Stream: config.PoolConfig{MaxConcurrent: 1},
		Routing: config.RoutingConfig{
			Default: config.RouteConfig{Primary: "fake"},
		},
	})
	fakeLLM := &reactChainLLMClient{
		responses: []string{
			`{"type":"tool_call","reason":"try blocked","tool":"demo.echo","args":{"text":"hello"}}`,
			`{"type":"final","answer":"blocked tool was observed"}`,
		},
	}
	llmExecutor.RegisterClient("fake", fakeLLM)

	toolRegistry := tool.NewRegistry()
	if err := toolRegistry.Register(echoTool{}); err != nil {
		t.Fatalf("register echo tool: %v", err)
	}
	agentRuntime := agent.NewAgentRuntime(llmExecutor, tool.NewExecutor(toolRegistry))
	engine := query.NewEngine(query.WithStrategy(NewPureReActStrategy(agentRuntime)))

	result, err := engine.Run(ctx, query.QueryRequest{
		RunID:        "run-2",
		Input:        "try a blocked tool",
		Mode:         query.ExecutionModePureReAct,
		Provider:     "fake",
		AllowedTools: []string{"other.tool"},
		Budget:       query.QueryBudget{MaxSteps: 4, MaxToolCalls: 2},
		Timeout:      time.Second,
	})
	if err != nil {
		t.Fatalf("query run returned error: %v", err)
	}
	if len(result.Steps) != 2 {
		t.Fatalf("steps length = %d, want 2", len(result.Steps))
	}
	if !result.Steps[0].IsError {
		t.Fatal("disallowed tool step should be marked as error")
	}
	if !strings.Contains(result.Steps[0].Content, "tool is not allowed") {
		t.Fatalf("unexpected disallowed observation: %q", result.Steps[0].Content)
	}
}

type reactChainLLMClient struct {
	mu        sync.Mutex
	responses []string
	requests  []llm.ChatRequest
}

func (c *reactChainLLMClient) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requests = append(c.requests, req)
	if len(c.requests) > len(c.responses) {
		return nil, errors.New("unexpected llm call")
	}
	return &llm.ChatResponse{Content: c.responses[len(c.requests)-1]}, nil
}

func (c *reactChainLLMClient) Stream(context.Context, llm.ChatRequest) (<-chan llm.StreamChunk, error) {
	return nil, errors.New("stream not implemented")
}

func (c *reactChainLLMClient) Close() error {
	return nil
}

func (c *reactChainLLMClient) callCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.requests)
}

func (c *reactChainLLMClient) userPrompt(index int) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if index >= len(c.requests) {
		return ""
	}
	for _, message := range c.requests[index].Messages {
		if message.Role == "user" {
			return message.Content
		}
	}
	return ""
}

type echoTool struct {
	tool.BaseTool
}

func (echoTool) Name() string {
	return "demo.echo"
}

func (echoTool) Description() string {
	return "returns the text argument"
}

func (echoTool) Schema() tool.Schema {
	return tool.ObjectSchema(map[string]interface{}{
		"text": tool.StringSchema("text to echo"),
	}, "text")
}

func (echoTool) IsReadOnly() bool {
	return true
}

func (echoTool) IsConcurrencySafe() bool {
	return true
}

func (echoTool) CheckPermission(context.Context, tool.ToolRequest) (tool.PermissionDecision, string, error) {
	return tool.PermissionAllow, "allowed in test", nil
}

func (echoTool) Execute(_ context.Context, req tool.ToolRequest) (tool.ToolOutput, error) {
	text, _ := req.Input["text"].(string)
	return tool.ToolOutput{
		Content: "echo: " + text,
		Data: map[string]interface{}{
			"text": text,
		},
	}, nil
}
