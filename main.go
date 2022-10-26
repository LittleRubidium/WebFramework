package main

import (
	"context"
	"github.com/gohade/hade/framework/gin"
	"github.com/gohade/hade/framework/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := gin.New()
	//core.Use(middleware.Test1(),middleware.Test2())
	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	//开启一个协程处理请求，否则main函数会阻塞在这里
	go server.ListenAndServe()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
