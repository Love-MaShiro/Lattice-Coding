package safety

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/safety")
	{
		api.POST("/check/path", checkPath)
		api.POST("/check/command", checkCommand)
		api.POST("/check/permission", checkPermission)
	}
}

func checkPath(c *gin.Context) {}

func checkCommand(c *gin.Context) {}

func checkPermission(c *gin.Context) {}
