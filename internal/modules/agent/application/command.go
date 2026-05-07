package application

type CreateAgentCommand struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	AgentType       string  `json:"agent_type"`
	ModelConfigID   uint64  `json:"model_config_id"`
	SystemPrompt    string  `json:"system_prompt"`
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	MaxTokens       int     `json:"max_tokens"`
	MaxContextTurns int     `json:"max_context_turns"`
	MaxSteps        int     `json:"max_steps"`
	Enabled         bool    `json:"enabled"`
}

type UpdateAgentCommand struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	AgentType       string  `json:"agent_type"`
	ModelConfigID   uint64  `json:"model_config_id"`
	SystemPrompt    string  `json:"system_prompt"`
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	MaxTokens       int     `json:"max_tokens"`
	MaxContextTurns int     `json:"max_context_turns"`
	MaxSteps        int     `json:"max_steps"`
	Enabled         *bool   `json:"enabled"`
}
