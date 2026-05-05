package llm

import (
	"context"
	"strings"
)

type RouterConfig struct {
	DefaultPrimary string
	FallbackList   []string
	MaxAttempts    int
	EnableFallback bool
}

type Router struct {
	registry *ClientRegistry
	config   RouterConfig
}

func NewRouter(registry *ClientRegistry, cfg RouterConfig) *Router {
	return &Router{
		registry: registry,
		config:   cfg,
	}
}

func (r *Router) Route(ctx context.Context, req ChatRequest) (*ChatResponse, CallResult) {
	result := CallResult{
		Provider: req.Provider,
		Model:    req.Model,
	}

	providers := r.buildProviderList(req.Provider)

	for i, p := range providers {
		client, ok := r.registry.Get(p)
		if !ok {
			continue
		}

		resp, err := client.Chat(ctx, req)
		if err != nil {
			result.Error = err
			result.Provider = p
			if r.shouldFallback(i) {
				result.Fallback = true
				continue
			}
			result.Success = false
			return nil, result
		}

		result.Provider = p
		result.Success = true
		result.LatencyMs = resp.LatencyMs
		result.Tokens = 0
		if resp.Usage != nil {
			result.Tokens = resp.Usage.TotalTokens
		}
		return resp, result
	}

	result.Success = false
	result.Error = ErrAllFallbackFail
	return nil, result
}

func (r *Router) RouteStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, CallResult) {
	result := CallResult{
		Provider: req.Provider,
		Model:    req.Model,
	}

	providers := r.buildProviderList(req.Provider)

	for i, p := range providers {
		client, ok := r.registry.Get(p)
		if !ok {
			continue
		}

		ch, err := client.Stream(ctx, req)
		if err != nil {
			result.Error = err
			result.Provider = p
			if r.shouldFallback(i) {
				result.Fallback = true
				continue
			}
			result.Success = false
			return nil, result
		}

		result.Provider = p
		result.Success = true
		return ch, result
	}

	result.Success = false
	result.Error = ErrAllFallbackFail
	return nil, result
}

func (r *Router) buildProviderList(provider string) []string {
	if provider != "" {
		return []string{provider}
	}

	list := []string{r.config.DefaultPrimary}
	if r.config.EnableFallback {
		list = append(list, r.config.FallbackList...)
	}
	return list
}

func (r *Router) shouldFallback(idx int) bool {
	return r.config.EnableFallback && idx < r.config.MaxAttempts
}

func parseProviderModel(p string) (provider, model string) {
	parts := strings.SplitN(p, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return p, ""
}
