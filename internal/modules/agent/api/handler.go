package api

import (
	"strconv"

	"lattice-coding/internal/common/handler"
	"lattice-coding/internal/common/response"
	"lattice-coding/internal/modules/agent/application"

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

func (h *Handler) CreateAgent(c *gin.Context) {
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	cmd := ToCreateAgentCommand(&req)

	result, err := h.cmdSvc.CreateAgent(c.Request.Context(), cmd)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToAgentResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) GetAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.querySvc.GetAgent(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToAgentResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) ListAgents(c *gin.Context) {
	var query AgentPageQuery
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

	appQuery := &application.AgentPageQuery{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := h.querySvc.ListAgents(c.Request.Context(), appQuery)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	pageResp := ToAgentPageResponse(result)
	response.OkPage(c, pageResp.Items, pageResp.Total, pageResp.Page, pageResp.PageSize)
}

func (h *Handler) UpdateAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	var req UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	cmd := ToUpdateAgentCommand(&req)

	result, err := h.cmdSvc.UpdateAgent(c.Request.Context(), id, cmd)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToAgentResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) DeleteAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DeleteAgent(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok[any](c, nil)
}

func (h *Handler) EnableAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.EnableAgent(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok[any](c, nil)
}

func (h *Handler) DisableAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DisableAgent(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok[any](c, nil)
}
