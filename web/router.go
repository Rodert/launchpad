package web

import (
	"github.com/gin-gonic/gin"
	"launchpad/action/ping"
	"launchpad/internal/handler"
	"launchpad/logger"
)

func Router() *gin.Engine {
	router := gin.New()

	// ping group
	rootGroup := router.Group("")
	{
		// 测试接口
		pingGroup := rootGroup.Group("ping")
		{
			pingGroup.GET("info", handler.TRPathParamHandler(ping.GetPing))
		}
	}

	// api
	apiV1Group := router.Group("/api/v1")
	apiV1Group.Use(logger.SetModule("api"), logger.AccessLog())

	kytGroup := apiV1Group.Group("demo")
	{
		kytGroup.GET("info", handler.TRPathParamHandler(ping.GetPing))
	}

	return router
}
