package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	_ = api.Group("/v1/knowledge")
}
