package agent

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/agents")
	{
		r.GET("", listAgents)
		r.GET("/:id", getAgent)
		r.POST("", createAgent)
		r.PUT("/:id", updateAgent)
		r.DELETE("/:id", deleteAgent)
		r.POST("/:id/run", runAgent)
	}
}

func listAgents(c *gin.Context) {}

func getAgent(c *gin.Context) {}

func createAgent(c *gin.Context) {}

func updateAgent(c *gin.Context) {}

func deleteAgent(c *gin.Context) {}

func runAgent(c *gin.Context) {}
