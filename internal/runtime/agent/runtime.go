package agent

import (
	"context"
	"errors"

	"lattice-coding/internal/runtime/llm"
	"lattice-coding/internal/runtime/tool"
)

type AgentRuntime struct {
	strategies      map[string]ExecutionStrategy
	defaultStrategy string
}

func NewAgentRuntime(llmExecutor *llm.Executor, toolExecutor *tool.Executor) *AgentRuntime {
	runtime := &AgentRuntime{
		strategies:      map[string]ExecutionStrategy{},
		defaultStrategy: StrategyReAct,
	}
	runtime.RegisterStrategy(NewReActStrategy(llmExecutor, toolExecutor, nil))
	return runtime
}

func (r *AgentRuntime) RegisterStrategy(strategy ExecutionStrategy) {
	if strategy == nil || strategy.Name() == "" {
		return
	}
	r.strategies[strategy.Name()] = strategy
}

func (r *AgentRuntime) Strategy(name string) (ExecutionStrategy, bool) {
	strategy, ok := r.strategies[name]
	return strategy, ok
}

func (r *AgentRuntime) Run(ctx context.Context, req Request) (*Result, error) {
	if r == nil {
		return nil, errors.New("agent runtime is nil")
	}
	name := req.Strategy
	if name == "" {
		name = r.defaultStrategy
	}
	strategy, ok := r.Strategy(name)
	if !ok {
		return nil, errors.New("agent strategy not registered: " + name)
	}
	return strategy.Execute(ctx, req)
}
