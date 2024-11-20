package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"blog/api/router"
	"blog/config"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run() {
	// 初始化路由
	r := router.Register()

	// 设置gin运行模式
	gin.SetMode(gin.ReleaseMode)
	logrus.Info("server starting...")

	// 获取配置
	conf := config.GetConfig()
	port := conf.ApiBase.ListenPort

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// 启动服务器
	go func() {
		fmt.Printf("server starting at port %d\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("listen error: %s\n", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Info("shutting down server...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("server shutdown error: %v", err)
	}

	logrus.Info("server exited")
	os.Exit(0)
}
