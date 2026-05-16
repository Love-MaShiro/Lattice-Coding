package persistence

import "lattice-coding/internal/modules/workflow/domain"

func workflowToPO(workflow *domain.WorkflowDefinition, po *WorkflowPO) {
	po.ID = workflow.ID
	po.Name = workflow.Name
	po.Description = workflow.Description
	po.Status = string(workflow.Status)
	po.Version = workflow.Version
	po.Meta = workflow.Meta
}

func poToWorkflow(po *WorkflowPO) *domain.WorkflowDefinition {
	return &domain.WorkflowDefinition{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		Status:      domain.WorkflowStatus(po.Status),
		Version:     po.Version,
		Meta:        po.Meta,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

func nodeToPO(workflowID uint64, node *domain.NodeDefinition, configRaw string, po *WorkflowNodePO) {
	po.ID = node.ID
	po.WorkflowID = workflowID
	po.NodeKey = node.NodeKey
	po.Name = node.Name
	po.Type = string(node.Type)
	po.Config = configRaw
	po.Position = node.Position
	po.SortOrder = node.SortOrder
	po.Meta = node.Meta
}

func poToNode(po *WorkflowNodePO, parser domain.NodeConfigParser) (*domain.NodeDefinition, error) {
	config, err := parser.Parse(domain.NodeType(po.Type), po.Config)
	if err != nil {
		return nil, err
	}
	return &domain.NodeDefinition{
		ID:         po.ID,
		WorkflowID: po.WorkflowID,
		NodeKey:    po.NodeKey,
		Name:       po.Name,
		Type:       domain.NodeType(po.Type),
		Config:     config,
		Position:   po.Position,
		SortOrder:  po.SortOrder,
		Meta:       po.Meta,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
	}, nil
}

func edgeToPO(workflowID uint64, edge *domain.EdgeDefinition, po *WorkflowEdgePO) {
	po.ID = edge.ID
	po.WorkflowID = workflowID
	po.EdgeKey = edge.EdgeKey
	po.SourceKey = edge.SourceKey
	po.TargetKey = edge.TargetKey
	po.Condition = edge.Condition
	po.SortOrder = edge.SortOrder
	po.Meta = edge.Meta
}

func poToEdge(po *WorkflowEdgePO) *domain.EdgeDefinition {
	return &domain.EdgeDefinition{
		ID:         po.ID,
		WorkflowID: po.WorkflowID,
		EdgeKey:    po.EdgeKey,
		SourceKey:  po.SourceKey,
		TargetKey:  po.TargetKey,
		Condition:  po.Condition,
		SortOrder:  po.SortOrder,
		Meta:       po.Meta,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
	}
}
