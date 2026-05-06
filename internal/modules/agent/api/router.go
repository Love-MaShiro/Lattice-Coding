package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	agents := api.Group("/v1/agents")
	{
		agents.POST("", h.CreateAgent)
		agents.GET("", h.ListAgents)
		agents.GET("/:id", h.GetAgent)
		agents.PUT("/:id", h.UpdateAgent)
		agents.DELETE("/:id", h.DeleteAgent)
		agents.POST("/:id/enable", h.EnableAgent)
		agents.POST("/:id/disable", h.DisableAgent)
	}
}
