package main

import (
	"os"
	"os/signal"
	"syscall"

	lib "github.com/qq65326/apns_push/lib"
	push "github.com/qq65326/apns_push/push"
	server "github.com/qq65326/apns_push/server"
)

var (
	Host = "0.0.0.0"
	Port = "8080"
)

func main() {
	var sigCh = make(chan os.Signal, 1)

	mainServer := server.MainServer(Host, Port)
	mainServer.Start()
	// 为了不阻塞后续流程，启用协程开启推送
	go push.StartWorker()

	// 监听系统信号
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		<-sigCh
		mainServer.Stop()
		push.Stop()
		lib.MainWaitGroup.Wait()
		break
	}

}
