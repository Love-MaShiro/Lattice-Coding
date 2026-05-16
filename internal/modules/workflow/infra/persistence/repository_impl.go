package persistence

import (
	"context"

	commondb "lattice-coding/internal/common/db"
	"lattice-coding/internal/modules/workflow/domain"

	"gorm.io/gorm"
)

type WorkflowRepositoryImpl struct {
	db     *gorm.DB
	parser domain.NodeConfigParser
}

func NewWorkflowRepositoryImpl(db *gorm.DB, parser domain.NodeConfigParser) domain.WorkflowRepository {
	if parser == nil {
		parser = domain.NewJSONNodeConfigParser()
	}
	return &WorkflowRepositoryImpl{db: db, parser: parser}
}

func (r *WorkflowRepositoryImpl) CreateWithGraph(ctx context.Context, workflow *domain.WorkflowDefinition) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		po := &WorkflowPO{}
		workflowToPO(workflow, po)
		if err := tx.Create(po).Error; err != nil {
			return err
		}
		workflow.ID = po.ID
		workflow.CreatedAt = po.CreatedAt
		workflow.UpdatedAt = po.UpdatedAt
		if err := r.createNodes(ctx, tx, workflow.ID, workflow.Nodes); err != nil {
			return err
		}
		if err := r.createEdges(ctx, tx, workflow.ID, workflow.Edges); err != nil {
			return err
		}
		return nil
	})
}

func (r *WorkflowRepositoryImpl) UpdateWithGraph(ctx context.Context, workflow *domain.WorkflowDefinition) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		var current WorkflowPO
		if err := tx.First(&current, workflow.ID).Error; err != nil {
			return err
		}
		po := &WorkflowPO{}
		workflowToPO(workflow, po)
		if err := tx.Model(&WorkflowPO{}).
			Where("id = ?", workflow.ID).
			Updates(map[string]interface{}{
				"name":        po.Name,
				"description": po.Description,
				"status":      po.Status,
				"version":     po.Version,
				"meta":        po.Meta,
			}).Error; err != nil {
			return err
		}
		if err := tx.Where("workflow_id = ?", workflow.ID).Delete(&WorkflowNodePO{}).Error; err != nil {
			return err
		}
		if err := tx.Where("workflow_id = ?", workflow.ID).Delete(&WorkflowEdgePO{}).Error; err != nil {
			return err
		}
		if err := r.createNodes(ctx, tx, workflow.ID, workflow.Nodes); err != nil {
			return err
		}
		if err := r.createEdges(ctx, tx, workflow.ID, workflow.Edges); err != nil {
			return err
		}
		return nil
	})
}

func (r *WorkflowRepositoryImpl) FindByIDWithGraph(ctx context.Context, id uint64) (*domain.WorkflowDefinition, error) {
	var workflowPO WorkflowPO
	if err := r.db.WithContext(ctx).First(&workflowPO, id).Error; err != nil {
		return nil, err
	}
	workflow := poToWorkflow(&workflowPO)

	var nodePOs []WorkflowNodePO
	if err := r.db.WithContext(ctx).
		Where("workflow_id = ?", id).
		Order("sort_order ASC, id ASC").
		Find(&nodePOs).Error; err != nil {
		return nil, err
	}
	workflow.Nodes = make([]domain.NodeDefinition, 0, len(nodePOs))
	for i := range nodePOs {
		node, err := poToNode(&nodePOs[i], r.parser)
		if err != nil {
			return nil, err
		}
		workflow.Nodes = append(workflow.Nodes, *node)
	}

	var edgePOs []WorkflowEdgePO
	if err := r.db.WithContext(ctx).
		Where("workflow_id = ?", id).
		Order("sort_order ASC, id ASC").
		Find(&edgePOs).Error; err != nil {
		return nil, err
	}
	workflow.Edges = make([]domain.EdgeDefinition, 0, len(edgePOs))
	for i := range edgePOs {
		workflow.Edges = append(workflow.Edges, *poToEdge(&edgePOs[i]))
	}

	return workflow, nil
}

func (r *WorkflowRepositoryImpl) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.WorkflowDefinition], error) {
	query := r.db.WithContext(ctx).Model(&WorkflowPO{})
	if req.Status != "" {
		query = query.Where("status = ?", string(req.Status))
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var pos []WorkflowPO
	if err := query.Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("id DESC").
		Find(&pos).Error; err != nil {
		return nil, err
	}

	items := make([]*domain.WorkflowDefinition, len(pos))
	for i := range pos {
		items[i] = poToWorkflow(&pos[i])
	}
	return &domain.PageResult[*domain.WorkflowDefinition]{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (r *WorkflowRepositoryImpl) DeleteWithGraph(ctx context.Context, id uint64) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		var current WorkflowPO
		if err := tx.First(&current, id).Error; err != nil {
			return err
		}
		if err := tx.Where("workflow_id = ?", id).Delete(&WorkflowNodePO{}).Error; err != nil {
			return err
		}
		if err := tx.Where("workflow_id = ?", id).Delete(&WorkflowEdgePO{}).Error; err != nil {
			return err
		}
		return tx.Delete(&WorkflowPO{}, id).Error
	})
}

func (r *WorkflowRepositoryImpl) createNodes(ctx context.Context, tx *gorm.DB, workflowID uint64, nodes []domain.NodeDefinition) error {
	if len(nodes) == 0 {
		return nil
	}
	pos := make([]WorkflowNodePO, len(nodes))
	for i := range nodes {
		configRaw, err := r.parser.Marshal(nodes[i].Config)
		if err != nil {
			return err
		}
		nodeToPO(workflowID, &nodes[i], configRaw, &pos[i])
	}
	if err := tx.WithContext(ctx).Create(&pos).Error; err != nil {
		return err
	}
	for i := range pos {
		nodes[i].ID = pos[i].ID
		nodes[i].WorkflowID = workflowID
		nodes[i].CreatedAt = pos[i].CreatedAt
		nodes[i].UpdatedAt = pos[i].UpdatedAt
	}
	return nil
}

func (r *WorkflowRepositoryImpl) createEdges(ctx context.Context, tx *gorm.DB, workflowID uint64, edges []domain.EdgeDefinition) error {
	if len(edges) == 0 {
		return nil
	}
	pos := make([]WorkflowEdgePO, len(edges))
	for i := range edges {
		edgeToPO(workflowID, &edges[i], &pos[i])
	}
	if err := tx.WithContext(ctx).Create(&pos).Error; err != nil {
		return err
	}
	for i := range pos {
		edges[i].ID = pos[i].ID
		edges[i].WorkflowID = workflowID
		edges[i].CreatedAt = pos[i].CreatedAt
		edges[i].UpdatedAt = pos[i].UpdatedAt
	}
	return nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&WorkflowPO{}, &WorkflowNodePO{}, &WorkflowEdgePO{})
}
