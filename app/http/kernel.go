package http

import (
	"github.com/gohade/hade/app/http/middleware/cors"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/gin"
)

func NewHttpEngine(c framework.Container) (*gin.Engine, error) {
	//设置为Release，默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	//默认启动一个引擎
	r := gin.New()
	r.Use(cors.Cors)
	r.SetContainer(c)

	r.Use(gin.Recovery())

	Routes(r)

	return r, nil
}
