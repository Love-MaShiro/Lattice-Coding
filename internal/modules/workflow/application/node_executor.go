package application

import (
	"context"
	"fmt"
	"sync"

	"lattice-coding/internal/modules/workflow/domain"
)

type NodeExecutionInput struct {
	Workflow *domain.WorkflowDefinition
	Node     domain.NodeDefinition
	RunState *domain.RunState
	Data     map[string]interface{}
}

type NodeExecutionOutput struct {
	Data map[string]interface{}
}

type NodeExecutor interface {
	Type() domain.NodeType
	Execute(ctx context.Context, input NodeExecutionInput) (*NodeExecutionOutput, error)
}

type NodeRegistry struct {
	mu        sync.RWMutex
	executors map[domain.NodeType]NodeExecutor
}

func NewNodeRegistry(executors ...NodeExecutor) *NodeRegistry {
	registry := &NodeRegistry{executors: map[domain.NodeType]NodeExecutor{}}
	for _, executor := range executors {
		registry.Register(executor)
	}
	return registry
}

func (r *NodeRegistry) Register(executor NodeExecutor) {
	if executor == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executors[executor.Type()] = executor
}

func (r *NodeRegistry) Get(nodeType domain.NodeType) (NodeExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	executor, ok := r.executors[nodeType]
	return executor, ok
}

func (r *NodeRegistry) Execute(ctx context.Context, input NodeExecutionInput) (*NodeExecutionOutput, error) {
	executor, ok := r.Get(input.Node.Type)
	if !ok {
		return nil, fmt.Errorf("node executor not registered: %s", input.Node.Type)
	}
	return executor.Execute(ctx, input)
}
