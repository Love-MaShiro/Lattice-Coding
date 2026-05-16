package prompt

import "context"

type Builder interface {
	Build(ctx context.Context, req Request) (*Prompt, error)
}

type TemplateRenderer interface {
	Render(ctx context.Context, template string, data map[string]interface{}) (string, error)
}

type IncludeResolver interface {
	Resolve(ctx context.Context, name string) (string, error)
}

type Request struct {
	System    string
	User      string
	Tools     []ToolPrompt
	Workflow  string
	Variables map[string]interface{}
}

type Prompt struct {
	System   string
	Messages []Message
	Metadata map[string]interface{}
}

type Message struct {
	Role    string
	Content string
}

type ToolPrompt struct {
	Name        string
	Description string
	Schema      map[string]interface{}
}
