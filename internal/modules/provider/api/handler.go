package api

import (
	"strconv"

	"lattice-coding/internal/common/handler"
	"lattice-coding/internal/common/response"
	"lattice-coding/internal/modules/provider/application"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cmdSvc         *application.CommandService
	querySvc       *application.QueryService
	healthSvc      *application.HealthService
	healthCheckSvc *application.HealthCheckService
	syncSvc        *application.SyncService
}

func NewHandler(
	cmdSvc *application.CommandService,
	querySvc *application.QueryService,
	healthSvc *application.HealthService,
	healthCheckSvc *application.HealthCheckService,
	syncSvc *application.SyncService,
) *Handler {
	return &Handler{
		cmdSvc:         cmdSvc,
		querySvc:       querySvc,
		healthSvc:      healthSvc,
		healthCheckSvc: healthCheckSvc,
		syncSvc:        syncSvc,
	}
}

func (h *Handler) CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	cmd := ToCreateProviderCommand(&req)

	result, err := h.cmdSvc.CreateProvider(c.Request.Context(), cmd)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToProviderResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) ListProviders(c *gin.Context) {
	var query ProviderPageQuery
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

	appQuery := &application.ProviderPageQuery{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := h.querySvc.ListProviders(c.Request.Context(), appQuery)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	pageResp := ToProviderPageResponse(result)
	response.OkPage(c, pageResp.Items, pageResp.Total, pageResp.Page, pageResp.PageSize)
}

func (h *Handler) GetProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.querySvc.GetProvider(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToProviderResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) UpdateProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	var req UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	cmd := ToUpdateProviderCommand(&req)

	result, err := h.cmdSvc.UpdateProvider(c.Request.Context(), id, cmd)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToProviderResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) DeleteProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DeleteProvider(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) EnableProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.EnableProvider(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) DisableProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DisableProvider(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) TestProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.healthSvc.TestProvider(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ProviderTestResponse{
		Success:   result.Success,
		LatencyMs: result.LatencyMs,
		Error:     result.ErrorMessage,
	})
}

func (h *Handler) SyncProviderModels(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.syncSvc.SyncModels(c.Request.Context(), id, false)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, result)
}

func (h *Handler) HealthCheckProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.healthCheckSvc.CheckProvider(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ProviderHealthResponse{
		ProviderID:   result.ProviderID,
		Status:       result.Status,
		LatencyMs:    result.LatencyMs,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:    result.CheckedAt,
	})
}

func (h *Handler) GetProviderHealth(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.healthCheckSvc.CheckProvider(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ProviderHealthResponse{
		ProviderID:   result.ProviderID,
		Status:       result.Status,
		LatencyMs:    result.LatencyMs,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:    result.CheckedAt,
	})
}

func (h *Handler) CreateModelConfig(c *gin.Context) {
	var req CreateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	cmd := ToCreateModelConfigCommand(&req)

	result, err := h.cmdSvc.CreateModelConfig(c.Request.Context(), cmd)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToModelConfigResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) ListModelConfigs(c *gin.Context) {
	var query ModelConfigPageQuery
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

	appQuery := &application.ModelConfigPageQuery{
		Page:       query.Page,
		PageSize:   query.PageSize,
		ProviderID: query.ProviderID,
	}

	result, err := h.querySvc.ListModelConfigs(c.Request.Context(), appQuery)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	pageResp := ToModelConfigPageResponse(result)
	response.OkPage(c, pageResp.Items, pageResp.Total, pageResp.Page, pageResp.PageSize)
}

func (h *Handler) GetModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.querySvc.GetModelConfig(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToModelConfigResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) UpdateModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	var req UpdateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, err)
		return
	}

	cmd := ToUpdateModelConfigCommand(&req)

	result, err := h.cmdSvc.UpdateModelConfig(c.Request.Context(), id, cmd)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	resp := ToModelConfigResponse(result)
	response.Ok(c, resp)
}

func (h *Handler) DeleteModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DeleteModelConfig(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) EnableModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.EnableModelConfig(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) DisableModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.DisableModelConfig(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) TestModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.healthSvc.TestModelConfig(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ModelTestResponse{
		Success:   result.Success,
		LatencyMs: result.LatencyMs,
		Error:     result.ErrorMessage,
	})
}

func (h *Handler) SetDefaultModelConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.cmdSvc.SetDefaultModelConfig(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, (*interface{})(nil))
}

func (h *Handler) GetModelConfigHealth(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.healthCheckSvc.CheckModelConfig(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	response.Ok(c, ProviderHealthResponse{
		ProviderID:    result.ProviderID,
		ModelConfigID: result.ModelConfigID,
		Status:        result.Status,
		LatencyMs:     result.LatencyMs,
		ErrorCode:     result.ErrorCode,
		ErrorMessage:  result.ErrorMessage,
		CheckedAt:     result.CheckedAt,
	})
}
