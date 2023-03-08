package http

import (
	"github.com/casbin/casbin/v3"
	"github.com/gohade/hade/app/http/consumer"
	"github.com/gohade/hade/app/http/middleware/auth"
	"github.com/gohade/hade/app/http/module/pk"
	"github.com/gohade/hade/app/http/module/ranklist"
	"github.com/gohade/hade/app/http/module/record"
	"github.com/gohade/hade/app/http/module/user/account"
	"github.com/gohade/hade/app/http/module/user/bot"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
	swaggerFiles "github.com/gohade/hade/framework/middleware/files"
	ginSwagger "github.com/gohade/hade/framework/middleware/gin-swagger"
	"github.com/gohade/hade/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	if configService.GetBool("app.swagger") {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	e, _ := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")
	r.Use(auth.AuthMiddleware(e))

	account.Register(r)
	bot.Register(r)
	consumer.Register(r)
	pk.Register(r)
	record.Register(r)
	ranklist.Register(r)
}
