package application

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"os"
	"strings"
	"time"

	"lattice-coding/internal/common/errors"
	redisutil "lattice-coding/internal/common/redis"
	"lattice-coding/internal/modules/chat/domain"
	"lattice-coding/internal/runtime/llm"
	"lattice-coding/internal/runtime/query"
)

type AgentGetter interface {
	GetAgentForChat(ctx context.Context, id uint64) (*AgentRuntimeDTO, error)
}

type MemoryConfig struct {
	CompressionThreshold int
	RetainAfterCompress  int
	CacheTTL             time.Duration
}

type CommandService struct {
	sessionRepo  domain.SessionRepository
	messageRepo  domain.MessageRepository
	agentGetter  AgentGetter
	queryEngine  *query.QueryEngine
	llmExecutor  *llm.Executor
	redisClient  *redisutil.Client
	memoryConfig MemoryConfig
}

var staticAgentToolNames = []string{
	"file.list",
	"file.read",
	"code.grep",
	"git.diff",
	"file.edit",
	"shell.run",
}

func NewCommandService(
	sessionRepo domain.SessionRepository,
	messageRepo domain.MessageRepository,
	agentGetter AgentGetter,
	queryEngine *query.QueryEngine,
	llmExecutor *llm.Executor,
	redisClient *redisutil.Client,
	memoryConfig MemoryConfig,
) *CommandService {
	memoryConfig = normalizeMemoryConfig(memoryConfig)
	return &CommandService{
		sessionRepo:  sessionRepo,
		messageRepo:  messageRepo,
		agentGetter:  agentGetter,
		queryEngine:  queryEngine,
		llmExecutor:  llmExecutor,
		redisClient:  redisClient,
		memoryConfig: memoryConfig,
	}
}

func normalizeMemoryConfig(cfg MemoryConfig) MemoryConfig {
	if cfg.CompressionThreshold <= 0 {
		cfg.CompressionThreshold = 80
	}
	if cfg.RetainAfterCompress <= 0 {
		cfg.RetainAfterCompress = 20
	}
	if cfg.RetainAfterCompress >= cfg.CompressionThreshold {
		cfg.RetainAfterCompress = cfg.CompressionThreshold / 4
		if cfg.RetainAfterCompress <= 0 {
			cfg.RetainAfterCompress = 1
		}
	}
	if cfg.CacheTTL <= 0 {
		cfg.CacheTTL = 24 * time.Hour
	}
	return cfg
}
func (s *CommandService) CreateSession(ctx context.Context, cmd *CreateSessionCommand) (*SessionDTO, error) {
	if cmd.AgentID == 0 {
		return nil, errors.InvalidArg("agent_id is required")
	}
	agent, err := s.getEnabledAgent(ctx, cmd.AgentID)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(cmd.Title)
	if title == "" {
		title = agent.Name
	}

	session := &domain.ChatSession{
		Title:         title,
		AgentID:       agent.ID,
		ModelConfigID: agent.ModelConfigID,
		Status:        domain.SessionStatusActive,
		Meta:          "{}",
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "create chat session failed")
	}
	return ToSessionDTO(session), nil
}

func (s *CommandService) CreateMessage(ctx context.Context, cmd *CreateMessageCommand) (*MessageDTO, error) {
	if cmd.SessionID == 0 {
		return nil, errors.InvalidArg("session_id is required")
	}
	if strings.TrimSpace(cmd.Content) == "" {
		return nil, errors.InvalidArg("content is required")
	}
	role := normalizeRole(cmd.Role)
	if !isAllowedRole(role) {
		return nil, errors.InvalidArg("invalid message role")
	}
	session, err := s.getSession(ctx, cmd.SessionID)
	if err != nil {
		return nil, err
	}

	message := &domain.ChatMessage{
		SessionID:  cmd.SessionID,
		Role:       role,
		Content:    cmd.Content,
		TokenCount: estimateTokens(cmd.Content),
		Meta:       normalizeJSON(cmd.Meta),
	}
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "create chat message failed")
	}
	_ = s.refreshContextCache(ctx, session)
	return ToMessageDTO(message), nil
}

