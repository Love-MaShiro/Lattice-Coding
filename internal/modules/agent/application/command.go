package application

type CreateAgentCommand struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ProviderID    uint64 `json:"provider_id"`
	ModelConfigID uint64 `json:"model_config_id"`
	SystemPrompt  string `json:"system_prompt"`
	Tools         string `json:"tools"`
	MaxSteps      int    `json:"max_steps"`
	Timeout       int    `json:"timeout"`
	Enabled       bool   `json:"enabled"`
}

type UpdateAgentCommand struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ProviderID    uint64 `json:"provider_id"`
	ModelConfigID uint64 `json:"model_config_id"`
	SystemPrompt  string `json:"system_prompt"`
	Tools         string `json:"tools"`
	MaxSteps      int    `json:"max_steps"`
	Timeout       int    `json:"timeout"`
	Enabled       *bool  `json:"enabled"`
}
