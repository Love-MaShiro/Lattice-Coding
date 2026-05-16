package builtin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	runtimetool "lattice-coding/internal/runtime/tool"
)

const (
	FileReadName = "file.read"
	FileEditName = "file.edit"
	FileListName = "file.list"
	CodeGrepName = "code.grep"
	GitDiffName  = "git.diff"
	ShellRunName = "shell.run"
)

var filteredDirs = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"dist":         {},
	"build":        {},
	".idea":        {},
	".vscode":      {},
}

type readOnlyTool struct {
	runtimetool.BaseTool
}

func (readOnlyTool) IsReadOnly() bool        { return true }
func (readOnlyTool) IsConcurrencySafe() bool { return true }
func (readOnlyTool) IsDestructive() bool     { return false }

func (readOnlyTool) CheckPermission(context.Context, runtimetool.ToolRequest) (runtimetool.PermissionDecision, string, error) {
	return runtimetool.PermissionAllow, "", nil
}

func requiredString(input map[string]interface{}, key string) (string, error) {
	value, ok := input[key]
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	text, ok := value.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return text, nil
}

func requiredStringAllowEmpty(input map[string]interface{}, key string) (string, error) {
	value, ok := input[key]
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return text, nil
}

func optionalString(input map[string]interface{}, key, defaultValue string) (string, error) {
	value, ok := input[key]
	if !ok || value == nil {
		return defaultValue, nil
	}
	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	if strings.TrimSpace(text) == "" {
		return defaultValue, nil
	}
	return text, nil
}

func optionalBool(input map[string]interface{}, key string, defaultValue bool) (bool, error) {
	value, ok := input[key]
	if !ok || value == nil {
		return defaultValue, nil
	}
	boolean, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return boolean, nil
}

func optionalInt(input map[string]interface{}, key string, defaultValue int) (int, error) {
	value, ok := input[key]
	if !ok || value == nil {
		return defaultValue, nil
	}
	switch typed := value.(type) {
	case int:
		return typed, nil
	case int64:
		return int(typed), nil
	case float64:
		if typed != float64(int(typed)) {
			return 0, fmt.Errorf("%s must be an integer", key)
		}
		return int(typed), nil
	default:
		return 0, fmt.Errorf("%s must be an integer", key)
	}
}

func resolveInsideWorkingDir(workingDir, requestedPath string) (string, string, error) {
	if strings.TrimSpace(workingDir) == "" {
		return "", "", errors.New("working_dir is required")
	}
	if strings.TrimSpace(requestedPath) == "" {
		requestedPath = "."
	}

	root, err := filepath.Abs(filepath.Clean(workingDir))
	if err != nil {
		return "", "", err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", "", err
	}

	target := filepath.Clean(requestedPath)
	if !filepath.IsAbs(target) {
		target = filepath.Join(root, target)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return "", "", err
	}
	if _, statErr := os.Stat(target); statErr == nil {
		target, err = filepath.EvalSymlinks(target)
		if err != nil {
			return "", "", err
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return "", "", statErr
	}

	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", "", fmt.Errorf("path escapes working directory: %s", requestedPath)
	}
	return target, filepath.ToSlash(rel), nil
}

func shouldSkipDir(name string) bool {
	_, ok := filteredDirs[name]
	return ok
}