func (s *CommandService) DeleteSession(ctx context.Context, id uint64) error {
	if _, err := s.getSession(ctx, id); err != nil {
		return err
	}
	if err := s.sessionRepo.DeleteByID(ctx, id); err != nil {
		return errors.DatabaseErrWithErr(err, "delete chat session failed")
	}
	return nil
}

func (s *CommandService) UpdateSessionSummary(ctx context.Context, cmd *UpdateSessionSummaryCommand) error {
	if cmd.SessionID == 0 {
		return errors.InvalidArg("session_id is required")
	}
	if _, err := s.getSession(ctx, cmd.SessionID); err != nil {
		return err
	}
	if err := s.sessionRepo.UpdateSummary(ctx, cmd.SessionID, cmd.Summary, cmd.SummarizedUntilMessageID); err != nil {
		return errors.DatabaseErrWithErr(err, "update chat session summary failed")
	}
	return nil
}

func (s *CommandService) Complete(ctx context.Context, cmd *CompletionCommand) (*CompletionDTO, error) {
	if strings.TrimSpace(cmd.Message) == "" {
		return nil, errors.InvalidArg("message is required")
	}

	session, agent, err := s.resolveSessionAndAgent(ctx, cmd)
	if err != nil {
		return nil, err
	}

	userMessage, err := s.saveUserMessage(ctx, session, cmd.Message)
	if err != nil {
		return nil, err
	}

	if err := s.maybeSummarize(ctx, session, agent); err != nil {
		return nil, errors.LLMErrWithErr(err, "compress conversation context failed")
	}

	messages, err := s.buildQueryHistoryMessages(ctx, session, userMessage.ID)
	if err != nil {
		return nil, err
	}

	if s.queryEngine == nil {
		return nil, errors.LLMErr("query engine is not initialized")
	}
	result, err := s.queryEngine.Run(ctx, query.QueryRequest{
		SessionID:     uintToString(session.ID),
		AgentID:       uintToString(agent.ID),
		Input:         cmd.Message,
		Messages:      messages,
		SystemPrompt:  agent.SystemPrompt,
		Summary:       session.Summary,
		Mode:          completionMode(cmd.Mode),
		ModelConfigID: agent.ModelConfigID,
		Temperature:   optionalPositiveFloat(agent.Temperature),
		TopP:          optionalPositiveFloat(agent.TopP),
		MaxTokens:     agent.MaxTokens,
		AllowedTools:  staticAgentTools(),
		WorkingDir:    defaultWorkingDir(),
	})
	if err != nil {
		return nil, errors.LLMErrWithErr(err, "query completion failed")
	}

	assistantContent := strings.TrimSpace(result.FinalAnswer)
	if assistantContent == "" {
		assistantContent = strings.TrimSpace(result.Content)
	}
	assistantMsg, err := s.saveAssistantMessageWithMeta(ctx, session, assistantContent, buildAssistantMeta(result))
	if err != nil {
		return nil, err
	}

	return &CompletionDTO{
		SessionID: session.ID,
		Message:   ToMessageDTO(assistantMsg),
		Content:   assistantContent,
	}, nil
}

