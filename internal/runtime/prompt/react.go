package prompt

import (
	"context"
	"strings"
)

func (b *PromptBuilder) BuildReActPrompt(ctx context.Context, req Request) (*Prompt, error) {
	systemPrompt, err := b.BuildSystemPrompt(ctx, req)
	if err != nil {
		return nil, err
	}
	var react strings.Builder
	react.WriteString(systemPrompt)
	tools, err := b.tools(ctx, req)
	if err != nil {
		return nil, err
	}
	react.WriteString("\n\n# ReAct JSON Mode\n")
	react.WriteString("Output exactly one raw JSON object and nothing else. Do not output markdown.\n")
	react.WriteString("Allowed formats:\n")
	react.WriteString(`{"type":"tool_call","reason":"short reason","tool":"tool.name","args":{}}`)
	react.WriteString("\n")
	react.WriteString(`{"type":"final","answer":"final answer"}`)
	react.WriteString("\nRules: keep reason short; do not reveal hidden reasoning; do not repeat the same tool call with the same args; do not invent tool names.\n")
	react.WriteString("\n# Available Tools\n")
	react.WriteString(formatTools(tools))
	return &Prompt{
		System: react.String(),
		Messages: []Message{
			{Role: "user", Content: req.UserInput},
		},
		Metadata: map[string]interface{}{"kind": "react"},
	}, nil
}
