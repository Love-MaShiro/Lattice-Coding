package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/chat")
	{
		r.POST("/completions", createChatCompletion)
		r.POST("/stream", createChatStream)
		r.GET("/messages/:conversation_id", getMessages)
		r.POST("/messages", createMessage)
	}
}

func createChatCompletion(c *gin.Context) {}

func createChatStream(c *gin.Context) {}

func getMessages(c *gin.Context) {}

func createMessage(c *gin.Context) {}
