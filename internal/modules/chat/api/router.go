package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	chat := api.Group("/v1/chat")
	{
		chat.POST("/sessions", h.CreateSession)
		chat.GET("/sessions", h.ListSessions)
		chat.GET("/sessions/:id", h.GetSession)
		chat.DELETE("/sessions/:id", h.DeleteSession)
		chat.POST("/sessions/:id/compact", h.CompactSession)
		chat.GET("/sessions/:id/messages", h.ListMessages)
		chat.POST("/messages", h.CreateMessage)
		chat.POST("/completions", h.CreateChatCompletion)
		chat.POST("/stream", h.CreateChatStream)
	}
}
