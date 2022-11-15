package http

import (
	"github.com/gohade/hade/app/http/module/demo"
	"github.com/gohade/hade/framework/gin"
	"github.com/gohade/hade/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	r.Use(static.Serve("/",static.LocalFile("./dist",false)))

	demo.Register(r)
}
