package safety

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/safety")
	{
		r.POST("/check/path", checkPath)
		r.POST("/check/command", checkCommand)
		r.POST("/check/permission", checkPermission)
	}
}

func checkPath(c *gin.Context) {}

func checkCommand(c *gin.Context) {}

func checkPermission(c *gin.Context) {}
