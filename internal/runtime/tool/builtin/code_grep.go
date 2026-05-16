package builtin

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	runtimetool "lattice-coding/internal/runtime/tool"
)

const (
	defaultGrepMaxResults = 100
	grepTimeout           = 10 * time.Second
)

type CodeGrepTool struct {
	readOnlyTool
}

func NewCodeGrepTool() *CodeGrepTool {
	return &CodeGrepTool{}
}

func (t *CodeGrepTool) Name() string { return CodeGrepName }

func (t *CodeGrepTool) Description() string {
	return "Search text in files under the working directory."
}

func (t *CodeGrepTool) Schema() runtimetool.Schema {
	return runtimetool.ObjectSchema(map[string]interface{}{
		"pattern":     runtimetool.StringSchema("Regular expression pattern to search for."),
		"path":        runtimetool.StringSchema("File or directory path, relative to working_dir unless absolute."),
		"max_results": runtimetool.NumberSchema("Maximum number of matches to return."),
	}, "pattern", "path")
}

func (t *CodeGrepTool) Validate(_ context.Context, input map[string]interface{}) error {
	if _, err := requiredString(input, "pattern"); err != nil {
		return err
	}
	if _, err := requiredString(input, "path"); err != nil {
		return err
	}
	maxResults, err := optionalInt(input, "max_results", defaultGrepMaxResults)
	if err != nil {
		return err
	}
	if maxResults <= 0 {
		return fmt.Errorf("max_results must be greater than 0")
	}
	return nil
}

func (t *CodeGrepTool) Execute(ctx context.Context, req runtimetool.ToolRequest) (runtimetool.ToolOutput, error) {
	pattern, _ := requiredString(req.Input, "pattern")
	searchPath, _ := requiredString(req.Input, "path")
	maxResults, _ := optionalInt(req.Input, "max_results", defaultGrepMaxResults)

	resolved, rel, err := resolveInsideWorkingDir(req.Context.WorkingDir, searchPath)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	if _, err := os.Stat(resolved); err != nil {
		return runtimetool.ToolOutput{}, err
	}

	matches, hasMore, engine, err := t.search(ctx, req.Context.WorkingDir, rel, pattern, maxResults)
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}

	return runtimetool.ToolOutput{
		Content: formatGrepMatches(matches, hasMore),
		Data: map[string]interface{}{
			"matches":     grepMatchesToData(matches),
			"has_more":    hasMore,
			"max_results": maxResults,
			"engine":      engine,
		},
	}, nil
}

func (t *CodeGrepTool) search(ctx context.Context, workingDir, path, pattern string, maxResults int) ([]grepMatch, bool, string, error) {
	if _, err := exec.LookPath("rg"); err == nil {
		matches, hasMore, err := runSearchCommand(ctx, workingDir, "rg", []string{
			"--line-number",
			"--no-heading",
			"--color", "never",
			"--glob", "!{.git,node_modules,dist,build,.idea,.vscode}/**",
			pattern,
			path,
		}, maxResults)
		if err == nil {
			return matches, hasMore, "rg", nil
		}
	}

	if _, err := exec.LookPath("grep"); err == nil {
		matches, hasMore, err := runSearchCommand(ctx, workingDir, "grep", []string{
			"-RIn",
			"--exclude-dir=.git",
			"--exclude-dir=node_modules",
			"--exclude-dir=dist",
			"--exclude-dir=build",
			"--exclude-dir=.idea",
			"--exclude-dir=.vscode",
			"--",
			pattern,
			path,
		}, maxResults)
		if err == nil {
			return matches, hasMore, "grep", nil
		}
	}

	matches, hasMore, err := scanWithGo(workingDir, path, pattern, maxResults)
	return matches, hasMore, "go", err
}

