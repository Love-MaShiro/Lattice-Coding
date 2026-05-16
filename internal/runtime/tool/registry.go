package tool

import (
	"sort"
	"sync"

	"lattice-coding/internal/common/errors"
)

type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *ToolRegistry {
	return &ToolRegistry{tools: map[string]Tool{}}
}

func NewToolRegistry() *ToolRegistry {
	return NewRegistry()
}

func (r *ToolRegistry) Register(t Tool) error {
	if t == nil {
		return errors.ToolInvalidParamsErr("tool is required")
	}
	name := t.Name()
	if name == "" {
		return errors.ToolInvalidParamsErr("tool name is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tools[name]; exists {
		return errors.ToolInvalidParamsErr("tool already registered: " + name)
	}
	r.tools[name] = t
	return nil
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

func (r *ToolRegistry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)

	tools := make([]Tool, 0, len(names))
	for _, name := range names {
		tools = append(tools, r.tools[name])
	}
	return tools
}

func (r *ToolRegistry) ListDescriptors() []ToolDescriptor {
	tools := r.List()
	descriptors := make([]ToolDescriptor, 0, len(tools))
	for _, t := range tools {
		descriptors = append(descriptors, DescriptorOf(t))
	}
	return descriptors
}

func DescriptorOf(t Tool) ToolDescriptor {
	if t == nil {
		return ToolDescriptor{}
	}
	return ToolDescriptor{
		Name:            t.Name(),
		Description:     t.Description(),
		Prompt:          t.Prompt(),
		Schema:          t.Schema(),
		ReadOnly:        t.IsReadOnly(),
		ConcurrencySafe: t.IsConcurrencySafe(),
		Destructive:     t.IsDestructive(),
	}
}
