package builtin

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	runtimetool "lattice-coding/internal/runtime/tool"
)

const defaultFileListLimit = 200

type FileListTool struct {
	readOnlyTool
	maxEntries int
}

func NewFileListTool() *FileListTool {
	return &FileListTool{maxEntries: defaultFileListLimit}
}

func (t *FileListTool) Name() string { return FileListName }

func (t *FileListTool) Description() string {
	return "List files and directories inside the working directory."
}

func (t *FileListTool) Schema() runtimetool.Schema {
	return runtimetool.ObjectSchema(map[string]interface{}{
		"dir_path":  runtimetool.StringSchema("Directory path, relative to working_dir unless absolute."),
		"recursive": runtimetool.BooleanSchema("Whether to recursively list child directories."),
	}, "dir_path")
}

func (t *FileListTool) Validate(_ context.Context, input map[string]interface{}) error {
	if _, err := requiredString(input, "dir_path"); err != nil {
		return err
	}
	_, err := optionalBool(input, "recursive", false)
	return err
}

func (t *FileListTool) Execute(_ context.Context, req runtimetool.ToolRequest) (runtimetool.ToolOutput, error) {
	dirPath, _ := requiredString(req.Input, "dir_path")
	recursive, _ := optionalBool(req.Input, "recursive", false)

	resolved, relRoot, err := resolveInsideWorkingDir(req.Context.WorkingDir, dirPath)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	if !info.IsDir() {
		return runtimetool.ToolOutput{}, fmt.Errorf("path is not a directory: %s", dirPath)
	}

	entries, truncated, err := t.listEntries(resolved, relRoot, recursive)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })

	content := formatFileEntries(entries, truncated)
	return runtimetool.ToolOutput{
		Content:   content,
		Truncated: truncated,
		Data: map[string]interface{}{
			"entries":   entriesToData(entries),
			"truncated": truncated,
			"limit":     t.maxEntries,
		},
	}, nil
}

func (t *FileListTool) listEntries(root, relRoot string, recursive bool) ([]fileEntry, bool, error) {
	limit := t.maxEntries
	if limit <= 0 {
		limit = defaultFileListLimit
	}

	entries := make([]fileEntry, 0)
	addEntry := func(path string, info fs.FileInfo) bool {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return false
		}
		displayPath := relRoot
		if rel != "." {
			displayPath = filepath.ToSlash(filepath.Join(relRoot, rel))
		}
		entries = append(entries, fileEntry{
			Path:  displayPath,
			IsDir: info.IsDir(),
			Size:  info.Size(),
		})
		return len(entries) < limit
	}

	truncated := false
	if recursive {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == root {
				return nil
			}
			if d.IsDir() && shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			if !addEntry(path, info) {
				truncated = true
				if d.IsDir() {
					return filepath.SkipDir
				}
				return errStopWalk
			}
			return nil
		})
		if err == errStopWalk {
			err = nil
		}
		return entries, truncated, err
	}

	children, err := os.ReadDir(root)
	if err != nil {
		return nil, false, err
	}
	for _, child := range children {
		if child.IsDir() && shouldSkipDir(child.Name()) {
			continue
		}
		info, err := child.Info()
		if err != nil {
			return nil, false, err
		}
		if !addEntry(filepath.Join(root, child.Name()), info) {
			truncated = true
			break
		}
	}
	return entries, truncated, nil
}

type fileEntry struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

var errStopWalk = fmt.Errorf("stop walk")

func formatFileEntries(entries []fileEntry, truncated bool) string {
	if len(entries) == 0 {
		if truncated {
			return "[truncated]"
		}
		return ""
	}
	lines := make([]string, 0, len(entries)+1)
	for _, entry := range entries {
		kind := "file"
		if entry.IsDir {
			kind = "dir"
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%d", kind, entry.Path, entry.Size))
	}
	if truncated {
		lines = append(lines, fmt.Sprintf("[truncated after %d entries]", len(entries)))
	}
	return joinLines(lines)
}

func entriesToData(entries []fileEntry) []map[string]interface{} {
	data := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		data = append(data, map[string]interface{}{
			"path":   entry.Path,
			"is_dir": entry.IsDir,
			"size":   entry.Size,
		})
	}
	return data
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	result := lines[0]
	for _, line := range lines[1:] {
		result += "\n" + line
	}
	return result
}
