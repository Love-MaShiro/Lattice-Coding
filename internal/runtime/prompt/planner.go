package prompt

import (
	"context"
	"strings"
)

func (b *PromptBuilder) BuildPlanGraphPrompt(ctx context.Context, req Request) (*Prompt, error) {
	systemPrompt, err := b.BuildSystemPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	tools, err := b.tools(ctx, req)
	if err != nil {
		return nil, err
	}
	var planner strings.Builder
	planner.WriteString(systemPrompt)
	planner.WriteString("\n\n# PlanGraph Mode\n")
	planner.WriteString("Output exactly one WorkflowSpec JSON object and nothing else. Do not output markdown.\n")
	planner.WriteString("Only DAG workflows are allowed. Allowed node types: llm, tool, condition, end. Allowed strategies: direct_chat, pure_react, fixed_workflow.\n")
	planner.WriteString("Do not use dangerous, destructive, shell, or unlisted tools. Every edge must reference existing nodes and cycles are forbidden.\n\n")
	planner.WriteString("# Available Tools\n")
	planner.WriteString(formatTools(nonDestructiveTools(tools)))
	return &Prompt{
		System: planner.String(),
		Messages: []Message{
			{Role: "user", Content: req.UserInput},
		},
		Metadata: map[string]interface{}{"kind": "plan_graph"},
	}, nil
}

func nonDestructiveTools(tools []ToolPrompt) []ToolPrompt {
	filtered := make([]ToolPrompt, 0, len(tools))
	for _, tool := range tools {
		if !tool.Destructive && !strings.Contains(strings.ToLower(tool.Name), "shell") {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}
