package prompt

import (
	"context"
	"strings"
)

func (b *PromptBuilder) BuildWorkflowNodePrompt(ctx context.Context, req Request) (*Prompt, error) {
	systemPrompt, err := b.BuildSystemPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	tools, err := b.tools(ctx, req)
	if err != nil {
		return nil, err
	}
	var workflow strings.Builder
	workflow.WriteString(systemPrompt)
	workflow.WriteString("\n\n# Workflow Node Mode\n")
	workflow.WriteString("Execute only the current node. Respect node config, workflow state, and allowed tools.\n")
	if req.NodeName != "" {
		workflow.WriteString("Node name: ")
		workflow.WriteString(req.NodeName)
		workflow.WriteString("\n")
	}
	if req.NodeType != "" {
		workflow.WriteString("Node type: ")
		workflow.WriteString(req.NodeType)
		workflow.WriteString("\n")
	}
	if req.Workflow != "" {
		workflow.WriteString("\n# Workflow Context\n")
		workflow.WriteString(req.Workflow)
	}
	if req.DeferredTools != "" {
		workflow.WriteString("\n\n# Deferred Tools\n")
		workflow.WriteString(req.DeferredTools)
	}
	workflow.WriteString("\n\n# Available Tools\n")
	workflow.WriteString(formatTools(tools))
	return &Prompt{
		System:   workflow.String(),
		Messages: []Message{{Role: "user", Content: req.UserInput}},
		Metadata: map[string]interface{}{"kind": "workflow_node"},
	}, nil
}
