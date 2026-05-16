package application

import (
	"context"
	"strings"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/workflow/domain"
)

type CommandService struct {
	repo   domain.WorkflowRepository
	parser domain.NodeConfigParser
}

func NewCommandService(repo domain.WorkflowRepository, parser domain.NodeConfigParser) *CommandService {
	if parser == nil {
		parser = domain.NewJSONNodeConfigParser()
	}
	return &CommandService{repo: repo, parser: parser}
}

func (s *CommandService) CreateWorkflow(ctx context.Context, cmd SaveWorkflowCommand) (*WorkflowDTO, error) {
	workflow, err := s.buildWorkflow(cmd)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateWithGraph(ctx, workflow); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "create workflow failed")
	}
	return ToWorkflowDTO(workflow), nil
}

func (s *CommandService) UpdateWorkflow(ctx context.Context, id uint64, cmd SaveWorkflowCommand) (*WorkflowDTO, error) {
	if id == 0 {
		return nil, errors.InvalidArg("workflow id is required")
	}
	workflow, err := s.buildWorkflow(cmd)
	if err != nil {
		return nil, err
	}
	workflow.ID = id
	if err := s.repo.UpdateWithGraph(ctx, workflow); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "update workflow failed")
	}
	result, err := s.repo.FindByIDWithGraph(ctx, id)
	if err != nil {
		return nil, errors.NotFoundErr("workflow not found")
	}
	return ToWorkflowDTO(result), nil
}

func (s *CommandService) DeleteWorkflow(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.InvalidArg("workflow id is required")
	}
	if err := s.repo.DeleteWithGraph(ctx, id); err != nil {
		return errors.DatabaseErrWithErr(err, "delete workflow failed")
	}
	return nil
}

func (s *CommandService) buildWorkflow(cmd SaveWorkflowCommand) (*domain.WorkflowDefinition, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, errors.InvalidArg("workflow name is required")
	}
	if cmd.Status == "" {
		cmd.Status = domain.WorkflowStatusDraft
	}
	if cmd.Version <= 0 {
		cmd.Version = 1
	}
	if cmd.Meta == "" {
		cmd.Meta = "{}"
	}
	nodes, err := s.buildNodes(cmd.Nodes)
	if err != nil {
		return nil, err
	}
	edges, err := buildEdges(cmd.Edges)
	if err != nil {
		return nil, err
	}
	if err := validateGraph(nodes, edges); err != nil {
		return nil, err
	}
	return &domain.WorkflowDefinition{
		Name:        strings.TrimSpace(cmd.Name),
		Description: cmd.Description,
		Status:      cmd.Status,
		Version:     cmd.Version,
		Meta:        cmd.Meta,
		Nodes:       nodes,
		Edges:       edges,
	}, nil
}

func (s *CommandService) buildNodes(commands []NodeCommand) ([]domain.NodeDefinition, error) {
	nodes := make([]domain.NodeDefinition, 0, len(commands))
	for i, cmd := range commands {
		if strings.TrimSpace(cmd.NodeKey) == "" {
			return nil, errors.InvalidArg("node_key is required")
		}
		if cmd.Type == "" {
			return nil, errors.InvalidArg("node type is required")
		}
		config, err := s.parser.Parse(cmd.Type, cmd.ConfigRaw)
		if err != nil {
			return nil, errors.InvalidArg("invalid node config: " + err.Error())
		}
		if cmd.Meta == "" {
			cmd.Meta = "{}"
		}
		if cmd.Position == "" {
			cmd.Position = "{}"
		}
		nodes = append(nodes, domain.NodeDefinition{
			NodeKey:   strings.TrimSpace(cmd.NodeKey),
			Name:      cmd.Name,
			Type:      cmd.Type,
			Config:    config,
			Position:  cmd.Position,
			SortOrder: firstNonZero(cmd.SortOrder, i+1),
			Meta:      cmd.Meta,
		})
	}
	return nodes, nil
}

func buildEdges(commands []EdgeCommand) ([]domain.EdgeDefinition, error) {
	edges := make([]domain.EdgeDefinition, 0, len(commands))
	for i, cmd := range commands {
		if strings.TrimSpace(cmd.SourceKey) == "" || strings.TrimSpace(cmd.TargetKey) == "" {
			return nil, errors.InvalidArg("edge source_key and target_key are required")
		}
		if cmd.EdgeKey == "" {
			cmd.EdgeKey = strings.TrimSpace(cmd.SourceKey) + "->" + strings.TrimSpace(cmd.TargetKey)
		}
		if cmd.Meta == "" {
			cmd.Meta = "{}"
		}
		edges = append(edges, domain.EdgeDefinition{
			EdgeKey:   strings.TrimSpace(cmd.EdgeKey),
			SourceKey: strings.TrimSpace(cmd.SourceKey),
			TargetKey: strings.TrimSpace(cmd.TargetKey),
			Condition: cmd.Condition,
			SortOrder: firstNonZero(cmd.SortOrder, i+1),
			Meta:      cmd.Meta,
		})
	}
	return edges, nil
}

func validateGraph(nodes []domain.NodeDefinition, edges []domain.EdgeDefinition) error {
	seenNodes := map[string]struct{}{}
	for _, node := range nodes {
		if _, ok := seenNodes[node.NodeKey]; ok {
			return errors.InvalidArg("duplicated node_key: " + node.NodeKey)
		}
		seenNodes[node.NodeKey] = struct{}{}
	}
	seenEdges := map[string]struct{}{}
	for _, edge := range edges {
		if _, ok := seenNodes[edge.SourceKey]; !ok {
			return errors.InvalidArg("edge source node not found: " + edge.SourceKey)
		}
		if _, ok := seenNodes[edge.TargetKey]; !ok {
			return errors.InvalidArg("edge target node not found: " + edge.TargetKey)
		}
		if _, ok := seenEdges[edge.EdgeKey]; ok {
			return errors.InvalidArg("duplicated edge_key: " + edge.EdgeKey)
		}
		seenEdges[edge.EdgeKey] = struct{}{}
	}
	return nil
}

func firstNonZero(value int, fallback int) int {
	if value != 0 {
		return value
	}
	return fallback
}
