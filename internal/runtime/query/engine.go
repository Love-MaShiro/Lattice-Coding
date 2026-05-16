package query

import (
	"context"
	"errors"
	"time"
)

type QueryStrategy interface {
	Mode() ExecutionMode
	Execute(ctx context.Context, state *QueryState) (*QueryResult, error)
}

type RunBinder interface {
	BindRun(ctx context.Context, req QueryRequest) (string, error)
}

type AgentConfigLoader interface {
	LoadAgentConfig(ctx context.Context, req QueryRequest) (*AgentConfig, error)
}

type ContextBuilder interface {
	BuildContext(ctx context.Context, req QueryRequest) (map[string]interface{}, error)
}

type RecoveryPolicy interface {
	Recover(ctx context.Context, state *QueryState, err error) (*QueryResult, error)
}

type QueryEngine struct {
	router       ExecutionModeRouter
	strategies   map[ExecutionMode]QueryStrategy
	runBinder    RunBinder
	agentLoader  AgentConfigLoader
	contextBuild ContextBuilder
	recovery     RecoveryPolicy
}

type Engine = QueryEngine

type EngineOption func(*QueryEngine)

func NewEngine(opts ...EngineOption) *QueryEngine {
	e := &QueryEngine{
		router:     NewExecutionModeRouter(),
		strategies: map[ExecutionMode]QueryStrategy{},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}
	return e
}

func WithRouter(router ExecutionModeRouter) EngineOption {
	return func(e *QueryEngine) {
		if router != nil {
			e.router = router
		}
	}
}

func WithStrategy(strategy QueryStrategy) EngineOption {
	return func(e *QueryEngine) {
		e.RegisterStrategy(strategy)
	}
}

func WithRunBinder(binder RunBinder) EngineOption {
	return func(e *QueryEngine) {
		e.runBinder = binder
	}
}

func WithAgentConfigLoader(loader AgentConfigLoader) EngineOption {
	return func(e *QueryEngine) {
		e.agentLoader = loader
	}
}

func WithContextBuilder(builder ContextBuilder) EngineOption {
	return func(e *QueryEngine) {
		e.contextBuild = builder
	}
}

func WithRecoveryPolicy(policy RecoveryPolicy) EngineOption {
	return func(e *QueryEngine) {
		e.recovery = policy
	}
}

func (e *QueryEngine) RegisterStrategy(strategy QueryStrategy) {
	if strategy == nil {
		return
	}
	e.strategies[strategy.Mode()] = strategy
}

func (e *QueryEngine) Run(ctx context.Context, req QueryRequest) (*QueryResult, error) {
	if e == nil {
		return nil, ErrQueryFailed.WithMessage("query engine is nil")
	}
	if req.Input == "" {
		return nil, ErrInvalidRequest.WithMessage("query input is required")
	}
	if req.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, req.Timeout)
		defer cancel()
	}

	if e.runBinder != nil {
		runID, err := e.runBinder.BindRun(ctx, req)
		if err != nil {
			return nil, errors.Join(ErrRunBindFailed, err)
		}
		if runID != "" {
			req.RunID = runID
		}
	}

	if e.agentLoader != nil {
		cfg, err := e.agentLoader.LoadAgentConfig(ctx, req)
		if err != nil {
			return nil, err
		}
		applyAgentConfig(&req, cfg)
	}

	if e.contextBuild != nil {
		metadata, err := e.contextBuild.BuildContext(ctx, req)
		if err != nil {
			return nil, errors.Join(ErrContextBuildFailed, err)
		}
		req.Metadata = mergeMetadata(req.Metadata, metadata)
	}

	req.Mode = e.router.Route(req)
	state := NewState(req)
	strategy, ok := e.strategies[req.Mode]
	if !ok {
		return nil, ErrStrategyNotFound.WithMessage("query strategy not found: " + req.Mode.String())
	}

	result, err := strategy.Execute(ctx, state)
	state.Finish()
	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			err = errors.Join(ErrQueryInterrupted, err)
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = errors.Join(ErrQueryTimeout, err)
		}
		if e.recovery != nil {
			if recovered, recoverErr := e.recovery.Recover(ctx, state, err); recoverErr == nil && recovered != nil {
				return recovered, nil
			}
		}
		return nil, err
	}
	if result == nil {
		return nil, ErrQueryFailed.WithMessage("query strategy returned nil result")
	}
	result.RunID = state.RunID
	result.Mode = req.Mode
	if result.Metadata == nil {
		result.Metadata = map[string]interface{}{}
	}
	result.Metadata["started_at"] = state.StartedAt
	result.Metadata["completed_at"] = time.Now()
	result.Metadata["budget"] = state.Budget.Snapshot()
	return result, nil
}

func (e *QueryEngine) Stream(ctx context.Context, req QueryRequest) (QueryStream, error) {
	out := make(chan StreamResult, 8)
	go func() {
		defer close(out)
		out <- StreamResult{Type: StreamEventStarted, RunID: req.RunID}
		result, err := e.Run(ctx, req)
		if err != nil {
			out <- StreamResult{Type: StreamEventError, RunID: req.RunID, Err: err, Done: true}
			return
		}
		out <- StreamResult{Type: StreamEventDone, RunID: result.RunID, Content: result.Content, Done: true, Metadata: result.Metadata}
	}()
	return out, nil
}

func applyAgentConfig(req *QueryRequest, cfg *AgentConfig) {
	if cfg == nil {
		return
	}
	if req.Mode == "" {
		req.Mode = cfg.Mode
	}
	if req.Provider == "" {
		req.Provider = cfg.Provider
	}
	if req.Model == "" {
		req.Model = cfg.Model
	}
	if req.ModelConfigID == 0 {
		req.ModelConfigID = cfg.ModelConfigID
	}
	if len(req.AllowedTools) == 0 {
		req.AllowedTools = cfg.AllowedTools
	}
	if req.Budget.MaxSteps == 0 {
		req.Budget.MaxSteps = cfg.Budget.MaxSteps
	}
	if req.Budget.MaxTokens == 0 {
		req.Budget.MaxTokens = cfg.Budget.MaxTokens
	}
	if req.Budget.MaxToolCalls == 0 {
		req.Budget.MaxToolCalls = cfg.Budget.MaxToolCalls
	}
}

func mergeMetadata(base map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
	if base == nil {
		base = map[string]interface{}{}
	}
	for key, value := range extra {
		base[key] = value
	}
	return base
}

type AgentConfig struct {
	Mode          ExecutionMode
	Provider      string
	Model         string
	ModelConfigID uint64
	AllowedTools  []string
	Budget        QueryBudget
	Metadata      map[string]interface{}
}
