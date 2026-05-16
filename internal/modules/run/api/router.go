package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	runs := api.Group("/v1/runs")
	{
		runs.GET("", h.ListRuns)
		runs.GET("/:id", h.GetRun)
		runs.GET("/:id/tool-invocations", h.ListToolInvocations)
	}
}
