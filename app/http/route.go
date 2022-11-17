package http

import (
	"github.com/gohade/hade/app/http/module/demo"

	"github.com/gohade/hade/app/http/middleware/cors"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
	swaggerFiles "github.com/gohade/hade/framework/middleware/files"
	ginSwagger "github.com/gohade/hade/framework/middleware/gin-swagger"
	"github.com/gohade/hade/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	r.Use(static.Serve("/",static.LocalFile("./dist",false)))

	if configService.GetBool("app.swagger") {
		r.GET("/swagger/*any",ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.Use(cors.Default())
	demo.Register(r)
}
