package prompt

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxIncludeDepth = 5

type IncludeResolver interface {
	Resolve(ctx context.Context, baseDir string, content string) (string, error)
}

type FileIncludeResolver struct {
	MaxDepth int
}

func NewFileIncludeResolver() FileIncludeResolver {
	return FileIncludeResolver{MaxDepth: maxIncludeDepth}
}

func (r FileIncludeResolver) Resolve(ctx context.Context, baseDir string, content string) (string, error) {
	maxDepth := r.MaxDepth
	if maxDepth <= 0 {
		maxDepth = maxIncludeDepth
	}
	visited := map[string]bool{}
	return resolveIncludes(ctx, baseDir, content, 0, maxDepth, visited)
}

func resolveIncludes(ctx context.Context, baseDir string, content string, depth int, maxDepth int, visited map[string]bool) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var out strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		path, ok := parseIncludeLine(line)
		if !ok {
			out.WriteString(line)
			out.WriteString("\n")
			continue
		}
		resolved, err := resolveIncludePath(baseDir, path)
		if err != nil {
			out.WriteString(includeComment(path, err.Error()))
			continue
		}
		if depth >= maxDepth {
			out.WriteString(includeComment(path, "max include depth reached"))
			continue
		}
		if visited[resolved] {
			out.WriteString(includeComment(path, "include cycle skipped"))
			continue
		}
		data, err := os.ReadFile(resolved)
		if err != nil {
			out.WriteString(includeComment(path, "missing include"))
			continue
		}
		visited[resolved] = true
		included, err := resolveIncludes(ctx, filepath.Dir(resolved), string(data), depth+1, maxDepth, visited)
		delete(visited, resolved)
		if err != nil {
			return "", err
		}
		out.WriteString(included)
		if !strings.HasSuffix(included, "\n") {
			out.WriteString("\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func parseIncludeLine(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "@") {
		return "", false
	}
	if strings.Contains(trimmed, " ") || strings.Contains(trimmed, "\t") {
		return "", false
	}
	path := strings.TrimPrefix(trimmed, "@")
	if path == "" {
		return "", false
	}
	return path, true
}

func resolveIncludePath(baseDir string, includePath string) (string, error) {
	if strings.HasPrefix(includePath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Abs(filepath.Join(home, strings.TrimPrefix(includePath, "~/")))
	}
	if filepath.IsAbs(includePath) {
		return filepath.Clean(includePath), nil
	}
	if strings.HasPrefix(includePath, "./") || strings.HasPrefix(includePath, "../") {
		return filepath.Abs(filepath.Join(baseDir, includePath))
	}
	return "", fmt.Errorf("unsupported include path")
}

func includeComment(path string, reason string) string {
	return fmt.Sprintf("<!-- include %s skipped: %s -->\n", path, reason)
}
