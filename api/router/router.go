package router

import (
	"blog/api/middleware"
	"blog/article"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register() *gin.Engine {
	r := gin.New()

	// 使用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// 或者使用 StaticFS 更灵活地配置
	r.StaticFS("/static", http.Dir("static"))

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
		api.GET("/articles", article.GetArticles)
		api.GET("/single-article/:path", article.GetSingleArticle)
	}

	return r
}
