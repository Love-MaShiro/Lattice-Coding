package domain

type AgentType string

const (
	AgentTypeCustomerService     AgentType = "customer_service"
	AgentTypeCodeReview          AgentType = "code_review"
	AgentTypeCodeFix             AgentType = "code_fix"
	AgentTypeTestGenerate        AgentType = "test_generate"
	AgentTypeExperimentReproduce AgentType = "experiment_reproduce"
)
