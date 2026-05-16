package prompt

import "context"

type Message struct {
	Role    string
	Content string
}

type Prompt struct {
	System   string
	Messages []Message
	Metadata map[string]interface{}
}

type Request struct {
	WorkingDir            string
	UserInput             string
	System                string
	Workflow              string
	NodeType              string
	NodeName              string
	Knowledge             string
	AgentConfig           string
	DeferredTools         string
	Shell                 string
	LocalInstructionFiles []string
	InstructionDirs       []string
	Variables             map[string]interface{}
	AllowedTools          []string
}

type ToolPrompt struct {
	Name        string
	Description string
	Schema      map[string]interface{}
	ReadOnly    bool
	Destructive bool
}

type ToolDescriber interface {
	DescribeTools(ctx context.Context, toolCtx ToolContext, allowed []string) ([]ToolPrompt, error)
}

type ToolContext struct {
	WorkingDir string
}
