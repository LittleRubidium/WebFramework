package http

import "github.com/gohade/hade/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	//设置为Release，默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	//默认启动一个引擎
	r := gin.Default()
	Routes(r)

	return r,nil
}
