package prompt

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var instructionFiles = []string{"LATTICE.md", "AGENTS.md", "CLAUDE.md"}

type ProjectInstructionLoader struct {
	IncludeResolver IncludeResolver
	GlobalPaths     []string
	UserPaths       []string
}

func NewProjectInstructionLoader() ProjectInstructionLoader {
	return ProjectInstructionLoader{
		IncludeResolver: NewFileIncludeResolver(),
		GlobalPaths:     defaultGlobalInstructionPaths(),
		UserPaths:       defaultUserInstructionPaths(),
	}
}

func (l ProjectInstructionLoader) Load(ctx context.Context, workingDir string) (string, error) {
	return l.LoadForRequest(ctx, Request{WorkingDir: workingDir})
}

func (l ProjectInstructionLoader) LoadForRequest(ctx context.Context, req Request) (string, error) {
	workingDir := req.WorkingDir
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	abs, err := filepath.Abs(workingDir)
	if err != nil {
		return "", err
	}
	var sections []string

	for _, path := range l.GlobalPaths {
		if section, ok, err := l.loadFile(ctx, path); err != nil {
			return "", err
		} else if ok {
			sections = append(sections, section)
		}
	}
	for _, path := range l.UserPaths {
		if section, ok, err := l.loadFile(ctx, path); err != nil {
			return "", err
		} else if ok {
			sections = append(sections, section)
		}
	}

	projectDirs := dirsFromRoot(abs)
	for _, dir := range projectDirs {
		dirSections, err := l.loadInstructionDir(ctx, dir)
		if err != nil {
			return "", err
		}
		sections = append(sections, dirSections...)

		rules, err := LoadRules(ctx, dir, l.IncludeResolver)
		if err != nil {
			return "", err
		}
		if rules != "" {
			sections = append(sections, rules)
		}
	}

	for _, path := range req.LocalInstructionFiles {
		if section, ok, err := l.loadFile(ctx, path); err != nil {
			return "", err
		} else if ok {
			sections = append(sections, section)
		}
	}
	for _, dir := range req.InstructionDirs {
		dirSections, err := l.loadInstructionDir(ctx, dir)
		if err != nil {
			return "", err
		}
		sections = append(sections, dirSections...)
	}

	return strings.TrimSpace(strings.Join(sections, "\n\n")), nil
}

func (l ProjectInstructionLoader) loadInstructionDir(ctx context.Context, dir string) ([]string, error) {
	var sections []string
	for _, name := range instructionFiles {
		path := filepath.Join(dir, name)
		if section, ok, err := l.loadFile(ctx, path); err != nil {
			return nil, err
		} else if ok {
			sections = append(sections, section)
		}
	}
	return sections, nil
}

func (l ProjectInstructionLoader) loadFile(ctx context.Context, path string) (string, bool, error) {
	if path == "" {
		return "", false, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", false, err
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return "", false, nil
	}
	content := string(data)
	if l.IncludeResolver != nil {
		content, err = l.IncludeResolver.Resolve(ctx, filepath.Dir(abs), content)
		if err != nil {
			return "", false, err
		}
	}
	return "## " + abs + "\n" + strings.TrimSpace(content), true, nil
}

func dirsFromRoot(start string) []string {
	dirs := dirsToRoot(start)
	for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}
	return dirs
}

func dirsToRoot(start string) []string {
	var dirs []string
	current := filepath.Clean(start)
	for {
		dirs = append(dirs, current)
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return dirs
}

func defaultGlobalInstructionPaths() []string {
	if runtime.GOOS == "windows" {
		programData := os.Getenv("ProgramData")
		if programData == "" {
			return nil
		}
		return []string{
			filepath.Join(programData, "Lattice-Coding", "CLAUDE.md"),
			filepath.Join(programData, "Lattice-Coding", "AGENTS.md"),
		}
	}
	return []string{
		"/etc/lattice/CLAUDE.md",
		"/etc/lattice/AGENTS.md",
	}
}

func defaultUserInstructionPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return nil
	}
	return []string{
		filepath.Join(home, ".lattice", "CLAUDE.md"),
		filepath.Join(home, ".lattice", "AGENTS.md"),
		filepath.Join(home, ".claude", "CLAUDE.md"),
		filepath.Join(home, "CLAUDE.md"),
	}
}
