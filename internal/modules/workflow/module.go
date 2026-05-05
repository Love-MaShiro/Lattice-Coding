package workflow

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/workflows")
	{
		r.GET("", listWorkflows)
		r.GET("/:id", getWorkflow)
		r.POST("", createWorkflow)
		r.PUT("/:id", updateWorkflow)
		r.DELETE("/:id", deleteWorkflow)
		r.POST("/:id/start", startWorkflow)
	}
}

func listWorkflows(c *gin.Context) {}

func getWorkflow(c *gin.Context) {}

func createWorkflow(c *gin.Context) {}

func updateWorkflow(c *gin.Context) {}

func deleteWorkflow(c *gin.Context) {}

func startWorkflow(c *gin.Context) {}
