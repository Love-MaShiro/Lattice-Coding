package knowledge

import (
	"lattice-coding/internal/modules/knowledge/api"
	contextapp "lattice-coding/internal/modules/knowledge/context/application"
	contextdomain "lattice-coding/internal/modules/knowledge/context/domain"
	contextpersistence "lattice-coding/internal/modules/knowledge/context/infra/persistence"
	packdomain "lattice-coding/internal/modules/knowledge/pack/domain"
	packpersistence "lattice-coding/internal/modules/knowledge/pack/infra/persistence"
	ragpersistence "lattice-coding/internal/modules/knowledge/rag/infra/persistence"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB *gorm.DB
}

type Module struct {
	ContextSourceRepo    contextdomain.SourceRepository
	ContextCandidateRepo contextdomain.CandidateRepository
	ContextPolicyRepo    contextdomain.PolicyRepository
	ContextSourceSvc     contextapp.SourceService
	ContextCandidateSvc  contextapp.CandidateService
	ContextPolicySvc     contextapp.PolicyService
	ContextPackSvc       contextapp.ContextPackService
	PackRepo             packdomain.PackRepository
}

func NewModule(p *ModuleProvider) *Module {
	_ = contextpersistence.Migrate(p.DB)
	_ = packpersistence.Migrate(p.DB)
	_ = ragpersistence.Migrate(p.DB)

	sourceRepo := contextpersistence.NewSourceRepositoryImpl(p.DB)
	candidateRepo := contextpersistence.NewCandidateRepositoryImpl(p.DB)
	policyRepo := contextpersistence.NewPolicyRepositoryImpl(p.DB)

	return &Module{
		ContextSourceRepo:    sourceRepo,
		ContextCandidateRepo: candidateRepo,
		ContextPolicyRepo:    policyRepo,
		ContextSourceSvc:     contextapp.NewSourceService(sourceRepo),
		ContextCandidateSvc:  contextapp.NewCandidateService(candidateRepo),
		ContextPolicySvc:     contextapp.NewPolicyService(policyRepo),
		ContextPackSvc:       contextapp.NewContextPackService(candidateRepo, policyRepo),
		PackRepo:             packpersistence.NewPackRepositoryImpl(p.DB),
	}
}

func RegisterRoutes(group *gin.RouterGroup, m *Module) {
	_ = m
	api.RegisterRoutes(group)
}
