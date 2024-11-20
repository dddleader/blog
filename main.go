package main

import (
	"blog/api"
	"blog/config"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 添加命令行参数
	var module string
	flag.StringVar(&module, "module", "api", "指定运行模块 (api)")
	flag.Parse()

	// 初始化配置
	config.Init()

	fmt.Printf("开始运行 %s 模块\n", module)

	// 根据模块参数启动相应服务
	switch module {
	case "api":
		server := api.New()
		server.Run()
	// case "site":
	// 	server := site.New()
	// 	server.Run()
	default:
		fmt.Println("退出：模块参数错误!")
		return
	}

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("服务器正在退出")
}
