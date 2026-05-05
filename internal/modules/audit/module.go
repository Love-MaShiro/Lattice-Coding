package audit

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/audit")
	{
		r.GET("/logs", listLogs)
		r.GET("/logs/:id", getLog)
		r.GET("/llm-calls", listLLMCalls)
		r.GET("/tool-calls", listToolCalls)
	}
}

func listLogs(c *gin.Context) {}

func getLog(c *gin.Context) {}

func listLLMCalls(c *gin.Context) {}

func listToolCalls(c *gin.Context) {}
