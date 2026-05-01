package agent

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/agents")
	{
		api.GET("", listAgents)
		api.GET("/:id", getAgent)
		api.POST("", createAgent)
		api.PUT("/:id", updateAgent)
		api.DELETE("/:id", deleteAgent)
		api.POST("/:id/run", runAgent)
	}
}

func listAgents(c *gin.Context) {}

func getAgent(c *gin.Context) {}

func createAgent(c *gin.Context) {}

func updateAgent(c *gin.Context) {}

func deleteAgent(c *gin.Context) {}

func runAgent(c *gin.Context) {}
