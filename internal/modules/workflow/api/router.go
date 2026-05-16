package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	api.POST("/v1/code-review/runs", h.RunCodeReview)

	r := api.Group("/v1/workflows")
	{
		r.POST("", h.CreateWorkflow)
		r.GET("", h.ListWorkflows)
		r.GET("/:id", h.GetWorkflow)
		r.PUT("/:id", h.UpdateWorkflow)
		r.DELETE("/:id", h.DeleteWorkflow)
	}
}
