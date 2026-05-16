package prompt

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GitContext struct {
	Branch        string
	RecentCommits string
	Status        string
}

func (g GitContext) Empty() bool {
	return g.Branch == "" && g.RecentCommits == "" && g.Status == ""
}

func (g GitContext) String() string {
	var b strings.Builder
	if g.Branch != "" {
		b.WriteString("Branch: ")
		b.WriteString(g.Branch)
		b.WriteString("\n")
	}
	if g.RecentCommits != "" {
		b.WriteString("Recent commits:\n")
		b.WriteString(g.RecentCommits)
		b.WriteString("\n")
	}
	if g.Status != "" {
		b.WriteString("Status:\n")
		b.WriteString(g.Status)
	}
	return strings.TrimSpace(b.String())
}

func LoadGitContext(ctx context.Context, workingDir string) GitContext {
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	if !isGitWorkTree(ctx, workingDir) {
		return GitContext{}
	}
	return GitContext{
		Branch:        runGit(ctx, workingDir, "branch", "--show-current"),
		RecentCommits: runGit(ctx, workingDir, "log", "--oneline", "-5"),
		Status:        runGit(ctx, workingDir, "status", "--short"),
	}
}

func isGitWorkTree(ctx context.Context, workingDir string) bool {
	return runGit(ctx, workingDir, "rev-parse", "--is-inside-work-tree") == "true"
}

func runGit(ctx context.Context, workingDir string, args ...string) string {
	gitCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(gitCtx, "git", args...)
	cmd.Dir = workingDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(stdout.String())
}
