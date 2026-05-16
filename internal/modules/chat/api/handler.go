package api

import (
	"strconv"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/common/handler"
	"lattice-coding/internal/common/response"
	"lattice-coding/internal/modules/chat/application"
	"lattice-coding/internal/runtime/query"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cmdSvc   *application.CommandService
	querySvc *application.QueryService
}

func NewHandler(cmdSvc *application.CommandService, querySvc *application.QueryService) *Handler {
	return &Handler{
		cmdSvc:   cmdSvc,
		querySvc: querySvc,
	}
}

func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.cmdSvc.CreateSession(c.Request.Context(), ToCreateSessionCommand(&req))
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ToSessionResponse(result))
}

func (h *Handler) ListSessions(c *gin.Context) {
	var query SessionPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handler.HandleError(c, err)
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	result, err := h.querySvc.ListSessions(c.Request.Context(), &application.SessionPageQuery{
		Page:     query.Page,
		PageSize: query.PageSize,
	})
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	pageResp := ToSessionPageResponse(result)
	response.OkPage(c, pageResp.Items, pageResp.Total, pageResp.Page, pageResp.PageSize)
}

func (h *Handler) GetSession(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.querySvc.GetSession(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ToSessionResponse(result))
}

func (h *Handler) DeleteSession(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DeleteSession(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok[any](c, nil)
}

func (h *Handler) CompactSession(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.cmdSvc.CompactSession(c.Request.Context(), &application.CompactSessionCommand{
		SessionID: id,
	})
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ToSessionResponse(result))
}

func (h *Handler) ListMessages(c *gin.Context) {
	sessionID, err := parseUintParam(c, "id")
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	var query MessageListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.querySvc.ListMessages(c.Request.Context(), &application.MessageQuery{
		SessionID: sessionID,
		Limit:     query.Limit,
	})
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ToMessageResponses(result))
}

func (h *Handler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.cmdSvc.CreateMessage(c.Request.Context(), ToCreateMessageCommand(&req))
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ToMessageResponse(result))
}

func (h *Handler) CreateChatCompletion(c *gin.Context) {
	var req CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.cmdSvc.Complete(c.Request.Context(), ToCompletionCommand(&req))
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ToCompletionResponse(result))
}

func (h *Handler) CreateChatStream(c *gin.Context) {
	var req CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(200)

	_, err := h.cmdSvc.StreamComplete(c.Request.Context(), ToCompletionCommand(&req), func(event query.StreamEvent) error {
		payload := gin.H{
			"type":    event.Type,
			"run_id":  event.RunID,
			"content": event.Content,
			"done":    event.Done,
		}
		if event.Err != nil {
			payload["message"] = event.Err.Error()
		}
		if event.Metadata != nil {
			payload["metadata"] = event.Metadata
		}
		c.SSEvent(string(event.Type), payload)
		c.Writer.Flush()
		return nil
	})
	if err != nil {
		c.SSEvent("error", gin.H{"message": err.Error()})
		c.Writer.Flush()
		return
	}
}

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, errors.InvalidArg(name + " is invalid")
	}
	return id, nil
}
