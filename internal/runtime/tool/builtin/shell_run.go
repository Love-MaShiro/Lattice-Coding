package builtin

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	runtimetool "lattice-coding/internal/runtime/tool"
)

const (
	defaultShellTimeoutSeconds = 30
	maxShellTimeoutSeconds     = 120
)

type ShellRunTool struct {
	runtimetool.BaseTool
}

func NewShellRunTool() *ShellRunTool {
	return &ShellRunTool{}
}

func (t *ShellRunTool) Name() string { return ShellRunName }

func (t *ShellRunTool) Description() string {
	return "Run a shell command inside the working directory."
}

func (t *ShellRunTool) Schema() runtimetool.Schema {
	return runtimetool.ObjectSchema(map[string]interface{}{
		"command":         runtimetool.StringSchema("Shell command to execute."),
		"timeout_seconds": runtimetool.NumberSchema("Optional timeout in seconds. Defaults to 30, max 120."),
	}, "command")
}

func (t *ShellRunTool) Validate(_ context.Context, input map[string]interface{}) error {
	if _, err := requiredString(input, "command"); err != nil {
		return err
	}
	timeoutSeconds, err := optionalInt(input, "timeout_seconds", defaultShellTimeoutSeconds)
	if err != nil {
		return err
	}
	if timeoutSeconds <= 0 {
		return fmt.Errorf("timeout_seconds must be greater than 0")
	}
	if timeoutSeconds > maxShellTimeoutSeconds {
		return fmt.Errorf("timeout_seconds must be less than or equal to %d", maxShellTimeoutSeconds)
	}
	return nil
}

func (t *ShellRunTool) IsReadOnly() bool        { return false }
func (t *ShellRunTool) IsConcurrencySafe() bool { return false }
func (t *ShellRunTool) IsDestructive() bool     { return false }

func (t *ShellRunTool) CheckPermission(context.Context, runtimetool.ToolRequest) (runtimetool.PermissionDecision, string, error) {
	return runtimetool.PermissionAllow, "", nil
}

func (t *ShellRunTool) Execute(ctx context.Context, req runtimetool.ToolRequest) (runtimetool.ToolOutput, error) {
	command, _ := requiredString(req.Input, "command")
	timeoutSeconds, _ := optionalInt(req.Input, "timeout_seconds", defaultShellTimeoutSeconds)
	workingDir, _, err := resolveInsideWorkingDir(req.Context.WorkingDir, ".")
	if err != nil {
		return runtimetool.ToolOutput{}, err
	}

	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	name, args := shellCommand(command)
	cmd := exec.CommandContext(cmdCtx, name, args...)
	cmd.Dir = workingDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdoutText := stdout.String()
	stderrText := stderr.String()
	exitCode := 0
	timedOut := cmdCtx.Err() == context.DeadlineExceeded
	if err != nil {
		exitCode = -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	output := shellRunOutput(command, stdoutText, stderrText, exitCode, timedOut, timeoutSeconds)
	if timedOut {
		return output, fmt.Errorf("shell command timed out after %d seconds", timeoutSeconds)
	}
	if err != nil {
		return output, fmt.Errorf("shell command exited with code %d", exitCode)
	}
	return output, nil
}

func shellCommand(command string) (string, []string) {
	if runtime.GOOS == "windows" {
		return "powershell.exe", []string{"-NoProfile", "-NonInteractive", "-Command", command}
	}
	return "sh", []string{"-c", command}
}

func shellRunOutput(command, stdoutText, stderrText string, exitCode int, timedOut bool, timeoutSeconds int) runtimetool.ToolOutput {
	content := formatShellOutput(stdoutText, stderrText)
	return runtimetool.ToolOutput{
		Content: content,
		Data: map[string]interface{}{
			"command":         command,
			"exit_code":       exitCode,
			"stdout":          stdoutText,
			"stderr":          stderrText,
			"timed_out":       timedOut,
			"timeout_seconds": timeoutSeconds,
		},
	}
}

func formatShellOutput(stdoutText, stderrText string) string {
	stdoutText = strings.TrimRight(stdoutText, "\n")
	stderrText = strings.TrimRight(stderrText, "\n")
	if stdoutText == "" && stderrText == "" {
		return "(no output)"
	}
	if stdoutText != "" && stderrText != "" {
		return stdoutText + "\n" + stderrText
	}
	if stdoutText != "" {
		return stdoutText
	}
	return stderrText
}
