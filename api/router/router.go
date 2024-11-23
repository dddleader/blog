package router

import (
	"blog/api/middleware"
	"blog/article"
	_ "blog/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api")
	{
		// 公开接口
		api.GET("/articles", article.GetArticles)
		api.GET("/articles/:id", article.GetSingleArticle)

		// 管理员接口
		admin := api.Group("/admin")
		{
			admin.POST("/login", article.AdminLogin)

			// 需要认证的接口
			authorized := admin.Group("")
			authorized.Use(middleware.AuthRequired())
			{
				// 文章管理接口
				authorized.POST("/articles", article.CreateArticle)
				authorized.PUT("/articles/:id", article.UpdateArticle)
				authorized.DELETE("/articles/:id", article.DeleteArticle)
			}
		}
	}

	return r
}
