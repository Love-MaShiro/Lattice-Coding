package prompt

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	rendered, err := renderTemplate("hello {{.Name}}, {{.Task}}", map[string]interface{}{
		"Name": "Lattice",
		"Task": "build",
	})
	if err != nil {
		t.Fatalf("renderTemplate returned error: %v", err)
	}
	if rendered != "hello Lattice, build" {
		t.Fatalf("rendered = %q", rendered)
	}
}

func TestIncludeResolver(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "part.md"), []byte("included content"), 0644); err != nil {
		t.Fatal(err)
	}
	resolver := NewFileIncludeResolver()
	got, err := resolver.Resolve(context.Background(), dir, "before\n@./part.md\nafter")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if !strings.Contains(got, "before\nincluded content\nafter") {
		t.Fatalf("resolved content = %q", got)
	}
}

func TestIncludeResolverCycleDoesNotLoop(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("A\n@./b.md"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("B\n@./a.md"), 0644); err != nil {
		t.Fatal(err)
	}
	resolver := NewFileIncludeResolver()
	got, err := resolver.Resolve(context.Background(), dir, "@./a.md")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if !strings.Contains(got, "include ./a.md skipped: include cycle skipped") {
		t.Fatalf("cycle comment missing: %q", got)
	}
	if strings.Count(got, "A") > 2 || strings.Count(got, "B") > 2 {
		t.Fatalf("cycle appears unresolved: %q", got)
	}
}

func TestLoadRulesSortedByFilename(t *testing.T) {
	dir := t.TempDir()
	rulesDir := filepath.Join(dir, ".lattice", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rulesDir, "b.md"), []byte("second"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rulesDir, "a.md"), []byte("first"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadRules(context.Background(), dir, nil)
	if err != nil {
		t.Fatalf("LoadRules returned error: %v", err)
	}
	firstIdx := strings.Index(got, "first")
	secondIdx := strings.Index(got, "second")
	if firstIdx < 0 || secondIdx < 0 || firstIdx > secondIdx {
		t.Fatalf("rules not sorted: %q", got)
	}
}

func TestProjectInstructionHierarchyLoadsClosestLater(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "parent")
	child := filepath.Join(parent, "child")
	commandDir := filepath.Join(root, "cmd")
	for _, dir := range []string{parent, child, commandDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(parent, "CLAUDE.md"), "parent rule")
	writeFile(t, filepath.Join(child, "CLAUDE.md"), "child rule")
	localFile := filepath.Join(root, "local.md")
	writeFile(t, localFile, "local rule")
	writeFile(t, filepath.Join(commandDir, "CLAUDE.md"), "command rule")

	loader := ProjectInstructionLoader{
		IncludeResolver: NewFileIncludeResolver(),
		GlobalPaths:     nil,
		UserPaths:       nil,
	}
	got, err := loader.LoadForRequest(context.Background(), Request{
		WorkingDir:            child,
		LocalInstructionFiles: []string{localFile},
		InstructionDirs:       []string{commandDir},
	})
	if err != nil {
		t.Fatalf("LoadForRequest returned error: %v", err)
	}

	assertOrdered(t, got, "parent rule", "child rule", "local rule", "command rule")
}

func TestLoadGitContextNonGitDirReturnsEmpty(t *testing.T) {
	got := LoadGitContext(context.Background(), t.TempDir())
	if !got.Empty() {
		t.Fatalf("git context should be empty: %+v", got)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertOrdered(t *testing.T, text string, values ...string) {
	t.Helper()
	previous := -1
	for _, value := range values {
		idx := strings.Index(text, value)
		if idx < 0 {
			t.Fatalf("%q not found in:\n%s", value, text)
		}
		if idx <= previous {
			t.Fatalf("%q was not loaded after previous value in:\n%s", value, text)
		}
		previous = idx
	}
}

func TestBuildReActPromptContainsJSONConstraints(t *testing.T) {
	builder := NewBuilder(WithToolDescriber(staticToolDescriber{
		tools: []ToolPrompt{{Name: "file.read", Description: "read file"}},
	}))
	prompt, err := builder.BuildReActPrompt(context.Background(), Request{
		WorkingDir:   t.TempDir(),
		UserInput:    "inspect README",
		AllowedTools: []string{"file.read"},
	})
	if err != nil {
		t.Fatalf("BuildReActPrompt returned error: %v", err)
	}
	checks := []string{
		"exactly one raw JSON object",
		"Do not output markdown",
		"Do not invent tool names",
		`"type":"tool_call"`,
		`"type":"final"`,
		"file.read",
	}
	for _, check := range checks {
		if !strings.Contains(prompt.System, check) {
			t.Fatalf("prompt missing %q:\n%s", check, prompt.System)
		}
	}
}

type staticToolDescriber struct {
	tools []ToolPrompt
}

func (d staticToolDescriber) DescribeTools(context.Context, ToolContext, []string) ([]ToolPrompt, error) {
	return d.tools, nil
}
