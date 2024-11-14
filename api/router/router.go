package router

import (
	"github.com/gin-gonic/gin"
)

func Register() *gin.Engine {
	r := gin.New()

	// 使用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// TODO: 添加API路由
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello, World!",
			})
		})
	}

	return r
}