func (s *CommandService) StreamComplete(ctx context.Context, cmd *CompletionCommand, onEvent func(event query.StreamEvent) error) (*CompletionDTO, error) {
	if strings.TrimSpace(cmd.Message) == "" {
		return nil, errors.InvalidArg("message is required")
	}

	session, agent, err := s.resolveSessionAndAgent(ctx, cmd)
	if err != nil {
		return nil, err
	}

	userMessage, err := s.saveUserMessage(ctx, session, cmd.Message)
	if err != nil {
		return nil, err
	}

	if err := s.maybeSummarize(ctx, session, agent); err != nil {
		return nil, errors.LLMErrWithErr(err, "compress conversation context failed")
	}

	messages, err := s.buildQueryHistoryMessages(ctx, session, userMessage.ID)
	if err != nil {
		return nil, err
	}

	if s.queryEngine == nil {
		return nil, errors.LLMErr("query engine is not initialized")
	}
	stream, err := s.queryEngine.Stream(ctx, query.QueryRequest{
		SessionID:     uintToString(session.ID),
		AgentID:       uintToString(agent.ID),
		Input:         cmd.Message,
		Messages:      messages,
		SystemPrompt:  agent.SystemPrompt,
		Summary:       session.Summary,
		Mode:          completionMode(cmd.Mode),
		ModelConfigID: agent.ModelConfigID,
		Temperature:   optionalPositiveFloat(agent.Temperature),
		TopP:          optionalPositiveFloat(agent.TopP),
		MaxTokens:     agent.MaxTokens,
		AllowedTools:  staticAgentTools(),
		WorkingDir:    defaultWorkingDir(),
		Stream:        true,
	})
	if err != nil {
		return nil, errors.LLMErrWithErr(err, "query stream failed")
	}

	var builder strings.Builder
	var finishedEvent *query.StreamEvent
	for event := range stream {
		if event.Err != nil {
			if onEvent != nil {
				_ = onEvent(event)
			}
			if stderrors.Is(event.Err, context.Canceled) {
				return nil, event.Err
			}
			return nil, errors.LLMErrWithErr(event.Err, "read query stream failed")
		}
		if event.Type == query.StreamEventLLMDelta {
			builder.WriteString(event.Content)
		}
		if event.Type == query.StreamEventRunFinished {
			copied := event
			finishedEvent = &copied
			continue
		}
		if onEvent != nil {
			if err := onEvent(event); err != nil {
				return nil, err
			}
		}
	}

	assistantContent := strings.TrimSpace(builder.String())
	assistantMsg, err := s.saveAssistantMessage(ctx, session, assistantContent)
	if err != nil {
		return nil, err
	}
	if finishedEvent != nil && onEvent != nil {
		finishedEvent.Content = assistantContent
		finishedEvent.Done = true
		if finishedEvent.Metadata == nil {
			finishedEvent.Metadata = map[string]interface{}{}
		}
		finishedEvent.Metadata["session_id"] = session.ID
		finishedEvent.Metadata["message_id"] = assistantMsg.ID
		if err := onEvent(*finishedEvent); err != nil {
			return nil, err
		}
	}

	return &CompletionDTO{
		SessionID: session.ID,
		Message:   ToMessageDTO(assistantMsg),
		Content:   assistantContent,
	}, nil
}

func (s *CommandService) CompactSession(ctx context.Context, cmd *CompactSessionCommand) (*SessionDTO, error) {
	if cmd.SessionID == 0 {
		return nil, errors.InvalidArg("session_id is required")
	}
	session, err := s.getSession(ctx, cmd.SessionID)
	if err != nil {
		return nil, err
	}
	agent, err := s.getEnabledAgent(ctx, session.AgentID)
	if err != nil {
		return nil, err
	}
	if s.llmExecutor == nil {
		return nil, errors.LLMErr("llm executor is not initialized")
	}

	activeMessages, err := s.messageRepo.FindBySessionIDAfterID(ctx, session.ID, session.SummarizedUntilMessageID, 0)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "load messages for compact failed")
	}
	if len(activeMessages) == 0 {
		return ToSessionDTO(session), nil
	}

	resp, result := s.llmExecutor.Chat(ctx, buildChatRequest(agent, []llm.Message{
		{Role: "user", Content: buildSummaryPrompt(session.Summary, activeMessages)},
	}))
	if !result.Success {
		return nil, errors.LLMErrWithErr(result.Error, "compact conversation context failed")
	}

	summary := strings.TrimSpace(resp.Content)
	untilID := activeMessages[len(activeMessages)-1].ID
	if err := s.sessionRepo.UpdateSummary(ctx, session.ID, summary, untilID); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "update compacted summary failed")
	}
	session.Summary = summary
	session.SummarizedUntilMessageID = untilID
	_ = s.refreshContextCache(ctx, session)
	return ToSessionDTO(session), nil
}

func (s *CommandService) saveUserMessage(ctx context.Context, session *domain.ChatSession, content string) (*domain.ChatMessage, error) {
	message := &domain.ChatMessage{
		SessionID:  session.ID,
		Role:       domain.MessageRoleUser,
		Content:    content,
		TokenCount: estimateTokens(content),
		Meta:       "{}",
	}
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "save user message failed")
	}
	_ = s.refreshContextCache(ctx, session)
	return message, nil
}

