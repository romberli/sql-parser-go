package router

import (
	"github.com/gin-gonic/gin"
	"github.com/romberli/go-template/api/v1/health"
)

// RegisterHealth is the sub-router for health
func RegisterHealth(group *gin.RouterGroup) {
	healthGroup := group.Group("/health")
	{
		healthGroup.GET("/ping", health.Ping)
	}
}
