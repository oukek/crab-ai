package router

import (
	"github.com/gin-gonic/gin"
	"oukek/crab-ai/app/middleware"
)

var routerInstance *gin.Engine

func init() {
	getRouterInstance()
}

// 获取路由实例，会进行初始化
func getRouterInstance() *gin.Engine {
	if routerInstance != nil {
		return routerInstance
	}
	routerInstance = gin.Default()

	// 注册中间件
	{
		// 异常处理
		routerInstance.Use(middleware.Recover)
		// 跨域
		routerInstance.Use(middleware.Cors())
	}
	return routerInstance
}

func GetRouter() *gin.Engine {
	r := routerInstance
	r.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})
	return routerInstance
}
