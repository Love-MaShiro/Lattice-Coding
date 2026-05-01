package provider

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/providers")
	{
		api.GET("", listProviders)
		api.GET("/:id", getProvider)
		api.POST("", createProvider)
		api.PUT("/:id", updateProvider)
		api.DELETE("/:id", deleteProvider)
		api.POST("/:id/test", testProvider)
	}
}

func listProviders(c *gin.Context) {}

func getProvider(c *gin.Context) {}

func createProvider(c *gin.Context) {}

func updateProvider(c *gin.Context) {}

func deleteProvider(c *gin.Context) {}

func testProvider(c *gin.Context) {}
