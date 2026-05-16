package agent

import "testing"

func TestParseReActActionToolCall(t *testing.T) {
	action, err := ParseReActAction(`{"type":"tool_call","reason":"inspect file","tool":"file.read","args":{"path":"README.md"}}`)
	if err != nil {
		t.Fatalf("ParseReActAction returned error: %v", err)
	}
	if action.Type != ReActActionToolCall {
		t.Fatalf("type = %q, want %q", action.Type, ReActActionToolCall)
	}
	if action.Tool != "file.read" {
		t.Fatalf("tool = %q, want file.read", action.Tool)
	}
	if action.Args["path"] != "README.md" {
		t.Fatalf("args[path] = %v, want README.md", action.Args["path"])
	}
}

func TestParseReActActionFinal(t *testing.T) {
	action, err := ParseReActAction(`{"type":"final","answer":"done"}`)
	if err != nil {
		t.Fatalf("ParseReActAction returned error: %v", err)
	}
	if action.Type != ReActActionFinal {
		t.Fatalf("type = %q, want %q", action.Type, ReActActionFinal)
	}
	if action.Answer != "done" {
		t.Fatalf("answer = %q, want done", action.Answer)
	}
}

func TestParseReActActionRejectsMarkdown(t *testing.T) {
	_, err := ParseReActAction("```json\n{\"type\":\"final\",\"answer\":\"done\"}\n```")
	if err == nil {
		t.Fatal("ParseReActAction returned nil error for markdown")
	}
}

func TestParseReActActionRejectsUnknownToolCallFields(t *testing.T) {
	_, err := ParseReActAction(`{"type":"tool_call","reason":"x","tool":"file.read","args":{},"extra":true}`)
	if err == nil {
		t.Fatal("ParseReActAction returned nil error for unknown field")
	}
}

func TestParseReActActionRequiresFinalAnswer(t *testing.T) {
	_, err := ParseReActAction(`{"type":"final","answer":""}`)
	if err == nil {
		t.Fatal("ParseReActAction returned nil error for empty final answer")
	}
}
