package domain

import (
	"encoding/json"
	"fmt"
)

type NodeConfig interface {
	NodeType() NodeType
	RawJSON() string
	sealedNodeConfig()
}

type BaseNodeConfig struct {
	Type NodeType        `json:"type,omitempty"`
	Raw  json.RawMessage `json:"-"`
}

func (c BaseNodeConfig) NodeType() NodeType {
	return c.Type
}

func (c BaseNodeConfig) RawJSON() string {
	if len(c.Raw) == 0 {
		return "{}"
	}
	return string(c.Raw)
}

func (c BaseNodeConfig) sealedNodeConfig() {}

type StartNodeConfig struct {
	BaseNodeConfig
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

type EndNodeConfig struct {
	BaseNodeConfig
	OutputSchema json.RawMessage `json:"output_schema,omitempty"`
}

type LLMNodeConfig struct {
	BaseNodeConfig
	ModelConfigID uint64          `json:"model_config_id,omitempty"`
	SystemPrompt  string          `json:"system_prompt,omitempty"`
	Prompt        string          `json:"prompt,omitempty"`
	Temperature   *float64        `json:"temperature,omitempty"`
	MaxTokens     int             `json:"max_tokens,omitempty"`
	OutputSchema  json.RawMessage `json:"output_schema,omitempty"`
}

type ToolNodeConfig struct {
	BaseNodeConfig
	ToolName string          `json:"tool_name,omitempty"`
	Args     json.RawMessage `json:"args,omitempty"`
	Timeout  int             `json:"timeout_ms,omitempty"`
}

type ConditionNodeConfig struct {
	BaseNodeConfig
	Expression string `json:"expression,omitempty"`
}

type KnowledgeRouteNodeConfig struct {
	BaseNodeConfig
	Strategy string `json:"strategy,omitempty"`
}

type KnowledgeRetrieveNodeConfig struct {
	BaseNodeConfig
	Route     string   `json:"route,omitempty"`
	SourceIDs []uint64 `json:"source_ids,omitempty"`
	Limit     int      `json:"limit,omitempty"`
}

type ContextBuildNodeConfig struct {
	BaseNodeConfig
	MaxTokens int     `json:"max_tokens,omitempty"`
	MinScore  float64 `json:"min_score,omitempty"`
	MaxItems  int     `json:"max_items,omitempty"`
}

type GenericNodeConfig struct {
	BaseNodeConfig
}

type NodeConfigParser interface {
	Parse(nodeType NodeType, raw string) (NodeConfig, error)
	Marshal(config NodeConfig) (string, error)
}

type JSONNodeConfigParser struct{}

func NewJSONNodeConfigParser() NodeConfigParser {
	return &JSONNodeConfigParser{}
}

func (p *JSONNodeConfigParser) Parse(nodeType NodeType, raw string) (NodeConfig, error) {
	if raw == "" {
		raw = "{}"
	}
	rawMessage := json.RawMessage(raw)
	base := BaseNodeConfig{Type: nodeType, Raw: rawMessage}

	switch nodeType {
	case NodeTypeStart:
		var cfg StartNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeEnd:
		var cfg EndNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeLLM:
		var cfg LLMNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeTool, NodeTypeWebSearch, NodeTypeMCPCall, NodeTypeCodeSearch, NodeTypeFileRead, NodeTypeShellCommand:
		var cfg ToolNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeCondition, NodeTypeParallel:
		var cfg ConditionNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeKnowledgeRoute:
		var cfg KnowledgeRouteNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeKnowledgeRetrieve:
		var cfg KnowledgeRetrieveNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	case NodeTypeContextBuild, NodeTypeContextCompress:
		var cfg ContextBuildNodeConfig
		if err := unmarshalConfig(rawMessage, &cfg, base); err != nil {
			return nil, err
		}
		return cfg, nil
	default:
		return GenericNodeConfig{BaseNodeConfig: base}, nil
	}
}

func (p *JSONNodeConfigParser) Marshal(config NodeConfig) (string, error) {
	if config == nil {
		return "{}", nil
	}
	raw := config.RawJSON()
	if raw == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(raw)) {
		return "", fmt.Errorf("invalid node config json")
	}
	return raw, nil
}

func unmarshalConfig[T interface{ SetBase(BaseNodeConfig) }](raw json.RawMessage, cfg T, base BaseNodeConfig) error {
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, cfg); err != nil {
			return err
		}
	}
	cfg.SetBase(base)
	return nil
}

func (c *StartNodeConfig) SetBase(base BaseNodeConfig)             { c.BaseNodeConfig = base }
func (c *EndNodeConfig) SetBase(base BaseNodeConfig)               { c.BaseNodeConfig = base }
func (c *LLMNodeConfig) SetBase(base BaseNodeConfig)               { c.BaseNodeConfig = base }
func (c *ToolNodeConfig) SetBase(base BaseNodeConfig)              { c.BaseNodeConfig = base }
func (c *ConditionNodeConfig) SetBase(base BaseNodeConfig)         { c.BaseNodeConfig = base }
func (c *KnowledgeRouteNodeConfig) SetBase(base BaseNodeConfig)    { c.BaseNodeConfig = base }
func (c *KnowledgeRetrieveNodeConfig) SetBase(base BaseNodeConfig) { c.BaseNodeConfig = base }
func (c *ContextBuildNodeConfig) SetBase(base BaseNodeConfig)      { c.BaseNodeConfig = base }
