package api

import (
	"lattice-coding/internal/common/handler"
	"lattice-coding/internal/common/response"
	"lattice-coding/internal/modules/run/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	querySvc *application.QueryService
}

func NewHandler(querySvc *application.QueryService) *Handler {
	return &Handler{querySvc: querySvc}
}

func (h *Handler) ListRuns(c *gin.Context) {
	page := response.ParsePageRequest(c)
	result, err := h.querySvc.ListRuns(c.Request.Context(), application.RunPageQuery{
		Page:     page.Page,
		PageSize: page.Size,
	})
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.OkPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

func (h *Handler) GetRun(c *gin.Context) {
	result, err := h.querySvc.GetRun(c.Request.Context(), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, result)
}

func (h *Handler) ListToolInvocations(c *gin.Context) {
	result, err := h.querySvc.ListToolInvocations(c.Request.Context(), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, result)
}