func runSearchCommand(ctx context.Context, workingDir, command string, args []string, maxResults int) ([]grepMatch, bool, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, grepTimeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, command, args...)
	cmd.Dir = workingDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if cmdCtx.Err() != nil {
		return nil, false, cmdCtx.Err()
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("%s failed: %s", command, strings.TrimSpace(stderr.String()))
	}

	return parseGrepOutput(stdout.String(), maxResults), grepOutputHasMore(stdout.String(), maxResults), nil
}

func parseGrepOutput(output string, maxResults int) []grepMatch {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	matches := make([]grepMatch, 0, minInt(len(lines), maxResults))
	for _, line := range lines {
		if line == "" || len(matches) >= maxResults {
			continue
		}
		match, ok := parseGrepLine(line)
		if ok {
			matches = append(matches, match)
		}
	}
	return matches
}

func grepOutputHasMore(output string, maxResults int) bool {
	if strings.TrimSpace(output) == "" {
		return false
	}
	return len(strings.Split(strings.TrimRight(output, "\n"), "\n")) > maxResults
}

func parseGrepLine(line string) (grepMatch, bool) {
	first := strings.Index(line, ":")
	if first < 0 {
		return grepMatch{}, false
	}
	second := strings.Index(line[first+1:], ":")
	if second < 0 {
		return grepMatch{}, false
	}
	second += first + 1
	lineNumber := 0
	for _, ch := range line[first+1 : second] {
		if ch < '0' || ch > '9' {
			return grepMatch{}, false
		}
		lineNumber = lineNumber*10 + int(ch-'0')
	}
	return grepMatch{
		File: filepath.ToSlash(line[:first]),
		Line: lineNumber,
		Text: line[second+1:],
	}, true
}

func scanWithGo(workingDir, path, pattern string, maxResults int) ([]grepMatch, bool, error) {
	expr, err := regexp.Compile(pattern)
	if err != nil {
		return nil, false, err
	}

	root := filepath.Join(workingDir, filepath.FromSlash(path))
	info, err := os.Stat(root)
	if err != nil {
		return nil, false, err
	}

	matches := make([]grepMatch, 0, maxResults)
	hasMore := false
	visitFile := func(filePath string) error {
		fileMatches, err := scanFile(workingDir, filePath, expr, maxResults-len(matches)+1)
		if err != nil {
			return err
		}
		for _, match := range fileMatches {
			if len(matches) >= maxResults {
				hasMore = true
				return errStopWalk
			}
			matches = append(matches, match)
		}
		return nil
	}

	if !info.IsDir() {
		err := visitFile(root)
		if err == errStopWalk {
			err = nil
		}
		return matches, hasMore, err
	}

	err = filepath.WalkDir(root, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if filePath != root && shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		return visitFile(filePath)
	})
	if err == errStopWalk {
		err = nil
	}
	return matches, hasMore, err
}

func scanFile(workingDir, filePath string, expr *regexp.Regexp, limit int) ([]grepMatch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	matches := make([]grepMatch, 0)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		text := scanner.Text()
		if expr.MatchString(text) {
			rel, err := filepath.Rel(workingDir, filePath)
			if err != nil {
				return nil, err
			}
			matches = append(matches, grepMatch{
				File: filepath.ToSlash(rel),
				Line: lineNumber,
				Text: text,
			})
			if len(matches) >= limit {
				return matches, nil
			}
		}
	}
	return matches, scanner.Err()
}

type grepMatch struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

func formatGrepMatches(matches []grepMatch, hasMore bool) string {
	lines := make([]string, 0, len(matches)+1)
	for _, match := range matches {
		lines = append(lines, fmt.Sprintf("%s:%d:%s", match.File, match.Line, match.Text))
	}
	if hasMore {
		lines = append(lines, "[more results omitted]")
	}
	return joinLines(lines)
}

func grepMatchesToData(matches []grepMatch) []map[string]interface{} {
	data := make([]map[string]interface{}, 0, len(matches))
	for _, match := range matches {
		data = append(data, map[string]interface{}{
			"file": match.File,
			"line": match.Line,
			"text": match.Text,
		})
	}
	return data
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
