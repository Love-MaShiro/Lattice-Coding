package prompt

import (
	"encoding/json"
	"strings"
)

func formatTools(tools []ToolPrompt) string {
	if len(tools) == 0 {
		return "No tools are available."
	}
	var b strings.Builder
	for _, tool := range tools {
		b.WriteString("- ")
		b.WriteString(tool.Name)
		if tool.Description != "" {
			b.WriteString(": ")
			b.WriteString(tool.Description)
		}
		if tool.Schema != nil {
			if data, err := json.Marshal(tool.Schema); err == nil {
				b.WriteString("\n  args_schema: ")
				b.Write(data)
			}
		}
		if tool.Destructive {
			b.WriteString("\n  destructive: true")
		}
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}

func filterTools(tools []ToolPrompt, allowed []string) []ToolPrompt {
	if len(allowed) == 0 {
		return tools
	}
	allowedSet := map[string]bool{}
	for _, name := range allowed {
		allowedSet[name] = true
	}
	filtered := make([]ToolPrompt, 0, len(tools))
	for _, tool := range tools {
		if allowedSet[tool.Name] {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}