func (s *CommandService) saveAssistantMessage(ctx context.Context, session *domain.ChatSession, content string) (*domain.ChatMessage, error) {
	return s.saveAssistantMessageWithMeta(ctx, session, content, "{}")
}

func (s *CommandService) saveAssistantMessageWithMeta(ctx context.Context, session *domain.ChatSession, content string, meta string) (*domain.ChatMessage, error) {
	message := &domain.ChatMessage{
		SessionID:  session.ID,
		Role:       domain.MessageRoleAssistant,
		Content:    content,
		TokenCount: estimateTokens(content),
		Meta:       normalizeJSON(meta),
	}
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "save assistant message failed")
	}
	_ = s.refreshContextCache(ctx, session)
	return message, nil
}

func (s *CommandService) resolveSessionAndAgent(ctx context.Context, cmd *CompletionCommand) (*domain.ChatSession, *AgentRuntimeDTO, error) {
	if cmd.SessionID > 0 {
		session, err := s.getSession(ctx, cmd.SessionID)
		if err != nil {
			return nil, nil, err
		}
		agentID := session.AgentID
		if cmd.AgentID > 0 {
			agentID = cmd.AgentID
		}
		agent, err := s.getEnabledAgent(ctx, agentID)
		if err != nil {
			return nil, nil, err
		}
		return session, agent, nil
	}

	if cmd.AgentID == 0 {
		return nil, nil, errors.InvalidArg("agent_id is required when session_id is empty")
	}
	agent, err := s.getEnabledAgent(ctx, cmd.AgentID)
	if err != nil {
		return nil, nil, err
	}
	session := &domain.ChatSession{
		Title:         buildTitle(cmd.Message),
		AgentID:       agent.ID,
		ModelConfigID: agent.ModelConfigID,
		Status:        domain.SessionStatusActive,
		Meta:          "{}",
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, errors.DatabaseErrWithErr(err, "create chat session failed")
	}
	return session, agent, nil
}

func (s *CommandService) buildModelMessages(ctx context.Context, session *domain.ChatSession, agent *AgentRuntimeDTO) ([]llm.Message, error) {
	result := make([]llm.Message, 0, s.memoryConfig.CompressionThreshold+2)
	if strings.TrimSpace(agent.SystemPrompt) != "" {
		result = append(result, llm.Message{Role: "system", Content: agent.SystemPrompt})
	}
	if strings.TrimSpace(session.Summary) != "" {
		result = append(result, llm.Message{
			Role:    "system",
			Content: "Previous conversation summary:\n" + session.Summary,
		})
	}

	messages, err := s.loadActiveContext(ctx, session)
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		result = append(result, toLLMMessage(message))
	}
	return result, nil
}

func (s *CommandService) buildQueryHistoryMessages(ctx context.Context, session *domain.ChatSession, currentUserMessageID uint64) ([]llm.Message, error) {
	messages, err := s.loadActiveContext(ctx, session)
	if err != nil {
		return nil, err
	}
	result := make([]llm.Message, 0, len(messages))
	for _, message := range messages {
		if message.ID == currentUserMessageID {
			continue
		}
		result = append(result, toLLMMessage(message))
	}
	return result, nil
}

func (s *CommandService) maybeSummarize(ctx context.Context, session *domain.ChatSession, agent *AgentRuntimeDTO) error {
	activeMessages, err := s.messageRepo.FindBySessionIDAfterID(ctx, session.ID, session.SummarizedUntilMessageID, 0)
	if err != nil {
		return err
	}
	if len(activeMessages) <= s.memoryConfig.CompressionThreshold {
		return nil
	}

	retain := s.memoryConfig.RetainAfterCompress
	if retain < 1 {
		retain = 1
	}
	if retain >= len(activeMessages) {
		retain = len(activeMessages) - 1
	}
	compressCount := len(activeMessages) - retain
	if compressCount <= 0 {
		return nil
	}

	messagesToCompress := activeMessages[:compressCount]
	resp, result := s.llmExecutor.Chat(ctx, buildChatRequest(agent, []llm.Message{
		{Role: "user", Content: buildSummaryPrompt(session.Summary, messagesToCompress)},
	}))
	if !result.Success {
		return result.Error
	}

	summary := strings.TrimSpace(resp.Content)
	untilID := messagesToCompress[len(messagesToCompress)-1].ID
	if err := s.sessionRepo.UpdateSummary(ctx, session.ID, summary, untilID); err != nil {
		return err
	}
	session.Summary = summary
	session.SummarizedUntilMessageID = untilID
	_ = s.refreshContextCache(ctx, session)
	return nil
}

