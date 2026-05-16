package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	runDomain "lattice-coding/internal/modules/run/domain"
	"lattice-coding/internal/runtime/llm"
	runtimetool "lattice-coding/internal/runtime/tool"
)

const codeReviewWorkflowID = "code_review_fixed"

type CodeReviewService struct {
	runRepo       runDomain.RunRepository
	toolExecutor  *runtimetool.ToolExecutor
	llmExecutor   *llm.Executor
	auditRecorder runtimetool.AuditRecorder
	dag           *DAGExecutor
}

type CodeReviewCommand struct {
	RunID         string
	AgentID       string
	SessionID     string
	WorkingDir    string
	ModelConfigID uint64
	MaxFiles      int
	MaxChars      int
}

type CodeReviewResult struct {
	RunID  string `json:"run_id"`
	Report string `json:"report"`
}

func NewCodeReviewService(runRepo runDomain.RunRepository, toolExecutor *runtimetool.ToolExecutor, llmExecutor *llm.Executor, auditRecorder runtimetool.AuditRecorder) *CodeReviewService {
	return &CodeReviewService{
		runRepo:       runRepo,
		toolExecutor:  toolExecutor,
		llmExecutor:   llmExecutor,
		auditRecorder: auditRecorder,
		dag:           NewDAGExecutor(),
	}
}

func (s *CodeReviewService) Run(ctx context.Context, cmd CodeReviewCommand) (*CodeReviewResult, error) {
	if cmd.RunID == "" {
		cmd.RunID = fmt.Sprintf("code-review-%d", time.Now().UnixNano())
	}
	if cmd.MaxFiles <= 0 {
		cmd.MaxFiles = 8
	}
	if cmd.MaxChars <= 0 {
		cmd.MaxChars = 60000
	}

	inputJSON, _ := json.Marshal(cmd)
	startedAt := time.Now()
	run := &runDomain.Run{
		ID:         cmd.RunID,
		AgentID:    cmd.AgentID,
		SessionID:  cmd.SessionID,
		WorkflowID: codeReviewWorkflowID,
		Status:     runDomain.RunStatusRunning,
		Input:      string(inputJSON),
		StartedAt:  startedAt,
	}
	if err := s.runRepo.Create(ctx, run); err != nil {
		return nil, err
	}

	state := &WorkflowState{
		RunID:      cmd.RunID,
		WorkingDir: cmd.WorkingDir,
		Input: map[string]interface{}{
			"model_config_id": cmd.ModelConfigID,
			"max_files":       cmd.MaxFiles,
			"max_chars":       cmd.MaxChars,
		},
		Values: map[string]interface{}{},
	}

	err := s.dag.Execute(ctx, s.nodes(cmd), s.edges(), state)
	completedAt := time.Now()
	run.CompletedAt = &completedAt
	if err != nil {
		run.Status = runDomain.RunStatusFailed
		run.Error = err.Error()
		_ = s.runRepo.Update(ctx, run)
		return nil, err
	}
	run.Status = runDomain.RunStatusCompleted
	run.Output = state.Output
	if err := s.runRepo.Update(ctx, run); err != nil {
		return nil, err
	}
	return &CodeReviewResult{RunID: cmd.RunID, Report: state.Output}, nil
}

func (s *CodeReviewService) nodes(cmd CodeReviewCommand) []DAGNode {
	return []DAGNode{
		{Key: "start", Run: s.auditOnlyNode("start")},
		{Key: "git_diff_node", Run: s.gitDiffNode},
		{Key: "changed_file_context_node", Run: s.changedFileContextNode},
		{Key: "related_code_search_node", Run: s.relatedCodeSearchNode},
		{Key: "llm_review_node", Run: s.llmReviewNode(cmd.ModelConfigID)},
		{Key: "final_report_node", Run: s.finalReportNode},
		{Key: "end", Run: s.auditOnlyNode("end")},
	}
}

