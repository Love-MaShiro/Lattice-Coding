package builtin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	runtimetool "lattice-coding/internal/runtime/tool"
)

type FileReadTool struct {
	readOnlyTool
	stateManager runtimetool.FileReadStateManager
}

func NewFileReadTool(stateManager runtimetool.FileReadStateManager) *FileReadTool {
	if stateManager == nil {
		stateManager = runtimetool.NewInMemoryFileReadStateManager()
	}
	return &FileReadTool{stateManager: stateManager}
}

func (t *FileReadTool) Name() string { return FileReadName }

func (t *FileReadTool) Description() string {
	return "Read a file inside the working directory and return line-numbered content."
}

func (t *FileReadTool) Schema() runtimetool.Schema {
	return runtimetool.ObjectSchema(map[string]interface{}{
		"file_path": runtimetool.StringSchema("Path to the file, relative to working_dir unless absolute."),
	}, "file_path")
}

func (t *FileReadTool) Validate(_ context.Context, input map[string]interface{}) error {
	_, err := requiredString(input, "file_path")
	return err
}

func (t *FileReadTool) Execute(_ context.Context, req runtimetool.ToolRequest) (runtimetool.ToolOutput, error) {
	filePath, _ := requiredString(req.Input, "file_path")
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

	content, err := os.ReadFile(resolved)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}

	checksum := recordFileReadState(t.stateManager, resolved, info, content)

	return runtimetool.ToolOutput{
		Content: lineNumberedContent(string(content)),
		Data: map[string]interface{}{
			"file":     rel,
			"mtime":    info.ModTime().UTC().Format("2006-01-02T15:04:05.000000000Z07:00"),
			"sha256":   checksum,
			"size":     info.Size(),
			"lineNums": true,
		},
	}, nil
}

func lineNumberedContent(content string) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	line := 1
	for scanner.Scan() {
		fmt.Fprintf(&builder, "%6d | %s\n", line, scanner.Text())
		line++
	}
	if content == "" {
		return ""
	}
	return strings.TrimRight(builder.String(), "\n")
}