func (s *CommandService) refreshContextCache(ctx context.Context, session *domain.ChatSession) error {
	if s.redisClient == nil || session == nil {
		return nil
	}
	messages, err := s.messageRepo.FindBySessionIDAfterID(ctx, session.ID, session.SummarizedUntilMessageID, 0)
	if err != nil {
		return err
	}
	return redisutil.Set(ctx, s.redisClient, contextCacheKey(session.ID, session.SummarizedUntilMessageID), ToMessageDTOs(messages), s.memoryConfig.CacheTTL)
}

func (s *CommandService) loadActiveContext(ctx context.Context, session *domain.ChatSession) ([]*MessageDTO, error) {
	cacheKey := contextCacheKey(session.ID, session.SummarizedUntilMessageID)
	if s.redisClient != nil {
		cached, err := redisutil.Get[[]*MessageDTO](ctx, s.redisClient, cacheKey)
		if err == nil && cached != nil {
			return *cached, nil
		}
	}

	messages, err := s.messageRepo.FindBySessionIDAfterID(ctx, session.ID, session.SummarizedUntilMessageID, 0)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "load active context failed")
	}
	dtos := ToMessageDTOs(messages)
	if s.redisClient != nil {
		_ = redisutil.Set(ctx, s.redisClient, cacheKey, dtos, s.memoryConfig.CacheTTL)
	}
	return dtos, nil
}

func (s *CommandService) getSession(ctx context.Context, id uint64) (*domain.ChatSession, error) {
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "find chat session failed")
	}
	if session == nil {
		return nil, errors.NotFoundErr("chat session not found")
	}
	return session, nil
}

func (s *CommandService) getEnabledAgent(ctx context.Context, id uint64) (*AgentRuntimeDTO, error) {
	agent, err := s.agentGetter.GetAgentForChat(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "find agent failed")
	}
	if agent == nil {
		return nil, errors.NotFoundErr("agent not found")
	}
	if !agent.Enabled {
		return nil, errors.ForbiddenErr("agent is disabled")
	}
	if agent.ModelConfigID == 0 {
		return nil, errors.InvalidArg("agent model_config_id is required")
	}
	return agent, nil
}

func toLLMMessage(message *MessageDTO) llm.Message {
	return llm.Message{Role: message.Role, Content: message.Content}
}

func buildChatRequest(agent *AgentRuntimeDTO, messages []llm.Message) llm.ChatRequest {
	req := llm.ChatRequest{
		ModelConfigID: agent.ModelConfigID,
		Messages:      messages,
		MaxTokens:     agent.MaxTokens,
	}
	if agent.Temperature > 0 {
		req.Temperature = &agent.Temperature
	}
	if agent.TopP > 0 {
		req.TopP = &agent.TopP
	}
	return req
}

