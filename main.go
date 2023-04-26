package main

import (
	"Gateway-Go/common/lib"
	"Gateway-Go/dao"
	"Gateway-Go/grpc_proxy_router"
	"Gateway-Go/http_proxy_router"
	"Gateway-Go/router"
	"Gateway-Go/tcp_proxy_router"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

//main.go混用，传入两个参数，config和endpoint
//endpoint dashboard后台管理 server代理服务器
//config ./conf/prod 对应配置文件夹

var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	//go run main.go -config=./conf/dev -endpoint=server
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}
	//判断endpoint
	if *endpoint == "dashboard" {
		lib.InitModule("./conf/dev/")
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		//执行server逻辑，启动代理服务器
		lib.InitModule("./conf/dev/")
		defer lib.Destroy()

		dao.ServiceManagerHandler.LoadOnce()
		dao.AppManagerHandler.LoadOnce()

		//使用协程防止多个服务器启动阻塞
		go func() {
			http_proxy_router.HttpServerRun()
		}()

		go func() {
			http_proxy_router.HttpsServerRun()
		}()

		go func() {
			tcp_proxy_router.TcpServerRun()
		}()

		go func() {
			grpc_proxy_router.GrpcServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		tcp_proxy_router.TcpServerStop()
		grpc_proxy_router.GrpcServerStop()
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}

}
