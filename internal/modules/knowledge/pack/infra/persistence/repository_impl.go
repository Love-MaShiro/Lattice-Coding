package persistence

import (
	"context"

	commondb "lattice-coding/internal/common/db"
	packdomain "lattice-coding/internal/modules/knowledge/pack/domain"

	"gorm.io/gorm"
)

type PackRepositoryImpl struct {
	db *gorm.DB
}

func NewPackRepositoryImpl(db *gorm.DB) packdomain.PackRepository {
	return &PackRepositoryImpl{db: db}
}

func (r *PackRepositoryImpl) CreateWithItems(ctx context.Context, pack *packdomain.KnowledgePack) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		po := &KnowledgePackPO{}
		packToPO(pack, po)
		if err := tx.Create(po).Error; err != nil {
			return err
		}
		pack.ID = po.ID
		pack.CreatedAt = po.CreatedAt
		pack.UpdatedAt = po.UpdatedAt
		if err := createItems(tx, pack.ID, pack.Items); err != nil {
			return err
		}
		return createCitations(tx, pack.ID, pack.Citations)
	})
}

func (r *PackRepositoryImpl) FindByIDWithItems(ctx context.Context, id uint64) (*packdomain.KnowledgePack, error) {
	var po KnowledgePackPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, err
	}
	return r.loadPackGraph(ctx, &po)
}

func (r *PackRepositoryImpl) FindByKeyWithItems(ctx context.Context, packKey string) (*packdomain.KnowledgePack, error) {
	var po KnowledgePackPO
	if err := r.db.WithContext(ctx).Where("pack_key = ?", packKey).First(&po).Error; err != nil {
		return nil, err
	}
	return r.loadPackGraph(ctx, &po)
}

func (r *PackRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		if err := tx.Where("pack_id = ?", id).Delete(&KnowledgeItemPO{}).Error; err != nil {
			return err
		}
		if err := tx.Where("pack_id = ?", id).Delete(&KnowledgeCitationPO{}).Error; err != nil {
			return err
		}
		return tx.Delete(&KnowledgePackPO{}, id).Error
	})
}

func (r *PackRepositoryImpl) loadPackGraph(ctx context.Context, po *KnowledgePackPO) (*packdomain.KnowledgePack, error) {
	pack := poToPack(po)
	var itemPOs []KnowledgeItemPO
	if err := r.db.WithContext(ctx).Where("pack_id = ?", pack.ID).Order("sort_order ASC, id ASC").Find(&itemPOs).Error; err != nil {
		return nil, err
	}
	pack.Items = make([]packdomain.KnowledgeItem, len(itemPOs))
	for i := range itemPOs {
		pack.Items[i] = poToItem(&itemPOs[i])
	}
	var citationPOs []KnowledgeCitationPO
	if err := r.db.WithContext(ctx).Where("pack_id = ?", pack.ID).Order("sort_order ASC, id ASC").Find(&citationPOs).Error; err != nil {
		return nil, err
	}
	pack.Citations = make([]packdomain.KnowledgeCitation, len(citationPOs))
	for i := range citationPOs {
		pack.Citations[i] = poToCitation(&citationPOs[i])
	}
	return pack, nil
}

func createItems(tx *gorm.DB, packID uint64, items []packdomain.KnowledgeItem) error {
	if len(items) == 0 {
		return nil
	}
	pos := make([]KnowledgeItemPO, len(items))
	for i := range items {
		itemToPO(packID, &items[i], &pos[i])
	}
	return tx.Create(&pos).Error
}

func createCitations(tx *gorm.DB, packID uint64, citations []packdomain.KnowledgeCitation) error {
	if len(citations) == 0 {
		return nil
	}
	pos := make([]KnowledgeCitationPO, len(citations))
	for i := range citations {
		citationToPO(packID, &citations[i], &pos[i])
	}
	return tx.Create(&pos).Error
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&KnowledgePackPO{}, &KnowledgeItemPO{}, &KnowledgeCitationPO{})
}
