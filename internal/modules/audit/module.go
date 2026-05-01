package audit

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/audit")
	{
		api.GET("/logs", listLogs)
		api.GET("/logs/:id", getLog)
		api.GET("/llm-calls", listLLMCalls)
		api.GET("/tool-calls", listToolCalls)
	}
}

func listLogs(c *gin.Context) {}

func getLog(c *gin.Context) {}

func listLLMCalls(c *gin.Context) {}

func listToolCalls(c *gin.Context) {}
