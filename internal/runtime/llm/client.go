package llm

import "context"

type LLMClient interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	Stream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)
	Close() error
}

type ClientRegistry struct {
	clients map[string]LLMClient
}

func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients: make(map[string]LLMClient),
	}
}

func (r *ClientRegistry) Register(name string, client LLMClient) {
	r.clients[name] = client
}

func (r *ClientRegistry) Get(name string) (LLMClient, bool) {
	client, ok := r.clients[name]
	return client, ok
}

func (r *ClientRegistry) List() []string {
	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}
	return names
}