func (s *CodeReviewService) edges() []DAGEdge {
	return []DAGEdge{
		{From: "start", To: "git_diff_node"},
		{From: "git_diff_node", To: "changed_file_context_node"},
		{From: "changed_file_context_node", To: "related_code_search_node"},
		{From: "related_code_search_node", To: "llm_review_node"},
		{From: "llm_review_node", To: "final_report_node"},
		{From: "final_report_node", To: "end"},
	}
}

func (s *CodeReviewService) auditOnlyNode(nodeID string) func(context.Context, *WorkflowState) error {
	return func(ctx context.Context, state *WorkflowState) error {
		s.recordNode(ctx, state, nodeID, false, "")
		return nil
	}
}

func (s *CodeReviewService) gitDiffNode(ctx context.Context, state *WorkflowState) error {
	s.recordNode(ctx, state, "git_diff_node", false, "started")
	result := s.executeTool(ctx, state, "git_diff_node", "git.diff", map[string]interface{}{})
	if result.IsError {
		s.recordNode(ctx, state, "git_diff_node", true, result.Error)
		return errors.New(result.Error)
	}
	state.Values["diff"] = result.Content
	s.recordNode(ctx, state, "git_diff_node", false, "finished")
	return nil
}

func (s *CodeReviewService) changedFileContextNode(ctx context.Context, state *WorkflowState) error {
	s.recordNode(ctx, state, "changed_file_context_node", false, "started")
	diff, _ := state.Values["diff"].(string)
	files := parseChangedFiles(diff)
	maxFiles, _ := state.Input["max_files"].(int)
	maxChars, _ := state.Input["max_chars"].(int)
	contexts := map[string]string{}
	usedChars := 0
	for _, file := range files {
		if len(contexts) >= maxFiles || usedChars >= maxChars {
			break
		}
		result := s.executeTool(ctx, state, "changed_file_context_node", "file.read", map[string]interface{}{"file_path": file})
		if result.IsError {
			contexts[file] = "ERROR: " + result.Error
			continue
		}
		content := result.Content
		if usedChars+len(content) > maxChars {
			content = content[:maxChars-usedChars]
		}
		contexts[file] = content
		usedChars += len(content)
	}
	state.Values["changed_files"] = files
	state.Values["file_contexts"] = contexts
	s.recordNode(ctx, state, "changed_file_context_node", false, "finished")
	return nil
}

func (s *CodeReviewService) relatedCodeSearchNode(ctx context.Context, state *WorkflowState) error {
	s.recordNode(ctx, state, "related_code_search_node", false, "started")
	files, _ := state.Values["changed_files"].([]string)
	patterns := buildSearchPatterns(files)
	results := map[string]string{}
	for _, pattern := range patterns {
		result := s.executeTool(ctx, state, "related_code_search_node", "code.grep", map[string]interface{}{
			"pattern":     pattern,
			"path":        ".",
			"max_results": 20,
		})
		if result.IsError {
			results[pattern] = "ERROR: " + result.Error
			continue
		}
		results[pattern] = result.Content
	}
	state.Values["related_search"] = results
	s.recordNode(ctx, state, "related_code_search_node", false, "finished")
	return nil
}

func (s *CodeReviewService) llmReviewNode(modelConfigID uint64) func(context.Context, *WorkflowState) error {
	return func(ctx context.Context, state *WorkflowState) error {
		s.recordNode(ctx, state, "llm_review_node", false, "started")
		if modelConfigID == 0 {
			return fmt.Errorf("model_config_id is required")
		}
		prompt := buildReviewPrompt(state)
		resp, call := s.llmExecutor.Chat(ctx, llm.ChatRequest{
			ModelConfigID: modelConfigID,
			Messages: []llm.Message{
				{Role: "system", Content: "You are a senior code reviewer. Be concise, specific, and risk-focused."},
				{Role: "user", Content: prompt},
			},
			MaxTokens: 4096,
		})
		if !call.Success || resp == nil {
			if call.Error != nil {
				s.recordNode(ctx, state, "llm_review_node", true, call.Error.Error())
				return call.Error
			}
			return fmt.Errorf("llm review failed")
		}
		state.Values["review"] = resp.Content
		s.recordNode(ctx, state, "llm_review_node", false, "finished")
		return nil
	}
}

