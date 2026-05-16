package application

import "lattice-coding/internal/modules/workflow/domain"

func ToWorkflowDTO(workflow *domain.WorkflowDefinition) *WorkflowDTO {
	if workflow == nil {
		return nil
	}
	return &WorkflowDTO{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Status:      workflow.Status,
		Version:     workflow.Version,
		Meta:        workflow.Meta,
		Nodes:       ToNodeDTOs(workflow.Nodes),
		Edges:       ToEdgeDTOs(workflow.Edges),
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	}
}

func ToWorkflowSummaryDTO(workflow *domain.WorkflowDefinition) *WorkflowDTO {
	dto := ToWorkflowDTO(workflow)
	if dto != nil {
		dto.Nodes = nil
		dto.Edges = nil
	}
	return dto
}

func ToNodeDTOs(nodes []domain.NodeDefinition) []NodeDTO {
	items := make([]NodeDTO, len(nodes))
	for i := range nodes {
		items[i] = NodeDTO{
			ID:         nodes[i].ID,
			WorkflowID: nodes[i].WorkflowID,
			NodeKey:    nodes[i].NodeKey,
			Name:       nodes[i].Name,
			Type:       nodes[i].Type,
			Config:     nodes[i].Config,
			Position:   nodes[i].Position,
			SortOrder:  nodes[i].SortOrder,
			Meta:       nodes[i].Meta,
			CreatedAt:  nodes[i].CreatedAt,
			UpdatedAt:  nodes[i].UpdatedAt,
		}
	}
	return items
}

func ToEdgeDTOs(edges []domain.EdgeDefinition) []EdgeDTO {
	items := make([]EdgeDTO, len(edges))
	for i := range edges {
		items[i] = EdgeDTO{
			ID:         edges[i].ID,
			WorkflowID: edges[i].WorkflowID,
			EdgeKey:    edges[i].EdgeKey,
			SourceKey:  edges[i].SourceKey,
			TargetKey:  edges[i].TargetKey,
			Condition:  edges[i].Condition,
			SortOrder:  edges[i].SortOrder,
			Meta:       edges[i].Meta,
			CreatedAt:  edges[i].CreatedAt,
			UpdatedAt:  edges[i].UpdatedAt,
		}
	}
	return items
}
