package application

type CreateProviderCommand struct {
	Name            string
	ProviderType    string
	BaseURL         string
	AuthType        string
	APIKey          string
	AuthConfig      string
	Config          string
	Enabled         bool
}

type UpdateProviderCommand struct {
	ID              uint64
	Name            string
	ProviderType    string
	BaseURL         string
	AuthType        string
	APIKey          string
	AuthConfig      string
	Config          string
	Enabled         *bool
}

type CreateModelConfigCommand struct {
	ProviderID   uint64
	Name         string
	Model        string
	ModelType    string
	Params       string
	Capabilities string
	IsDefault    bool
	Enabled      bool
}

type UpdateModelConfigCommand struct {
	ID           uint64
	Name         string
	Model        string
	ModelType    string
	Params       string
	Capabilities string
	IsDefault    *bool
	Enabled      *bool
}
