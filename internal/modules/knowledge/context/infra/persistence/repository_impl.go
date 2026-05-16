package persistence

import (
	"context"

	commondb "lattice-coding/internal/common/db"
	contextdomain "lattice-coding/internal/modules/knowledge/context/domain"

	"gorm.io/gorm"
)

type SourceRepositoryImpl struct {
	db *gorm.DB
}

func NewSourceRepositoryImpl(db *gorm.DB) contextdomain.SourceRepository {
	return &SourceRepositoryImpl{db: db}
}

func (r *SourceRepositoryImpl) Create(ctx context.Context, source *contextdomain.ContextSource) error {
	po := &ContextSourcePO{}
	sourceToPO(source, po)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	*source = *poToSource(po)
	return nil
}

func (r *SourceRepositoryImpl) FindByKey(ctx context.Context, sourceKey string) (*contextdomain.ContextSource, error) {
	var po ContextSourcePO
	if err := r.db.WithContext(ctx).Where("source_key = ?", sourceKey).First(&po).Error; err != nil {
		return nil, err
	}
	return poToSource(&po), nil
}

func (r *SourceRepositoryImpl) DeleteByKey(ctx context.Context, sourceKey string) error {
	return r.db.WithContext(ctx).Where("source_key = ?", sourceKey).Delete(&ContextSourcePO{}).Error
}

type CandidateRepositoryImpl struct {
	db *gorm.DB
}

func NewCandidateRepositoryImpl(db *gorm.DB) contextdomain.CandidateRepository {
	return &CandidateRepositoryImpl{db: db}
}

func (r *CandidateRepositoryImpl) CreateWithSignals(ctx context.Context, candidate *contextdomain.ContextCandidate) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		po := &ContextCandidatePO{}
		candidateToPO(candidate, po)
		if err := tx.Create(po).Error; err != nil {
			return err
		}
		candidate.ID = po.ID
		candidate.CreatedAt = po.CreatedAt
		candidate.UpdatedAt = po.UpdatedAt
		return createSignals(tx, candidate.ID, candidate.Signals)
	})
}

func (r *CandidateRepositoryImpl) FindByKeyWithSignals(ctx context.Context, candidateKey string) (*contextdomain.ContextCandidate, error) {
	var po ContextCandidatePO
	if err := r.db.WithContext(ctx).Where("candidate_key = ?", candidateKey).First(&po).Error; err != nil {
		return nil, err
	}
	return r.loadCandidateSignals(ctx, &po)
}

func (r *CandidateRepositoryImpl) FindBySourceKey(ctx context.Context, sourceKey string) ([]*contextdomain.ContextCandidate, error) {
	var pos []ContextCandidatePO
	if err := r.db.WithContext(ctx).Where("source_key = ?", sourceKey).Order("score DESC, id DESC").Find(&pos).Error; err != nil {
		return nil, err
	}
	items := make([]*contextdomain.ContextCandidate, len(pos))
	for i := range pos {
		items[i] = poToCandidate(&pos[i])
	}
	return items, nil
}

func (r *CandidateRepositoryImpl) DeleteByKey(ctx context.Context, candidateKey string) error {
	return commondb.Transaction(ctx, r.db, func(tx *gorm.DB) error {
		var po ContextCandidatePO
		if err := tx.Where("candidate_key = ?", candidateKey).First(&po).Error; err != nil {
			return err
		}
		if err := tx.Where("candidate_id = ?", po.ID).Delete(&ContextSignalPO{}).Error; err != nil {
			return err
		}
		return tx.Delete(&po).Error
	})
}

func (r *CandidateRepositoryImpl) loadCandidateSignals(ctx context.Context, po *ContextCandidatePO) (*contextdomain.ContextCandidate, error) {
	candidate := poToCandidate(po)
	var signalPOs []ContextSignalPO
	if err := r.db.WithContext(ctx).Where("candidate_id = ?", candidate.ID).Order("weight DESC, id ASC").Find(&signalPOs).Error; err != nil {
		return nil, err
	}
	candidate.Signals = make([]contextdomain.ContextSignal, len(signalPOs))
	for i := range signalPOs {
		candidate.Signals[i] = poToSignal(&signalPOs[i])
	}
	return candidate, nil
}

type PolicyRepositoryImpl struct {
	db *gorm.DB
}

func NewPolicyRepositoryImpl(db *gorm.DB) contextdomain.PolicyRepository {
	return &PolicyRepositoryImpl{db: db}
}

func (r *PolicyRepositoryImpl) Create(ctx context.Context, policy *contextdomain.ContextPolicy) error {
	po := &ContextPolicyPO{}
	policyToPO(policy, po)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	*policy = *poToPolicy(po)
	return nil
}

func (r *PolicyRepositoryImpl) Update(ctx context.Context, policy *contextdomain.ContextPolicy) error {
	po := &ContextPolicyPO{}
	policyToPO(policy, po)
	return r.db.WithContext(ctx).Model(po).Omit("created_at").Updates(po).Error
}

func (r *PolicyRepositoryImpl) FindByKey(ctx context.Context, policyKey string) (*contextdomain.ContextPolicy, error) {
	var po ContextPolicyPO
	if err := r.db.WithContext(ctx).Where("policy_key = ?", policyKey).First(&po).Error; err != nil {
		return nil, err
	}
	return poToPolicy(&po), nil
}

func (r *PolicyRepositoryImpl) DeleteByKey(ctx context.Context, policyKey string) error {
	return r.db.WithContext(ctx).Where("policy_key = ?", policyKey).Delete(&ContextPolicyPO{}).Error
}

func createSignals(tx *gorm.DB, candidateID uint64, signals []contextdomain.ContextSignal) error {
	if len(signals) == 0 {
		return nil
	}
	pos := make([]ContextSignalPO, len(signals))
	for i := range signals {
		signalToPO(candidateID, &signals[i], &pos[i])
	}
	return tx.Create(&pos).Error
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&ContextSourcePO{}, &ContextCandidatePO{}, &ContextSignalPO{}, &ContextPolicyPO{})
}
