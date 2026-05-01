package workflow

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/workflows")
	{
		api.GET("", listWorkflows)
		api.GET("/:id", getWorkflow)
		api.POST("", createWorkflow)
		api.PUT("/:id", updateWorkflow)
		api.DELETE("/:id", deleteWorkflow)
		api.POST("/:id/start", startWorkflow)
	}
}

func listWorkflows(c *gin.Context) {}

func getWorkflow(c *gin.Context) {}

func createWorkflow(c *gin.Context) {}

func updateWorkflow(c *gin.Context) {}

func deleteWorkflow(c *gin.Context) {}

func startWorkflow(c *gin.Context) {}
