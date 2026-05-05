package knowledge

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/knowledge")
	{
		r.GET("/documents", listDocuments)
		r.GET("/documents/:id", getDocument)
		r.POST("/documents", createDocument)
		r.PUT("/documents/:id", updateDocument)
		r.DELETE("/documents/:id", deleteDocument)
		r.POST("/search", searchDocuments)
	}
}

func listDocuments(c *gin.Context) {}

func getDocument(c *gin.Context) {}

func createDocument(c *gin.Context) {}

func updateDocument(c *gin.Context) {}

func deleteDocument(c *gin.Context) {}

func searchDocuments(c *gin.Context) {}
