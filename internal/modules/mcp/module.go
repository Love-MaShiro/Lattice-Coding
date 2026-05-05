package mcp

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/tools")
	{
		r.GET("", listTools)
		r.GET("/:id", getTool)
		r.POST("/:id/execute", executeTool)
	}
}

func listTools(c *gin.Context) {}

func getTool(c *gin.Context) {}

func executeTool(c *gin.Context) {}
