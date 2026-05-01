package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/chat")
	{
		api.POST("/completions", createChatCompletion)
		api.POST("/stream", createChatStream)
		api.GET("/messages/:conversation_id", getMessages)
		api.POST("/messages", createMessage)
	}
}

func createChatCompletion(c *gin.Context) {}

func createChatStream(c *gin.Context) {}

func getMessages(c *gin.Context) {}

func createMessage(c *gin.Context) {}
