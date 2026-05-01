package mcp

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/tools")
	{
		api.GET("", listTools)
		api.GET("/:id", getTool)
		api.POST("/:id/execute", executeTool)
	}
}

func listTools(c *gin.Context) {}

func getTool(c *gin.Context) {}

func executeTool(c *gin.Context) {}
