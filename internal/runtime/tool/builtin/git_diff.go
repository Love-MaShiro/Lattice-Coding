package builtin

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	runtimetool "lattice-coding/internal/runtime/tool"
)

const gitDiffTimeout = 10 * time.Second

type GitDiffTool struct {
	readOnlyTool
}

func NewGitDiffTool() *GitDiffTool {
	return &GitDiffTool{}
}

func (t *GitDiffTool) Name() string { return GitDiffName }

func (t *GitDiffTool) Description() string {
	return "Return git diff output for the working directory."
}

func (t *GitDiffTool) Schema() runtimetool.Schema {
	return runtimetool.ObjectSchema(map[string]interface{}{
		"staged": runtimetool.BooleanSchema("Whether to return staged changes."),
	})
}

func (t *GitDiffTool) Validate(_ context.Context, input map[string]interface{}) error {
	_, err := optionalBool(input, "staged", false)
	return err
}

func (t *GitDiffTool) Execute(ctx context.Context, req runtimetool.ToolRequest) (runtimetool.ToolOutput, error) {
	workingDir, _, err := resolveInsideWorkingDir(req.Context.WorkingDir, ".")
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}
	staged, _ := optionalBool(req.Input, "staged", false)

	args := []string{"diff"}
	if staged {
		args = append(args, "--staged")
	}

	cmdCtx, cancel := context.WithTimeout(ctx, gitDiffTimeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "git", args...)
	cmd.Dir = workingDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	combined := strings.TrimRight(stdout.String()+stderr.String(), "\n")
	if cmdCtx.Err() != nil {
		return runtimetool.ToolOutput{Content: combined}, cmdCtx.Err()
	}
	if err != nil {
		if combined == "" {
			combined = err.Error()
		}
		return runtimetool.ToolOutput{Content: combined}, fmt.Errorf("%s", combined)
	}

	return runtimetool.ToolOutput{
		Content: combined,
		Data: map[string]interface{}{
			"staged": staged,
		},
	}, nil
}
