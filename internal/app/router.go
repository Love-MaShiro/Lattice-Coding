package app

import (
	"net/http"

	"lattice-coding/internal/common/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(deps *Dependencies) (*gin.Engine, *gin.RouterGroup) {
	r := gin.New()

	r.Use(middleware.Logger(deps.Logger))
	r.Use(middleware.Recovery(deps.Logger))
	r.Use(middleware.CORS())
	r.Use(middleware.Trace())
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    "Lattice-coding is running",
		})
	})

	return r, api
}
