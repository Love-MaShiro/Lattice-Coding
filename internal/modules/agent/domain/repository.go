package domain

import "context"

type AgentRepository interface {
	Create(ctx context.Context, agent *Agent) error
	Update(ctx context.Context, agent *Agent) error
	FindByID(ctx context.Context, id uint64) (*Agent, error)
	FindPage(ctx context.Context, req *PageRequest) (*PageResult[*Agent], error)
	DeleteByID(ctx context.Context, id uint64) error
	ExistsByName(ctx context.Context, name string, excludeID uint64) (bool, error)
	UpdateEnabled(ctx context.Context, id uint64, enabled bool) error
}

type AgentToolRepository interface {
	BatchCreate(ctx context.Context, agentTools []*AgentTool) error
	DeleteByAgentID(ctx context.Context, agentID uint64) error
	FindByAgentID(ctx context.Context, agentID uint64) ([]*AgentTool, error)
	CountByAgentIDs(ctx context.Context, agentIDs []uint64) (map[uint64]int64, error)
	DeleteByAgentIDs(ctx context.Context, agentIDs []uint64) error
}

type AgentReferenceChecker interface {
	HasModelConfigReferences(ctx context.Context, modelConfigID uint64) (bool, error)
}