func (s *CodeReviewService) finalReportNode(ctx context.Context, state *WorkflowState) error {
	s.recordNode(ctx, state, "final_report_node", false, "started")
	review, _ := state.Values["review"].(string)
	files, _ := state.Values["changed_files"].([]string)
	state.Output = fmt.Sprintf("# AI Code Review\n\nChanged files: %d\n\n%s", len(files), review)
	s.recordNode(ctx, state, "final_report_node", false, "finished")
	return nil
}

func (s *CodeReviewService) executeTool(ctx context.Context, state *WorkflowState, nodeID, name string, input map[string]interface{}) runtimetool.ToolResult {
	return s.toolExecutor.Execute(ctx, runtimetool.ToolRequest{
		Name:  name,
		Input: input,
		Context: runtimetool.ToolContext{
			RunID:      state.RunID,
			WorkingDir: state.WorkingDir,
			Metadata: map[string]interface{}{
				"node_id": nodeID,
			},
		},
	})
}

func (s *CodeReviewService) recordNode(ctx context.Context, state *WorkflowState, nodeID string, isError bool, message string) {
	if s.auditRecorder == nil {
		return
	}
	eventType := "workflow_node_finished"
	if message == "started" {
		eventType = "workflow_node_started"
	}
	if isError {
		eventType = "workflow_node_failed"
	}
	_ = s.auditRecorder.Record(ctx, runtimetool.AuditEvent{
		EventType: eventType,
		Request: runtimetool.ToolRequest{
			Name: "workflow.node",
			Context: runtimetool.ToolContext{
				RunID: state.RunID,
				Metadata: map[string]interface{}{
					"node_id": nodeID,
				},
			},
		},
		Result: runtimetool.ToolResult{
			ToolName: "workflow.node",
			IsError:  isError,
			Content:  message,
			Error:    errorString(isError, message),
		},
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
	})
}

func parseChangedFiles(diff string) []string {
	seen := map[string]struct{}{}
	files := make([]string, 0)
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "diff --git ") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				file := strings.TrimPrefix(parts[3], "b/")
				if file != "/dev/null" {
					if _, ok := seen[file]; !ok {
						seen[file] = struct{}{}
						files = append(files, file)
					}
				}
			}
		}
	}
	return files
}

func buildSearchPatterns(files []string) []string {
	seen := map[string]struct{}{}
	patterns := make([]string, 0)
	for _, file := range files {
		base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		if base != "" {
			seen[regexp.QuoteMeta(base)] = struct{}{}
		}
		seen[`func\s+`] = struct{}{}
		seen[`import\s+`] = struct{}{}
	}
	for pattern := range seen {
		patterns = append(patterns, pattern)
	}
	if len(patterns) > 5 {
		return patterns[:5]
	}
	return patterns
}

func buildReviewPrompt(state *WorkflowState) string {
	diff, _ := state.Values["diff"].(string)
	files, _ := state.Values["changed_files"].([]string)
	fileContexts, _ := json.MarshalIndent(state.Values["file_contexts"], "", "  ")
	related, _ := json.MarshalIndent(state.Values["related_search"], "", "  ")
	return fmt.Sprintf(`Review the following code changes.

Return sections:
1. Severe issues
2. Potential bugs
3. Concurrency / transaction / security risks
4. Maintainability suggestions
5. Block merge? yes/no with reason

Changed files:
%s

Git diff:
%s

Changed file context:
%s

Related code search:
%s
`, strings.Join(files, "\n"), diff, string(fileContexts), string(related))
}

func errorString(isError bool, message string) string {
	if isError {
		return message
	}
	return ""
}
