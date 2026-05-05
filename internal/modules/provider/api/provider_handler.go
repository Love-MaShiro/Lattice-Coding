package api

import (
	"strconv"
	"time"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/common/handler"
	"lattice-coding/internal/common/response"
	"lattice-coding/internal/modules/provider/application"
	"lattice-coding/internal/runtime/llm"

	"github.com/gin-gonic/gin"
)

// ProviderHandler provider HTTP 处理器
type ProviderHandler struct {
	providerService *application.ProviderService
	llmFactory      *llm.LLMFactory
}

func NewProviderHandler(providerService *application.ProviderService, llmFactory *llm.LLMFactory) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
		llmFactory:      llmFactory,
	}
}

func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	ctx := c.Request.Context()

	var req application.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	provider, err := h.providerService.CreateProvider(ctx, &req)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, provider)
}

func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	var req application.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	provider, err := h.providerService.UpdateProvider(ctx, id, &req)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, provider)
}

func (h *ProviderHandler) GetProvider(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	provider, modelConfigs, err := h.providerService.GetProviderWithModelConfigs(ctx, id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, gin.H{"provider": provider, "model_configs": modelConfigs})
}

func (h *ProviderHandler) ListProviders(c *gin.Context) {
	ctx := c.Request.Context()

	providers, err := h.providerService.ListProviders(ctx)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, providers)
}

func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	if err := h.providerService.DeleteProvider(ctx, id); err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok[any](c, nil)
}

func (h *ProviderHandler) EnableProvider(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	if err := h.providerService.EnableProvider(ctx, id); err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok[any](c, nil)
}

func (h *ProviderHandler) DisableProvider(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	if err := h.providerService.DisableProvider(ctx, id); err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok[any](c, nil)
}

func (h *ProviderHandler) TestProvider(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	startTime := time.Now()
	success, testErr := h.llmFactory.TestProvider(ctx, id)
	latency := time.Since(startTime)

	result := gin.H{
		"success":    success,
		"latency_ms": latency.Milliseconds(),
	}

	if testErr != nil {
		result["error"] = testErr.Error()
	}

	response.Ok(c, result)
}

func (h *ProviderHandler) CreateModelConfig(c *gin.Context) {
	ctx := c.Request.Context()

	var req application.CreateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	modelConfig, err := h.providerService.CreateModelConfig(ctx, &req)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, modelConfig)
}

func (h *ProviderHandler) ListModelConfigs(c *gin.Context) {
	ctx := c.Request.Context()

	modelConfigs, err := h.providerService.ListModelConfigs(ctx)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	response.Ok(c, modelConfigs)
}

func (h *ProviderHandler) TestModelConfig(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errors.InvalidArg(err.Error()))
		return
	}

	startTime := time.Now()
	success, testErr := h.llmFactory.TestModel(ctx, id)
	latency := time.Since(startTime)

	result := gin.H{
		"success":    success,
		"latency_ms": latency.Milliseconds(),
	}

	if testErr != nil {
		result["error"] = testErr.Error()
	}

	response.Ok(c, result)
}
