package llm

import "sync"

type ClientRegistry struct {
	mu      sync.RWMutex
	clients map[string]LLMClient
}

func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients: make(map[string]LLMClient),
	}
}

func (r *ClientRegistry) Register(name string, client LLMClient) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[name] = client
}

func (r *ClientRegistry) Get(name string) (LLMClient, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	client, ok := r.clients[name]
	return client, ok
}

func (r *ClientRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}
	return names
}
