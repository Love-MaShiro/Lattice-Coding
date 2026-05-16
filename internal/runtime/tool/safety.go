package tool

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type SafetyDecision string

const (
	SafetyAllow        SafetyDecision = "allow"
	SafetyDeny         SafetyDecision = "deny"
	SafetyNeedApproval SafetyDecision = "need_approval"
)

type SafetyCheckResult struct {
	Decision SafetyDecision         `json:"decision"`
	Reason   string                 `json:"reason,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type SafetyChecker interface {
	Check(ctx context.Context, req ToolRequest, descriptor ToolDescriptor) (SafetyCheckResult, error)
}

type NoopSafetyChecker struct{}

func (NoopSafetyChecker) Check(context.Context, ToolRequest, ToolDescriptor) (SafetyCheckResult, error) {
	return SafetyCheckResult{Decision: SafetyAllow}, nil
}

type RuleBasedSafetyChecker struct{}

func NewRuleBasedSafetyChecker() RuleBasedSafetyChecker {
	return RuleBasedSafetyChecker{}
}

func (RuleBasedSafetyChecker) Check(_ context.Context, req ToolRequest, _ ToolDescriptor) (SafetyCheckResult, error) {
	if req.Name == "file.edit" {
		return CheckFileEditSafety(), nil
	}
	if req.Name != "shell.run" {
		return SafetyCheckResult{Decision: SafetyAllow}, nil
	}
	value, ok := req.Input["command"]
	if !ok {
		return SafetyCheckResult{Decision: SafetyDeny, Reason: "command is required"}, nil
	}
	command, ok := value.(string)
	if !ok || strings.TrimSpace(command) == "" {
		return SafetyCheckResult{Decision: SafetyDeny, Reason: "command must be a non-empty string"}, nil
	}
	return CheckShellCommandSafety(command), nil
}

func CheckFileEditSafety() SafetyCheckResult {
	env := strings.ToLower(strings.TrimSpace(firstNonEmptyEnv("LATTICE_ENV", "APP_ENV", "GO_ENV")))
	if env == "prod" || env == "production" {
		return SafetyCheckResult{
			Decision: SafetyNeedApproval,
			Reason:   "file.edit requires approval in production",
		}
	}
	return SafetyCheckResult{
		Decision: SafetyAllow,
		Reason:   "file.edit allowed in development mode",
	}
}

func CheckShellCommandSafety(command string) SafetyCheckResult {
	normalized := strings.ToLower(normalizeShellCommand(command))
	if normalized == "" {
		return SafetyCheckResult{Decision: SafetyDeny, Reason: "command is empty"}
	}

	denyPatterns := []struct {
		pattern string
		reason  string
	}{
		{`(^|[;&|]\s*)sudo(\s|$)`, "sudo is not allowed"},
		{`(^|[;&|]\s*)rm\s+[^;&|]*-[^\s;&|]*r[^\s;&|]*f|(^|[;&|]\s*)rm\s+[^;&|]*-[^\s;&|]*f[^\s;&|]*r`, "rm -rf is not allowed"},
		{`(^|[;&|]\s*)chmod\s+-R(\s|$)`, "chmod -R is not allowed"},
		{`curl\b.*\|\s*(sh|bash)\b`, "curl pipe to shell is not allowed"},
		{`wget\b.*\|\s*(sh|bash)\b`, "wget pipe to shell is not allowed"},
		{`(^|[;&|]\s*)mkfs(\.|\s|$)`, "mkfs is not allowed"},
		{`(^|[;&|]\s*)dd(\s|$)`, "dd is not allowed"},
		{`(^|[;&|]\s*)shutdown(\s|$)`, "shutdown is not allowed"},
		{`(^|[;&|]\s*)reboot(\s|$)`, "reboot is not allowed"},
	}
	for _, item := range denyPatterns {
		if regexp.MustCompile(item.pattern).MatchString(normalized) {
			return SafetyCheckResult{Decision: SafetyDeny, Reason: item.reason}
		}
	}

	allowPatterns := []string{
		`^ls(\s|$)`,
		`^pwd$`,
		`^cat\s+.+`,
		`^grep\s+.+`,
		`^rg\s+.+`,
		`^find\s+.+`,
		`^git\s+status(\s|$)`,
		`^git\s+diff(\s|$)`,
		`^go\s+test(\s|$)`,
		`^npm\s+test(\s|$)`,
	}
	for _, pattern := range allowPatterns {
		if regexp.MustCompile(pattern).MatchString(normalized) {
			return SafetyCheckResult{Decision: SafetyAllow}
		}
	}

	return SafetyCheckResult{
		Decision: SafetyNeedApproval,
		Reason:   fmt.Sprintf("command requires approval: %s", command),
	}
}

func normalizeShellCommand(command string) string {
	fields := strings.Fields(command)
	return strings.TrimSpace(strings.Join(fields, " "))
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}
