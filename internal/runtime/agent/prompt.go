package agent

import (
	"encoding/json"
	"strings"

	"lattice-coding/internal/runtime/tool"
)

func BuildReActMessages(req Request, tools []tool.ToolDescriptor, steps []ReActStep) []string {
	return []string{
		buildReActSystemPrompt(tools),
		buildReActUserPrompt(req, steps),
	}
}

func buildReActSystemPrompt(tools []tool.ToolDescriptor) string {
	var b strings.Builder
	b.WriteString("You are an agent runtime using the ReAct pattern.\n")
	b.WriteString("You must respond with exactly one raw JSON object and nothing else.\n")
	b.WriteString("Do not use markdown, code fences, comments, or explanatory text outside JSON.\n")
	b.WriteString("Do not invent tool names. Use only the tools listed below.\n")
	b.WriteString("Keep reason short.\n\n")
	b.WriteString("To call a tool, output:\n")
	b.WriteString(`{"type":"tool_call","reason":"short reason","tool":"tool.name","args":{}}`)
	b.WriteString("\n\nTo finish, output:\n")
	b.WriteString(`{"type":"final","answer":"final answer"}`)
	b.WriteString("\n\nAvailable tools:\n")
	if len(tools) == 0 {
		b.WriteString("none\n")
		return b.String()
	}
	for _, descriptor := range tools {
		b.WriteString("- ")
		b.WriteString(descriptor.Name)
		if descriptor.Description != "" {
			b.WriteString(": ")
			b.WriteString(descriptor.Description)
		}
		if descriptor.Schema != nil {
			if data, err := json.Marshal(descriptor.Schema); err == nil {
				b.WriteString("\n  args_schema: ")
				b.Write(data)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func buildReActUserPrompt(req Request, steps []ReActStep) string {
	var b strings.Builder
	b.WriteString("Task:\n")
	b.WriteString(req.Input)
	if len(steps) == 0 {
		return b.String()
	}

	b.WriteString("\n\nPrevious steps:\n")
	for _, step := range steps {
		b.WriteString("Iteration ")
		b.WriteString(intString(step.Iteration))
		b.WriteString("\nReason: ")
		b.WriteString(step.Reason)
		b.WriteString("\nAction: ")
		b.WriteString(step.Action)
		if step.ActionInput != nil {
			if data, err := json.Marshal(step.ActionInput); err == nil {
				b.WriteString("\nAction input: ")
				b.Write(data)
			}
		}
		b.WriteString("\nObservation: ")
		b.WriteString(step.Observation)
		b.WriteString("\n")
	}
	return b.String()
}

func intString(value int) string {
	data, _ := json.Marshal(value)
	return string(data)
}
