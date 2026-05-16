package api

import (
	"strconv"

	"lattice-coding/internal/common/handler"
	"lattice-coding/internal/common/response"
	"lattice-coding/internal/modules/workflow/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cmdSvc        *application.CommandService
	querySvc      *application.QueryService
	codeReviewSvc *application.CodeReviewService
}

func NewHandler(cmdSvc *application.CommandService, querySvc *application.QueryService, codeReviewSvc *application.CodeReviewService) *Handler {
	return &Handler{cmdSvc: cmdSvc, querySvc: querySvc, codeReviewSvc: codeReviewSvc}
}

func (h *Handler) CreateWorkflow(c *gin.Context) {
	var req SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}
	result, err := h.cmdSvc.CreateWorkflow(c.Request.Context(), req.ToApplication())
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, ToWorkflowResponse(result))
}

func (h *Handler) ListWorkflows(c *gin.Context) {
	var query WorkflowPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handler.HandleError(c, err)
		return
	}
	result, err := h.querySvc.ListWorkflows(c.Request.Context(), query.ToApplication())
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.OkPage(c, ToWorkflowResponses(result.Items), result.Total, result.Page, result.PageSize)
}

func (h *Handler) GetWorkflow(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	result, err := h.querySvc.GetWorkflow(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, ToWorkflowResponse(result))
}

func (h *Handler) UpdateWorkflow(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	var req SaveWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}
	result, err := h.cmdSvc.UpdateWorkflow(c.Request.Context(), id, req.ToApplication())
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, ToWorkflowResponse(result))
}

func (h *Handler) DeleteWorkflow(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	if err := h.cmdSvc.DeleteWorkflow(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok[any](c, nil)
}

func (h *Handler) RunCodeReview(c *gin.Context) {
	var req CodeReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}
	result, err := h.codeReviewSvc.Run(c.Request.Context(), req.ToApplication())
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, result)
}

func parseID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
