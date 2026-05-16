package builtin

import (
	"context"
	"fmt"
	"os"
	"strings"

	runtimetool "lattice-coding/internal/runtime/tool"
)

type FileEditTool struct {
	runtimetool.BaseTool
	stateManager runtimetool.FileReadStateManager
}

func NewFileEditTool(stateManager runtimetool.FileReadStateManager) *FileEditTool {
	if stateManager == nil {
		stateManager = runtimetool.NewInMemoryFileReadStateManager()
	}
	return &FileEditTool{stateManager: stateManager}
}

func (t *FileEditTool) Name() string { return FileEditName }

func (t *FileEditTool) Description() string {
	return "Replace one unique string in a file after read-before-edit and mtime checks."
}

func (t *FileEditTool) Schema() runtimetool.Schema {
	return runtimetool.ObjectSchema(map[string]interface{}{
		"file_path":  runtimetool.StringSchema("Path to the file, relative to working_dir unless absolute."),
		"old_string": runtimetool.StringSchema("Existing text to replace. Must appear exactly once."),
		"new_string": runtimetool.StringSchema("Replacement text."),
	}, "file_path", "old_string", "new_string")
}

func (t *FileEditTool) Validate(_ context.Context, input map[string]interface{}) error {
	if _, err := requiredString(input, "file_path"); err != nil {
		return err
	}
	if _, err := requiredString(input, "old_string"); err != nil {
		return err
	}
	_, err := requiredStringAllowEmpty(input, "new_string")
	return err
}

func (t *FileEditTool) IsReadOnly() bool        { return false }
func (t *FileEditTool) IsConcurrencySafe() bool { return false }
func (t *FileEditTool) IsDestructive() bool     { return false }

func (t *FileEditTool) CheckPermission(context.Context, runtimetool.ToolRequest) (runtimetool.PermissionDecision, string, error) {
	return runtimetool.PermissionAllow, "", nil
}

func (t *FileEditTool) Execute(_ context.Context, req runtimetool.ToolRequest) (runtimetool.ToolOutput, error) {
	filePath, _ := requiredString(req.Input, "file_path")
	oldString, _ := requiredString(req.Input, "old_string")
	newString, _ := requiredStringAllowEmpty(req.Input, "new_string")

	resolved, rel, err := resolveInsideWorkingDir(req.Context.WorkingDir, filePath)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	if info.IsDir() {
		return runtimetool.ToolOutput{}, fmt.Errorf("path is a directory: %s", filePath)
	}

	readState, ok := t.stateManager.Get(resolved)
	if !ok {
		return runtimetool.ToolOutput{}, fmt.Errorf("file must be read with file.read before editing: %s", rel)
	}
	if readState.MTime != info.ModTime().UnixNano() {
		return runtimetool.ToolOutput{}, fmt.Errorf("file changed since last read; please run file.read again before editing: %s", rel)
	}

	content, err := os.ReadFile(resolved)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	if currentChecksum := fileChecksum(content); readState.Checksum != "" && readState.Checksum != currentChecksum {
		return runtimetool.ToolOutput{}, fmt.Errorf("file content changed since last read; please run file.read again before editing: %s", rel)
	}

	text := string(content)
	count := strings.Count(text, oldString)
	if count == 0 {
		return runtimetool.ToolOutput{}, fmt.Errorf("old_string was not found in %s; please run file.read again before editing", rel)
	}
	if count > 1 {
		return runtimetool.ToolOutput{}, fmt.Errorf("old_string appears %d times in %s; provide more precise context", count, rel)
	}

	updated := strings.Replace(text, oldString, newString, 1)
	if err := os.WriteFile(resolved, []byte(updated), info.Mode().Perm()); err != nil {
		return runtimetool.ToolOutput{}, err
	}

	updatedInfo, err := os.Stat(resolved)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	updatedChecksum := recordFileReadState(t.stateManager, resolved, updatedInfo, []byte(updated))

	contentText := fmt.Sprintf(
		"Edited %s\nReplaced 1 occurrence\n\nDiff summary:\n- %s\n+ %s",
		rel,
		summarizeEditText(oldString),
		summarizeEditText(newString),
	)
	return runtimetool.ToolOutput{
		Content: contentText,
		Data: map[string]interface{}{
			"file":       rel,
			"old_sha256": readState.Checksum,
			"new_sha256": updatedChecksum,
			"replaced":   1,
			"diff": map[string]interface{}{
				"old": oldString,
				"new": newString,
			},
		},
	}, nil
}

func summarizeEditText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\n", "\\n")
	if len(text) <= 160 {
		return text
	}
	return text[:157] + "..."
}