func optionalPositiveFloat(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

type assistantTraceItem struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Level       string                 `json:"level"`
	Iteration   int                    `json:"iteration,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Reason      string                 `json:"reason,omitempty"`
	ActionInput map[string]interface{} `json:"action_input,omitempty"`
	Observation string                 `json:"observation,omitempty"`
	IsError     bool                   `json:"is_error,omitempty"`
}

func buildAssistantMeta(result *query.QueryResult) string {
	if result == nil || len(result.Steps) == 0 {
		return "{}"
	}
	trace := make([]assistantTraceItem, 0, len(result.Steps))
	for i, step := range result.Steps {
		item := assistantTraceItem{
			ID:          i + 1,
			Title:       traceTitle(step),
			Content:     traceContent(step),
			Level:       traceLevel(step),
			Iteration:   step.Iteration,
			Action:      step.Name,
			Observation: step.Content,
			IsError:     step.IsError,
		}
		if step.Metadata != nil {
			if reason, ok := step.Metadata["reason"].(string); ok {
				item.Reason = reason
			}
			if action, ok := step.Metadata["action"].(string); ok && action != "" {
				item.Action = action
			}
			if input, ok := step.Metadata["action_input"].(map[string]interface{}); ok {
				item.ActionInput = input
			}
		}
		trace = append(trace, item)
	}
	data, err := json.Marshal(map[string]interface{}{"trace": trace})
	if err != nil {
		return "{}"
	}
	return string(data)
}

func traceTitle(step query.StepResult) string {
	if step.IsError {
		return "执行失败"
	}
	if step.Name == "final" {
		return "最终回答"
	}
	if strings.TrimSpace(step.Name) != "" {
		return "调用工具: " + step.Name
	}
	return "思考"
}

func traceContent(step query.StepResult) string {
	if step.Metadata != nil {
		if reason, ok := step.Metadata["reason"].(string); ok && strings.TrimSpace(reason) != "" {
			return reason
		}
	}
	if strings.TrimSpace(step.Content) != "" {
		return truncateTraceText(step.Content)
	}
	return "模型完成了一步推理。"
}

func traceLevel(step query.StepResult) string {
	if step.IsError {
		return "danger"
	}
	if step.Name == "final" {
		return "success"
	}
	return "info"
}

func truncateTraceText(text string) string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) <= 500 {
		return string(runes)
	}
	return string(runes[:500]) + "..."
}

func completionMode(mode string) query.ExecutionMode {
	switch strings.TrimSpace(mode) {
	case query.ExecutionModeFixedWorkflow.String(), "workflow":
		return query.ExecutionModeFixedWorkflow
	case query.ExecutionModePlanGraph.String(), "plan_graph":
		return query.ExecutionModePlanGraph
	case query.ExecutionModePureReAct.String(), "react":
		return query.ExecutionModePureReAct
	default:
		return query.ExecutionModeDirectChat
	}
}

func staticAgentTools() []string {
	tools := make([]string, len(staticAgentToolNames))
	copy(tools, staticAgentToolNames)
	return tools
}

func defaultWorkingDir() string {
	if workspace := strings.TrimSpace(os.Getenv("LATTICE_WORKSPACE_DIR")); workspace != "" {
		return workspace
	}
	wd, err := os.Getwd()
	if err != nil || strings.TrimSpace(wd) == "" {
		return "."
	}
	return wd
}

func buildSummaryPrompt(previous string, messages []*domain.ChatMessage) string {
	var b strings.Builder
	b.WriteString("Summarize the following conversation messages for future context. Keep user goals, decisions, constraints, facts, open questions, and any commitments. Do not invent details. Be concise but complete.\n")
	if previous != "" {
		b.WriteString("\nExisting summary:\n")
		b.WriteString(previous)
	}
	b.WriteString("\nMessages to compress:\n")
	for _, message := range messages {
		b.WriteString(string(message.Role))
		b.WriteString(": ")
		b.WriteString(message.Content)
		b.WriteString("\n")
	}
	return b.String()
}

func normalizeJSON(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return "{}"
	}
	return raw
}

func normalizeRole(role string) domain.MessageRole {
	if strings.TrimSpace(role) == "" {
		return domain.MessageRoleUser
	}
	return domain.MessageRole(role)
}

func contextCacheKey(sessionID uint64, summarizedUntilID uint64) string {
	return "chat:session:context:" + uintToString(sessionID) + ":" + uintToString(summarizedUntilID)
}

func buildTitle(message string) string {
	runes := []rune(strings.TrimSpace(message))
	if len(runes) > 40 {
		runes = runes[:40]
	}
	if len(runes) == 0 {
		return "New Chat"
	}
	return string(runes)
}

func estimateTokens(content string) int {
	runes := len([]rune(content))
	if runes == 0 {
		return 0
	}
	return (runes + 3) / 4
}

func isAllowedRole(role domain.MessageRole) bool {
	return role == domain.MessageRoleSystem ||
		role == domain.MessageRoleUser ||
		role == domain.MessageRoleAssistant ||
		role == domain.MessageRoleTool
}

func uintToString(v uint64) string {
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for v > 0 {
		buf = append(buf, byte('0'+v%10))
		v /= 10
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
